// Copyright 2022 Liuxiangchao iwind.liu@gmail.com. All rights reserved. Official site: https://goedge.cn .

package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	stringutil "github.com/iwind/TeaGo/utils/string"
)

type UpgradeFileWriter struct {
	rawWriter io.Writer
	written   int64
}

func NewUpgradeFileWriter(rawWriter io.Writer) *UpgradeFileWriter {
	return &UpgradeFileWriter{rawWriter: rawWriter}
}

func (this *UpgradeFileWriter) Write(p []byte) (n int, err error) {
	n, err = this.rawWriter.Write(p)
	this.written += int64(n)
	return
}

func (this *UpgradeFileWriter) TotalWritten() int64 {
	return this.written
}

type UpgradeManager struct {
	client *http.Client

	component string

	newVersion    string
	contentLength int64
	isDownloading bool
	writer        *UpgradeFileWriter
	body          io.ReadCloser
	isCancelled   bool

	// 新增字段
	notifier       UpdateNotifier
	log            *UpgradeLog
	logManager     *UpgradeLogManager
	startTime      time.Time
	downloadedSize int64
	destFile       string // 目标文件路径
	resumeSupport  bool   // 是否支持断点续传
}

func NewUpgradeManager(component string) *UpgradeManager {
	notifier := NewMultiNotifier()
	notifier.AddNotifier(NewLogNotifier())
	// 可以根据配置添加更多通知器
	// notifier.AddNotifier(NewWebhookNotifier("http://your-webhook-url"))

	return &UpgradeManager{
		component: component,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, // TODO: 应该移除，使用正确的证书验证
				},
			},
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       30 * time.Minute,
		},
		notifier:   notifier,
		logManager: SharedUpgradeLogManager(),
	}
}

// SetNotifier 设置通知器
func (this *UpgradeManager) SetNotifier(notifier UpdateNotifier) {
	this.notifier = notifier
}

func (this *UpgradeManager) Start() error {
	if this.isDownloading {
		return NewUpgradeError(StageDownload, ErrCodeAlreadyRunning, "another process is running", nil)
	}

	this.isDownloading = true
	this.startTime = time.Now()

	// 创建升级日志
	this.log = &UpgradeLog{
		Component:  this.component,
		OldVersion: teaconst.Version,
		Status:     StatusPending,
		StartTime:  this.startTime,
	}
	_ = this.logManager.CreateLog(this.log)

	defer func() {
		this.client.CloseIdleConnections()
		this.isDownloading = false

		// 更新日志
		this.log.EndTime = time.Now()
		_ = this.logManager.UpdateLog(this.log)
	}()

	// 检查unzip
	unzipExe, _ := exec.LookPath("unzip")
	if len(unzipExe) == 0 {
		err := NewUpgradeError(StageInstall, ErrCodeInstallFailed, "can not find 'unzip' command", nil)
		this.handleError(err)
		return err
	}

	// 检查cp
	cpExe, _ := exec.LookPath("cp")
	if len(cpExe) == 0 {
		err := NewUpgradeError(StageInstall, ErrCodeInstallFailed, "can not find 'cp' command", nil)
		this.handleError(err)
		return err
	}

	// 检查新版本
	downloadURL, err := this.checkNewVersion()
	if err != nil {
		this.handleError(err)
		return err
	}

	// 通知开始
	this.notifier.NotifyStart(this.component, this.newVersion)
	this.log.NewVersion = this.newVersion
	this.log.Status = StatusDownloading
	_ = this.logManager.UpdateLog(this.log)

	// 下载和安装
	err = this.downloadAndInstall(downloadURL, unzipExe, cpExe)
	if err != nil {
		this.handleError(err)
		return err
	}

	// 成功
	duration := time.Since(this.startTime)
	this.notifier.NotifySuccess(this.component, this.newVersion, duration)
	this.log.Status = StatusSuccess
	this.log.DownloadSpeed = float64(this.log.DownloadSize) / duration.Seconds() / 1024 / 1024 // MB/s
	_ = this.logManager.UpdateLog(this.log)

	return nil
}

