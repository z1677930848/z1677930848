package profile

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction actions.Action

// RunGet 显示个人设置
func (this *IndexAction) RunGet(params struct {
	Auth *helpers.UserShouldAuth
}) {
	if !params.Auth.IsUser() {
		this.RedirectURL("/user")
		return
	}

	adminId := params.Auth.AdminId()

	this.Data["title"] = "个人设置"
	this.Data["adminId"] = adminId

	// TODO: 从 RPC 获取用户信息

	this.Show()
}

// RunPost 更新个人信息
func (this *IndexAction) RunPost(params struct {
	Email    string
	Phone    string
	Password string

	Auth *helpers.UserShouldAuth
}) {
	if !params.Auth.IsUser() {
		this.Fail("请先登录")
		return
	}

	// TODO: 调用 RPC 更新用户信息

	this.Success()
}
