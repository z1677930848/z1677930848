package upgrade

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/tasks"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/actions"
)

// InstallAction 执行更新
type InstallAction struct {
	actionutils.ParentAction
}

func (this *InstallAction) RunPost(params struct {
	Must *actions.Must
}) {
	// 在后台执行更新
	go func() {
		err := tasks.DownloadAndInstallUpdate()
		if err != nil {
			// 记录错误日志
			// logs.Println("[UPGRADE]update failed:", err)
		}
	}()

	this.Data["message"] = "更新任务已启动，系统将在后台下载并安装新版本"
	this.Success()
}
