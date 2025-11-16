package upgrade

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
)

// HistoryAction 获取升级历史
type HistoryAction struct {
	actionutils.ParentAction
}

func (this *HistoryAction) RunPost(params struct {
	Limit int
}) {
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Limit > 50 {
		params.Limit = 50
	}

	logManager := utils.SharedUpgradeLogManager()
	logs, err := logManager.GetLogs("admin", params.Limit)
	if err != nil {
		this.Fail("获取升级历史失败: " + err.Error())
		return
	}

	var historyList []map[string]interface{}
	for _, log := range logs {
		item := map[string]interface{}{
			"id":            log.ID,
			"oldVersion":    log.OldVersion,
			"newVersion":    log.NewVersion,
			"status":        log.Status,
			"startTime":     log.StartTime.Format("2006-01-02 15:04:05"),
			"downloadSpeed": log.DownloadSpeed,
			"errorMessage":  log.ErrorMessage,
		}

		if !log.EndTime.IsZero() {
			item["endTime"] = log.EndTime.Format("2006-01-02 15:04:05")
			item["duration"] = log.Duration
		}

		historyList = append(historyList, item)
	}

	this.Data["history"] = historyList
	this.Success()
}
