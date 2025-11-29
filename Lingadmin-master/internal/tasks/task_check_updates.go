// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package tasks

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
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
	"github.com/TeaOSLab/EdgeAdmin/internal/events"
	"github.com/TeaOSLab/EdgeAdmin/internal/goman"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/systemconfigs"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	stringutil "github.com/iwind/TeaGo/utils/string"
)

func init() {
	events.On(events.EventStart, func() {
		var task = NewCheckUpdatesTask()
		goman.New(func() {
			task.Start()
		})

		// 鍚姩涓存椂鏂囦欢娓呯悊浠诲姟
		utils.ScheduleCleanupTask()
	})
}

type CheckUpdatesTask struct {
	ticker     *time.Ticker
	logManager *utils.UpgradeLogManager
	notifier   utils.UpdateNotifier
	cleaner    *utils.TempFileCleaner
}

func NewCheckUpdatesTask() *CheckUpdatesTask {
	// 鍒涘缓澶氶€氶亾閫氱煡鍣?
	multiNotifier := utils.NewMultiNotifier()
	multiNotifier.AddNotifier(utils.NewLogNotifier())
	// 鍙牴鎹厤缃坊鍔犳洿澶氶€氱煡鍣?
	// multiNotifier.AddNotifier(utils.NewWebhookNotifier("http://your-webhook-url"))

	return &CheckUpdatesTask{
		logManager: utils.SharedUpgradeLogManager(),
		notifier:   multiNotifier,
		cleaner:    utils.NewTempFileCleaner(),
	}
}

func (this *CheckUpdatesTask) Start() {
	// 鍚姩鍚庣珛鍗虫鏌ヤ竴娆?
	err := this.Loop()
	if err != nil {
		logs.Println("[TASK][CHECK_UPDATES_TASK]" + err.Error())
	}

	// 鐒跺悗姣?灏忔椂妫€鏌ヤ竴娆?
	this.ticker = time.NewTicker(6 * time.Hour)
	for range this.ticker.C {
		err := this.Loop()
		if err != nil {
			logs.Println("[TASK][CHECK_UPDATES_TASK]" + err.Error())
		}
	}
}

func (this *CheckUpdatesTask) Loop() error {
	// 妫€鏌ユ槸鍚﹀紑鍚?
	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		return err
	}
	valueResp, err := rpcClient.SysSettingRPC().ReadSysSetting(rpcClient.Context(0), &pb.ReadSysSettingRequest{Code: systemconfigs.SettingCodeCheckUpdates})
	if err != nil {
		return err
	}
	var valueJSON = valueResp.ValueJSON
	var config = &systemconfigs.CheckUpdatesConfig{AutoCheck: false}
	if len(valueJSON) > 0 {
		err = json.Unmarshal(valueJSON, config)
		if err != nil {
			return errors.New("decode config failed: " + err.Error())
		}
		if !config.AutoCheck {
			return nil
		}
	} else {
		return nil
	}

	// 寮€濮嬫鏌?
	type Response struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}

	// 鐩墠鏀寔Linux
	if runtime.GOOS != "linux" {
		return nil
	}

	var apiURL = teaconst.UpdatesURL
	apiURL = strings.ReplaceAll(apiURL, "${os}", runtime.GOOS)
	apiURL = strings.ReplaceAll(apiURL, "${arch}", runtime.GOARCH)

	logs.Println("[TASK][CHECK_UPDATES_TASK]checking updates from:", apiURL)

	resp, err := http.Get(apiURL)
	if err != nil {
		return errors.New("read api failed: " + err.Error())
	}

	defer func() {
		_ = resp.Body.Close()
	}()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New("read api failed: " + err.Error())
	}

	var apiResponse = &Response{}
	err = json.Unmarshal(data, apiResponse)
	if err != nil {
		return errors.New("decode version data failed: " + err.Error())
	}

	if apiResponse.Code != 200 {
		return errors.New("invalid response: " + apiResponse.Message)
	}

	var m = maps.NewMap(apiResponse.Data)
	var dlHost = m.GetString("host")
	var versions = m.GetSlice("versions")
	if len(versions) > 0 {
		for _, version := range versions {
			var vMap = maps.NewMap(version)
			if vMap.GetString("code") == "admin" {
				var latestVersion = vMap.GetString("version")
				var changelog = vMap.GetString("changelog")
				var description = vMap.GetString("description")
				var downloadURL = dlHost + vMap.GetString("url")
				var fileSHA256 = vMap.GetString("sha256")

				logs.Println("[TASK][CHECK_UPDATES_TASK]current version:", teaconst.Version, "latest version:", latestVersion)

				if stringutil.VersionCompare(teaconst.Version, latestVersion) < 0 {
					teaconst.NewVersionCode = latestVersion
					teaconst.NewVersionDownloadURL = downloadURL

					// 淇濆瓨鏇存柊淇℃伅鍒版枃浠?
					updateInfo := map[string]interface{}{
						"version":        latestVersion,
						"currentVersion": teaconst.Version,
						"downloadURL":    downloadURL,
						"changelog":      changelog,
						"description":    description,
						"sha256":         fileSHA256,
						"checkTime":      time.Now().Format("2006-01-02 15:04:05"),
					}
					updateInfoJSON, _ := json.MarshalIndent(updateInfo, "", "  ")
					_ = os.WriteFile(Tea.ConfigFile("update_info.json"), updateInfoJSON, 0644)

					logs.Println("[TASK][CHECK_UPDATES_TASK]new version available:", latestVersion)
					logs.Println("[TASK][CHECK_UPDATES_TASK]download url:", downloadURL)
					logs.Println("[TASK][CHECK_UPDATES_TASK]changelog:", changelog)

					return nil
				} else {
					logs.Println("[TASK][CHECK_UPDATES_TASK]no updates available, current version is latest")
					teaconst.NewVersionCode = ""
					teaconst.NewVersionDownloadURL = ""
				}
			}
		}
	}

	return nil
}

