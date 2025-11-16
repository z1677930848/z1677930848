package stats

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction actions.Action

// RunGet 显示统计数据
func (this *IndexAction) RunGet(params struct {
	Auth *helpers.UserShouldAuth
}) {
	if !params.Auth.IsUser() {
		this.RedirectURL("/user")
		return
	}

	adminId := params.Auth.AdminId()

	this.Data["title"] = "统计数据"
	this.Data["adminId"] = adminId

	// TODO: 从 RPC 获取统计数据
	stats := map[string]interface{}{
		"todayRequests":  12543,
		"todayTraffic":   3254789120,
		"weekRequests":   87654,
		"weekTraffic":    22341234567,
		"monthRequests":  345678,
		"monthTraffic":   98765432109,
	}

	this.Data["stats"] = stats
	this.Show()
}
