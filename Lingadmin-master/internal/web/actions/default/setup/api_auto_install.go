package setup

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var apiBinaryNames = []string{"ling-api", "sk-api", "edge-api"}

func ensureLocalAPINodeInstalled(apiNodeDir string) error {
	if len(apiNodeDir) == 0 {
		return errors.New("API 节点目录为空")
	}

	// 确保目录存在
	err := os.MkdirAll(filepath.Join(apiNodeDir, "configs"), 0755)
	if err != nil {
		return fmt.Errorf("创建 API configs 目录失败：%w", err)
	}
	err = os.MkdirAll(filepath.Join(apiNodeDir, "bin"), 0755)
	if err != nil {
		return fmt.Errorf("创建 API bin 目录失败：%w", err)
	}

	// 已存在可执行文件则跳过
	for _, name := range apiBinaryNames {
		stat, statErr := os.Stat(filepath.Join(apiNodeDir, "bin", name))
		if statErr == nil && stat.Mode().IsRegular() {
			return nil
		}
	}

	downloadURL, version, expectedMD5, expectedSHA256, err := fetchLatestComponentPackage("api")
	if err != nil {
		return err
	}
	if len(downloadURL) == 0 {
		return errors.New("获取到的下载地址为空")
	}

	currentStatusText = "正在下载 API 组件 v" + version
	archivePath, err := downloadFileAndVerify(downloadURL, expectedMD5, expectedSHA256)
	if err != nil {
		return err
	}
	defer func() {
		_ = os.Remove(archivePath)
	}()

	currentStatusText = "正在部署 API 组件"
	installedName, err := extractBinaryFromArchive(archivePath, filepath.Join(apiNodeDir, "bin"), apiBinaryNames)
	if err != nil {
		return err
	}

	// 确保可执行权限
	_ = os.Chmod(filepath.Join(apiNodeDir, "bin", installedName), 0755)
	logs.Println("[SETUP]installed api binary:", installedName)
	return nil
}

func fetchLatestComponentPackage(componentCode string) (downloadURL string, version string, md5Value string, sha256Value string, err error) {
	currentStatusText = "正在获取组件下载信息"

	url := strings.TrimSpace(os.Getenv("LINGCDN_UPDATES_URL"))
	if len(url) == 0 {
		url = teaconst.UpdatesURL
	}
	url = strings.ReplaceAll(url, "${os}", runtime.GOOS)
	url = strings.ReplaceAll(url, "${arch}", runtime.GOARCH)

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", "", "", "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", "", fmt.Errorf("读取版本信息失败：%w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return "", "", "", "", fmt.Errorf("读取版本信息失败：HTTP %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", "", fmt.Errorf("读取版本信息失败：%w", err)
	}

	var m = maps.Map{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		return "", "", "", "", fmt.Errorf("解析版本信息失败：%w", err)
	}
	if m.GetInt("code") != 200 {
		msg := m.GetString("message")
		if len(msg) == 0 {
			msg = "版本接口返回异常"
		}
		return "", "", "", "", errors.New(msg)
	}

	dataMap := m.GetMap("data")
	host := dataMap.GetString("host")
	versions := dataMap.GetSlice("versions")

	var downloadPath string
	for _, v := range versions {
		vMap := maps.NewMap(v)
		if vMap.GetString("code") == componentCode {
			version = vMap.GetString("version")
			downloadPath = vMap.GetString("url")
			md5Value = vMap.GetString("md5")
			sha256Value = vMap.GetString("sha256")
			break
		}
	}

	if len(downloadPath) == 0 {
		return "", "", "", "", errors.New("未找到组件 '" + componentCode + "' 的下载信息")
	}

	downloadURL = joinURL(host, downloadPath)
	return downloadURL, version, md5Value, sha256Value, nil
}

func joinURL(host string, downloadPath string) string {
	if strings.HasPrefix(downloadPath, "http://") || strings.HasPrefix(downloadPath, "https://") {
		return downloadPath
	}
	if len(host) == 0 {
		return downloadPath
	}
	if strings.HasSuffix(host, "/") && strings.HasPrefix(downloadPath, "/") {
		return host + downloadPath[1:]
	}
	if !strings.HasSuffix(host, "/") && !strings.HasPrefix(downloadPath, "/") {
		return host + "/" + downloadPath
	}
	return host + downloadPath
}