// DownloadAndInstallUpdate 涓嬭浇骞跺畨瑁呮洿鏂帮紙鏀硅繘鐗堬級
func DownloadAndInstallUpdate() error {
	startTime := time.Now()
	logs.Println("[UPDATE]starting update process...")

	// 鍒涘缓鍗囩骇鏃ュ織
	logManager := utils.SharedUpgradeLogManager()
	upgradeLog := &utils.UpgradeLog{
		Component:  "admin",
		OldVersion: teaconst.Version,
		Status:     utils.StatusPending,
		StartTime:  startTime,
	}
	_ = logManager.CreateLog(upgradeLog)

	// 鍒涘缓涓存椂鏂囦欢娓呯悊鍣?
	cleaner := utils.NewTempFileCleaner()
	defer func() {
		if err := cleaner.CleanupAll(); err != nil {
			logs.Println("[UPDATE]cleanup failed:", err)
		}
	}()

	// 鍒涘缓閫氱煡鍣?
	notifier := utils.NewMultiNotifier()
	notifier.AddNotifier(utils.NewLogNotifier())
	notifier.AddNotifier(utils.NewConsoleNotifier())

	// 璇诲彇鏇存柊淇℃伅
	updateInfoData, err := os.ReadFile(Tea.ConfigFile("update_info.json"))
	if err != nil {
		upgradeErr := utils.NewUpgradeError(utils.StageCheckVersion, utils.ErrCodeNetworkFailed,
			"read update info failed", err)
		handleUpgradeError(upgradeLog, logManager, notifier, upgradeErr)
		return upgradeErr
	}

	var updateInfo map[string]interface{}
	err = json.Unmarshal(updateInfoData, &updateInfo)
	if err != nil {
		upgradeErr := utils.NewUpgradeError(utils.StageCheckVersion, utils.ErrCodeInvalidResponse,
			"parse update info failed", err)
		handleUpgradeError(upgradeLog, logManager, notifier, upgradeErr)
		return upgradeErr
	}

	downloadURL := updateInfo["downloadURL"].(string)
	expectedSHA256 := updateInfo["sha256"].(string)
	version := updateInfo["version"].(string)

	upgradeLog.NewVersion = version
	upgradeLog.DownloadURL = downloadURL
	_ = logManager.UpdateLog(upgradeLog)

	// 閫氱煡寮€濮?
	notifier.NotifyStart("admin", version)

	logs.Println("[UPDATE]downloading version:", version)
	logs.Println("[UPDATE]download url:", downloadURL)

	// 鍒涘缓涓存椂鐩綍
	tmpDir := Tea.ConfigFile("tmp")
	_ = os.MkdirAll(tmpDir, 0755)

	// 涓嬭浇鏂囦欢
	downloadFilePath := filepath.Join(tmpDir, fmt.Sprintf("ling-admin-v%s.zip", version))
	cleaner.AddFile(downloadFilePath)

	upgradeLog.Status = utils.StatusDownloading
	_ = logManager.UpdateLog(upgradeLog)

	err = downloadFileWithProgress(downloadURL, downloadFilePath, func(progress float32, speed float64) {
		message := fmt.Sprintf("Downloading: %.1f MB/s", speed)
		notifier.NotifyProgress("admin", progress*0.6, message)
	})
	if err != nil {
		upgradeErr := utils.NewUpgradeError(utils.StageDownload, utils.ErrCodeDownloadFailed,
			"download failed", err)
		handleUpgradeError(upgradeLog, logManager, notifier, upgradeErr)
		return upgradeErr
	}

	stat, _ := os.Stat(downloadFilePath)
	upgradeLog.DownloadSize = stat.Size()
	_ = logManager.UpdateLog(upgradeLog)

	logs.Println("[UPDATE]download completed")
	notifier.NotifyProgress("admin", 0.65, "Download complete, verifying...")

	// 楠岃瘉SHA256
	upgradeLog.Status = utils.StatusVerifying
	_ = logManager.UpdateLog(upgradeLog)

	if expectedSHA256 != "" {
		actualSHA256, err := calculateSHA256(downloadFilePath)
		if err != nil {
			upgradeErr := utils.NewUpgradeError(utils.StageVerify, utils.ErrCodeVerifyFailed,
				"calculate sha256 failed", err)
			handleUpgradeError(upgradeLog, logManager, notifier, upgradeErr)
			return upgradeErr
		}
		if actualSHA256 != expectedSHA256 {
			upgradeErr := utils.NewUpgradeError(utils.StageVerify, utils.ErrCodeVerifyFailed,
				fmt.Sprintf("sha256 mismatch: expected %s, got %s", expectedSHA256, actualSHA256), nil)
			handleUpgradeError(upgradeLog, logManager, notifier, upgradeErr)
			return upgradeErr
		}
		logs.Println("[UPDATE]sha256 verification passed")
		notifier.NotifyProgress("admin", 0.7, "Verification passed")
	}

	// 瑙ｅ帇鏂囦欢
	extractDir := filepath.Join(tmpDir, "extract")
	_ = os.RemoveAll(extractDir)
	_ = os.MkdirAll(extractDir, 0755)
	cleaner.AddDir(extractDir)

	notifier.NotifyProgress("admin", 0.75, "Extracting files...")
	err = unzip(downloadFilePath, extractDir)
	if err != nil {
		upgradeErr := utils.NewUpgradeError(utils.StageUnzip, utils.ErrCodeUnzipFailed,
			"unzip failed", err)
		handleUpgradeError(upgradeLog, logManager, notifier, upgradeErr)
		return upgradeErr
	}

	logs.Println("[UPDATE]extract completed")
	notifier.NotifyProgress("admin", 0.85, "Installing...")

	// 鎵惧埌浜岃繘鍒舵枃浠?
	binaryPath := filepath.Join(extractDir, "ling-admin")
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		upgradeErr := utils.NewUpgradeError(utils.StageInstall, utils.ErrCodeInstallFailed,
			"binary file not found in package", nil)
		handleUpgradeError(upgradeLog, logManager, notifier, upgradeErr)
		return upgradeErr
	}

	upgradeLog.Status = utils.StatusInstalling
	_ = logManager.UpdateLog(upgradeLog)

	// 澶囦唤褰撳墠鐗堟湰
	currentBinary, err := os.Executable()
	if err != nil {
		upgradeErr := utils.NewUpgradeError(utils.StageBackup, utils.ErrCodeBackupFailed,
			"get current binary path failed", err)
		handleUpgradeError(upgradeLog, logManager, notifier, upgradeErr)
		return upgradeErr
	}

	backupPath := currentBinary + ".backup." + teaconst.Version
	err = copyFile(currentBinary, backupPath)
	if err != nil {
		logs.Println("[UPDATE]backup failed:", err.Error(), "- continuing anyway")
	} else {
		logs.Println("[UPDATE]current version backed up to:", backupPath)
		upgradeLog.BackupPath = backupPath
		_ = logManager.UpdateLog(upgradeLog)
		// 澶囦唤鏂囦欢淇濈暀7澶╁悗娓呯悊
		cleaner.AddFileWithDelay(backupPath, 7*24*time.Hour)
	}

	notifier.NotifyProgress("admin", 0.9, "Replacing binary...")

	// 鏇挎崲浜岃繘鍒舵枃浠?
	err = os.Chmod(binaryPath, 0755)
	if err != nil {
		upgradeErr := utils.NewUpgradeError(utils.StageInstall, utils.ErrCodePermissionDenied,
			"chmod new binary failed", err)
		handleUpgradeError(upgradeLog, logManager, notifier, upgradeErr)
		return upgradeErr
	}

	// 鍏堝皾璇曠洿鎺ヨ鐩?
	err = copyFile(binaryPath, currentBinary)
	if err != nil {
		upgradeErr := utils.NewUpgradeError(utils.StageInstall, utils.ErrCodeInstallFailed,
			"replace binary failed", err)
		handleUpgradeError(upgradeLog, logManager, notifier, upgradeErr)
		return upgradeErr
	}

	logs.Println("[UPDATE]binary updated successfully")

	// 鏇存柊web鐩綍锛堝鏋滃瓨鍦級
	webSrcDir := filepath.Join(extractDir, "web")
	if _, err := os.Stat(webSrcDir); err == nil {
		webDestDir := Tea.Root + "/web"
		_ = os.RemoveAll(webDestDir)
		err = copyDir(webSrcDir, webDestDir)
		if err != nil {
			logs.Println("[UPDATE]web update failed:", err.Error())
		} else {
			logs.Println("[UPDATE]web directory updated")
		}
	}

	notifier.NotifyProgress("admin", 0.95, "Update complete, restarting...")

	// 鏇存柊鎴愬姛
	duration := time.Since(startTime)
	upgradeLog.Status = utils.StatusSuccess
	upgradeLog.EndTime = time.Now()
	upgradeLog.DownloadSpeed = float64(upgradeLog.DownloadSize) / duration.Seconds() / 1024 / 1024
	_ = logManager.UpdateLog(upgradeLog)

	logs.Println("[UPDATE]update completed successfully, version:", version)
	logs.Println("[UPDATE]restarting service...")

	notifier.NotifySuccess("admin", version, duration)

	// 閲嶅惎鏈嶅姟
	return restartService()
}

