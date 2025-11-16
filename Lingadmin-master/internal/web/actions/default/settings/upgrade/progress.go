package upgrade

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
)

// ProgressAction 获取更新进度
type ProgressAction struct {
	actionutils.ParentAction
}

func (this *ProgressAction) RunPost(params struct{}) {
	logManager := utils.SharedUpgradeLogManager()

	// 获取最新的升级日志
	latestLog, err := logManager.GetLatestLog("admin")
	if err != nil || latestLog == nil {
		this.Data["status"] = "idle"
		this.Success()
		return
	}

	this.Data["status"] = string(latestLog.Status)
	this.Data["progress"] = map[string]interface{}{
		"oldVersion":    latestLog.OldVersion,
		"newVersion":    latestLog.NewVersion,
		"status":        latestLog.Status,
		"startTime":     latestLog.StartTime.Format("2006-01-02 15:04:05"),
		"downloadSpeed": latestLog.DownloadSpeed,
		"downloadSize":  latestLog.DownloadSize,
		"errorMessage":  latestLog.ErrorMessage,
		"errorStage":    latestLog.ErrorStage,
	}

	if !latestLog.EndTime.IsZero() {
		this.Data["progress"].(map[string]interface{})["endTime"] = latestLog.EndTime.Format("2006-01-02 15:04:05")
		this.Data["progress"].(map[string]interface{})["duration"] = latestLog.Duration
	}

	this.Success()
}