func downloadFileAndVerify(downloadURL string, expectedMD5 string, expectedSHA256 string) (string, error) {
	currentStatusText = "正在下载组件包"

	suffix := ""
	lower := strings.ToLower(downloadURL)
	if strings.HasSuffix(lower, ".tar.gz") {
		suffix = ".tar.gz"
	} else if strings.HasSuffix(lower, ".zip") {
		suffix = ".zip"
	}

	tmpFile, err := os.CreateTemp("", "lingcdn-download-*"+suffix)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = tmpFile.Close()
	}()

	client := &http.Client{Timeout: 30 * time.Minute}
	req, err := http.NewRequest(http.MethodGet, downloadURL, nil)
	if err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", fmt.Errorf("下载失败：%w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		_ = os.Remove(tmpFile.Name())
		return "", fmt.Errorf("下载失败：HTTP %d", resp.StatusCode)
	}

	md5Sum := md5.New()
	sha256Sum := sha256.New()
	writer := io.MultiWriter(tmpFile, md5Sum, sha256Sum)
	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", fmt.Errorf("下载失败：%w", err)
	}

	actualMD5 := hex.EncodeToString(md5Sum.Sum(nil))
	actualSHA256 := hex.EncodeToString(sha256Sum.Sum(nil))

	expectedMD5 = strings.ToLower(strings.TrimSpace(expectedMD5))
	expectedSHA256 = strings.ToLower(strings.TrimSpace(expectedSHA256))

	if len(expectedSHA256) > 0 && expectedSHA256 != "null" && actualSHA256 != expectedSHA256 {
		_ = os.Remove(tmpFile.Name())
		return "", fmt.Errorf("文件校验失败（sha256），期望 %s，实际 %s", expectedSHA256, actualSHA256)
	}
	if len(expectedMD5) > 0 && expectedMD5 != "null" && actualMD5 != expectedMD5 {
		_ = os.Remove(tmpFile.Name())
		return "", fmt.Errorf("文件校验失败（md5），期望 %s，实际 %s", expectedMD5, actualMD5)
	}

	return tmpFile.Name(), nil
}

func extractBinaryFromArchive(archivePath string, destDir string, candidates []string) (string, error) {
	err := os.MkdirAll(destDir, 0755)
	if err != nil {
		return "", err
	}

	lower := strings.ToLower(archivePath)
	if strings.HasSuffix(lower, ".zip") {
		return extractBinaryFromZip(archivePath, destDir, candidates)
	}
	if strings.HasSuffix(lower, ".tar.gz") {
		return extractBinaryFromTarGz(archivePath, destDir, candidates)
	}

	// 尝试按zip处理
	if name, zipErr := extractBinaryFromZip(archivePath, destDir, candidates); zipErr == nil {
		return name, nil
	}
	// 再尝试tar.gz
	return extractBinaryFromTarGz(archivePath, destDir, candidates)
}

func extractBinaryFromZip(zipPath string, destDir string, candidates []string) (string, error) {
	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = zr.Close()
	}()

	var target *zip.File
	var targetName string
	for _, name := range candidates {
		for _, f := range zr.File {
			if f.FileInfo().IsDir() {
				continue
			}
			if path.Base(f.Name) == name {
				target = f
				targetName = name
				break
			}
		}
		if target != nil {
			break
		}
	}
	if target == nil {
		return "", errors.New("解压失败：在 zip 包中未找到 API 可执行文件")
	}

	rc, err := target.Open()
	if err != nil {
		return "", err
	}
	defer func() {
		_ = rc.Close()
	}()

	dstPath := filepath.Join(destDir, targetName)
	tmpPath := dstPath + ".tmp"
	_ = os.Remove(tmpPath)

	fp, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(fp, rc)
	_ = fp.Close()
	if err != nil {
		_ = os.Remove(tmpPath)
		return "", err
	}

	_ = os.Remove(dstPath)
	err = os.Rename(tmpPath, dstPath)
	if err != nil {
		_ = os.Remove(tmpPath)
		return "", err
	}

	return targetName, nil
}

func extractBinaryFromTarGz(tarGzPath string, destDir string, candidates []string) (string, error) {
	fp, err := os.Open(tarGzPath)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = fp.Close()
	}()

	gz, err := gzip.NewReader(fp)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = gz.Close()
	}()

	tr := tar.NewReader(gz)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		if h.Typeflag != tar.TypeReg {
			continue
		}

		base := path.Base(h.Name)
		for _, name := range candidates {
			if base == name {
				dstPath := filepath.Join(destDir, name)
				tmpPath := dstPath + ".tmp"
				_ = os.Remove(tmpPath)

				out, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
				if err != nil {
					return "", err
				}
				_, err = io.Copy(out, tr)
				_ = out.Close()
				if err != nil {
					_ = os.Remove(tmpPath)
					return "", err
				}

				_ = os.Remove(dstPath)
				err = os.Rename(tmpPath, dstPath)
				if err != nil {
					_ = os.Remove(tmpPath)
					return "", err
				}

				return name, nil
			}
		}
	}

	return "", errors.New("can not find api binary in tar.gz")
}