// handleUpgradeError 澶勭悊鍗囩骇閿欒
func handleUpgradeError(log *utils.UpgradeLog, logManager *utils.UpgradeLogManager,
	notifier utils.UpdateNotifier, err *utils.UpgradeError) {
	log.Status = utils.StatusFailed
	log.ErrorCode = int(err.Code)
	log.ErrorMessage = err.Message
	log.ErrorStage = string(err.Stage)
	log.EndTime = time.Now()
	_ = logManager.UpdateLog(log)

	notifier.NotifyFailed(log.Component, log.NewVersion, err)
}

// downloadFileWithProgress 涓嬭浇鏂囦欢骞舵樉绀鸿繘搴?
func downloadFileWithProgress(url, dest string, progressCallback func(progress float32, speed float64)) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("http status: %d", resp.StatusCode)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	contentLength := resp.ContentLength
	downloaded := int64(0)
	startTime := time.Now()
	lastNotifyTime := startTime

	buffer := make([]byte, 32*1024)
	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			_, writeErr := out.Write(buffer[:n])
			if writeErr != nil {
				return writeErr
			}

			downloaded += int64(n)

			// 姣忕鏇存柊涓€娆¤繘搴?
			if time.Since(lastNotifyTime) >= time.Second {
				progress := float32(downloaded) / float32(contentLength)
				speed := float64(downloaded) / time.Since(startTime).Seconds() / 1024 / 1024
				progressCallback(progress, speed)
				lastNotifyTime = time.Now()
			}
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}

	// 鏈€鍚庝竴娆¤繘搴︽洿鏂?
	if contentLength > 0 {
		speed := float64(downloaded) / time.Since(startTime).Seconds() / 1024 / 1024
		progressCallback(1.0, speed)
	}

	return nil
}

func calculateSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

func copyFile(src, dst string) error {
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

	// 鍚屾鍒扮鐩?
	err = destFile.Sync()
	return err
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		return copyFile(path, dstPath)
	})
}

func restartService() error {
	// 尝试使用systemctl重启
	cmd := exec.Command("systemctl", "restart", teaconst.SystemdServiceName)
	err := cmd.Run()
	if err == nil {
		return nil
	}

	// 如果systemctl失败，则直接拉起新进程并退出当前进程
	logs.Println("[UPDATE]systemctl restart failed, trying direct restart")

	exePath, pathErr := os.Executable()
	if pathErr != nil {
		return fmt.Errorf("locate executable failed: %w", pathErr)
	}

	newCmd := exec.Command(exePath, os.Args[1:]...)
	newCmd.Stdout = os.Stdout
	newCmd.Stderr = os.Stderr
	newCmd.Env = os.Environ()
	if startErr := newCmd.Start(); startErr != nil {
		return fmt.Errorf("start new process failed: %w", startErr)
	}

	// 延迟退出，保证当前请求完成
	time.AfterFunc(1*time.Second, func() {
		os.Exit(0)
	})

	return nil
}