// checkNewVersion 检查新版本
func (this *UpgradeManager) checkNewVersion() (string, error) {
	this.notifier.NotifyProgress(this.component, 0.0, "Checking for new version...")

	var url = teaconst.UpdatesURL
	var osName = runtime.GOOS
	if Tea.IsTesting() && osName == "darwin" {
		osName = "linux"
	}
	url = strings.ReplaceAll(url, "${os}", osName)
	url = strings.ReplaceAll(url, "${arch}", runtime.GOARCH)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", NewUpgradeError(StageCheckVersion, ErrCodeNetworkFailed, "create url request failed", err)
	}

	resp, err := this.client.Do(req)
	if err != nil {
		return "", NewUpgradeError(StageCheckVersion, ErrCodeNetworkFailed, "read latest version failed", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", NewUpgradeError(StageCheckVersion, ErrCodeInvalidResponse,
			fmt.Sprintf("invalid response code '%d'", resp.StatusCode), nil)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", NewUpgradeError(StageCheckVersion, ErrCodeNetworkFailed, "read response body failed", err)
	}

	var m = maps.Map{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		return "", NewUpgradeError(StageCheckVersion, ErrCodeInvalidResponse,
			"invalid response data: "+string(data), err)
	}

	var code = m.GetInt("code")
	if code != 200 {
		return "", NewUpgradeError(StageCheckVersion, ErrCodeInvalidResponse, m.GetString("message"), nil)
	}

	var dataMap = m.GetMap("data")
	var downloadHost = dataMap.GetString("host")
	var versions = dataMap.GetSlice("versions")
	var downloadPath = ""
	for _, component := range versions {
		var componentMap = maps.NewMap(component)
		if componentMap.Has("version") {
			if componentMap.GetString("code") == this.component {
				var version = componentMap.GetString("version")
				if stringutil.VersionCompare(version, teaconst.Version) > 0 {
					this.newVersion = version
					downloadPath = componentMap.GetString("url")
					break
				}
			}
		}
	}

	if len(downloadPath) == 0 {
		return "", NewUpgradeError(StageCheckVersion, ErrCodeNoNewVersion, "no latest version to download", nil)
	}

	downloadURL := downloadHost + downloadPath
	this.log.DownloadURL = downloadURL
	_ = this.logManager.UpdateLog(this.log)

	return downloadURL, nil
}

// downloadAndInstall 下载并安装
func (this *UpgradeManager) downloadAndInstall(downloadURL, unzipExe, cpExe string) error {
	// 准备下载路径
	var tmpDir = os.TempDir()
	var filename = filepath.Base(downloadURL)
	this.destFile = filepath.Join(tmpDir, filename)

	// 创建临时文件清理器
	cleaner := NewTempFileCleaner()
	cleaner.AddFile(this.destFile)
	defer func() {
		if err := cleaner.CleanupAll(); err != nil {
			logs.Println("[UPGRADE]failed to cleanup temp files:", err)
		}
	}()

	// 检查是否支持断点续传
	existingSize := int64(0)
	if stat, err := os.Stat(this.destFile); err == nil {
		existingSize = stat.Size()
	}

	// 下载文件
	err := this.downloadFile(downloadURL, existingSize)
	if err != nil {
		return err
	}

	this.log.DownloadSize = this.downloadedSize
	_ = this.logManager.UpdateLog(this.log)

	// 解压
	this.notifier.NotifyProgress(this.component, 0.8, "Extracting files...")
	this.log.Status = StatusVerifying
	_ = this.logManager.UpdateLog(this.log)

	unzipDir, err := this.unzipFile(unzipExe)
	if err != nil {
		return err
	}
	cleaner.AddDir(unzipDir)

	// 查找安装文件
	installationFiles, err := filepath.Glob(unzipDir + "/edge-" + this.component + "/*")
	if err != nil {
		return NewUpgradeError(StageInstall, ErrCodeInstallFailed, "lookup installation files failed", err)
	}

	// 备份并安装
	this.notifier.NotifyProgress(this.component, 0.9, "Installing files...")
	this.log.Status = StatusInstalling
	_ = this.logManager.UpdateLog(this.log)

	err = this.backupAndInstall(installationFiles, cpExe)
	if err != nil {
		return err
	}

	this.notifier.NotifyProgress(this.component, 1.0, "Installation complete")

	return nil
}

