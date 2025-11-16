package dashboard

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction actions.Action

// RunGet 显示用户仪表盘
func (this *IndexAction) RunGet(params struct {
	Auth *helpers.UserShouldAuth
}) {
	// 检查是否登录
	if !params.Auth.IsUser() {
		this.RedirectURL("/user")
		return
	}

	adminId := params.Auth.AdminId()

	this.Data["title"] = "用户仪表盘"
	this.Data["adminId"] = adminId

	// TODO: 从 RPC 获取用户的统计数据
	this.Data["totalDomains"] = 5
	this.Data["totalRequests"] = "12.5K"
	this.Data["totalTraffic"] = "3.2GB"
	this.Data["sslCerts"] = 3

	this.Show()
}