// downloadFile 下载文件（支持断点续传）
func (this *UpgradeManager) downloadFile(downloadURL string, existingSize int64) error {
	req, err := http.NewRequest(http.MethodGet, downloadURL, nil)
	if err != nil {
		return NewUpgradeError(StageDownload, ErrCodeNetworkFailed, "create download request failed", err)
	}

	// 如果文件已存在，尝试断点续传
	if existingSize > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", existingSize))
		logs.Println("[UPGRADE]resuming download from:", existingSize)
	}

	resp, err := this.client.Do(req)
	if err != nil {
		return NewUpgradeError(StageDownload, ErrCodeDownloadFailed, "download failed: "+downloadURL, err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	// 检查响应状态
	if resp.StatusCode == http.StatusPartialContent {
		// 支持断点续传
		this.resumeSupport = true
		logs.Println("[UPGRADE]resume supported, continuing download")
	} else if resp.StatusCode == http.StatusOK {
		// 不支持断点续传，重新下载
		this.resumeSupport = false
		existingSize = 0
		_ = os.Remove(this.destFile)
		logs.Println("[UPGRADE]resume not supported, starting from beginning")
	} else {
		return NewUpgradeError(StageDownload, ErrCodeDownloadFailed,
			fmt.Sprintf("invalid response code '%d'", resp.StatusCode), nil)
	}

	this.contentLength = resp.ContentLength + existingSize
	this.body = resp.Body

	// 打开文件（追加或创建）
	var openFlags = os.O_CREATE | os.O_WRONLY
	if existingSize > 0 && this.resumeSupport {
		openFlags |= os.O_APPEND
	} else {
		openFlags |= os.O_TRUNC
	}

	fp, err := os.OpenFile(this.destFile, openFlags, 0644)
	if err != nil {
		return NewUpgradeError(StageDownload, ErrCodeDownloadFailed, "create file failed", err)
	}
	defer fp.Close()

	this.writer = NewUpgradeFileWriter(fp)
	this.downloadedSize = existingSize

	// 下载并实时更新进度
	buffer := make([]byte, 32*1024) // 32KB buffer
	lastNotifyTime := time.Now()

	for {
		if this.isCancelled {
			return NewUpgradeError(StageDownload, ErrCodeCancelled, "download cancelled", nil)
		}

		n, err := resp.Body.Read(buffer)
		if n > 0 {
			_, writeErr := this.writer.Write(buffer[:n])
			if writeErr != nil {
				return NewUpgradeError(StageDownload, ErrCodeDownloadFailed, "write file failed", writeErr)
			}

			this.downloadedSize += int64(n)

			// 每秒更新一次进度
			if time.Since(lastNotifyTime) >= time.Second {
				progress := this.Progress()
				if progress > 0 {
					speed := float64(this.downloadedSize-existingSize) / time.Since(this.startTime).Seconds() / 1024 / 1024
					message := fmt.Sprintf("Downloading: %.1f MB/s", speed)
					this.notifier.NotifyProgress(this.component, progress*0.7, message) // 下载占70%进度
				}
				lastNotifyTime = time.Now()
			}
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			return NewUpgradeError(StageDownload, ErrCodeDownloadFailed, "download failed", err)
		}
	}

	return nil
}

// unzipFile 解压文件
func (this *UpgradeManager) unzipFile(unzipExe string) (string, error) {
	var tmpDir = os.TempDir()
	var unzipDir = tmpDir + "/edge-" + this.component + "-tmp"

	stat, err := os.Stat(unzipDir)
	if err == nil && stat.IsDir() {
		err = os.RemoveAll(unzipDir)
		if err != nil {
			return "", NewUpgradeError(StageUnzip, ErrCodeUnzipFailed,
				"remove old dir '"+unzipDir+"' failed", err)
		}
	}

	var unzipCmd = exec.Command(unzipExe, "-q", "-o", this.destFile, "-d", unzipDir)
	var unzipStderr = &bytes.Buffer{}
	unzipCmd.Stderr = unzipStderr
	err = unzipCmd.Run()
	if err != nil {
		return "", NewUpgradeError(StageUnzip, ErrCodeUnzipFailed,
			"unzip installation file failed: "+unzipStderr.String(), err)
	}

	return unzipDir, nil
}

// backupAndInstall 备份并安装
func (this *UpgradeManager) backupAndInstall(installationFiles []string, cpExe string) error {
	currentExe, err := os.Executable()
	if err != nil {
		return NewUpgradeError(StageInstall, ErrCodeInstallFailed,
			"reveal current executable file path failed", err)
	}

	// 备份当前版本
	backupPath := currentExe + ".backup." + teaconst.Version
	err = copyFileSecure(currentExe, backupPath)
	if err != nil {
		logs.Println("[UPGRADE]backup failed:", err, "- continuing anyway")
	} else {
		logs.Println("[UPGRADE]backup created:", backupPath)
		this.log.BackupPath = backupPath
		_ = this.logManager.UpdateLog(this.log)
	}

	var targetDir = filepath.Dir(filepath.Dir(currentExe))
	if !Tea.IsTesting() {
		for _, installationFile := range installationFiles {
			var cpCmd = exec.Command(cpExe, "-R", "-f", installationFile, targetDir)
			var cpStderr = &bytes.Buffer{}
			cpCmd.Stderr = cpStderr
			err = cpCmd.Run()
			if err != nil {
				return NewUpgradeError(StageInstall, ErrCodeInstallFailed,
					"overwrite installation files failed: '"+cpCmd.String()+"': "+cpStderr.String(), err)
			}
		}
	}

	return nil
}

// handleError 处理错误
func (this *UpgradeManager) handleError(err error) {
	if upgradeErr, ok := err.(*UpgradeError); ok {
		this.log.Status = StatusFailed
		this.log.ErrorCode = int(upgradeErr.Code)
		this.log.ErrorMessage = upgradeErr.Message
		this.log.ErrorStage = string(upgradeErr.Stage)
		_ = this.logManager.UpdateLog(this.log)

		this.notifier.NotifyFailed(this.component, this.newVersion, err)
	} else {
		this.log.Status = StatusFailed
		this.log.ErrorMessage = err.Error()
		_ = this.logManager.UpdateLog(this.log)

		this.notifier.NotifyFailed(this.component, this.newVersion, err)
	}
}

func (this *UpgradeManager) IsDownloading() bool {
	return this.isDownloading
}

func (this *UpgradeManager) Progress() float32 {
	if this.contentLength <= 0 {
		return -1
	}
	if this.downloadedSize <= 0 {
		return -1
	}
	return float32(this.downloadedSize) / float32(this.contentLength)
}

func (this *UpgradeManager) NewVersion() string {
	return this.newVersion
}

func (this *UpgradeManager) Cancel() error {
	this.isCancelled = true
	this.isDownloading = false

	if this.body != nil {
		_ = this.body.Close()
	}

	this.log.Status = StatusCancelled
	_ = this.logManager.UpdateLog(this.log)

	this.notifier.NotifyCancelled(this.component, this.newVersion)

	return nil
}

// copyFileSecure 安全地复制文件
func copyFileSecure(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// 同步到磁盘
	err = destFile.Sync()
	if err != nil {
		return err
	}

	return nil
}
