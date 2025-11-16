package domains

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction actions.Action

// RunGet 显示域名列表
func (this *IndexAction) RunGet(params struct {
	Auth *helpers.UserShouldAuth
}) {
	// 检查是否登录
	if !params.Auth.IsUser() {
		this.RedirectURL("/user")
		return
	}

	adminId := params.Auth.AdminId()

	this.Data["title"] = "域名管理"
	this.Data["adminId"] = adminId

	// TODO: 从 RPC 获取用户的域名列表
	domains := []map[string]interface{}{
		{
			"id":     1,
			"name":   "www.example.com",
			"status": "active",
			"ssl":    true,
		},
		{
			"id":     2,
			"name":   "api.example.com",
			"status": "active",
			"ssl":    true,
		},
	}

	this.Data["domains"] = domains
	this.Show()
}

type CreateAction actions.Action

func (this *CreateAction) RunGet(params struct {
	Auth *helpers.UserShouldAuth
}) {
	if !params.Auth.IsUser() {
		this.RedirectURL("/user")
		return
	}

	this.Data["title"] = "添加域名"
	this.Show()
}

func (this *CreateAction) RunPost(params struct {
	Domain string

	Must *actions.Must
	Auth *helpers.UserShouldAuth
}) {
	if !params.Auth.IsUser() {
		this.Fail("请先登录")
		return
	}

	params.Must.
		Field("domain", params.Domain).
		Require("请输入域名")

	// TODO: 调用 RPC 创建域名

	this.Success()
}

type DeleteAction actions.Action

func (this *DeleteAction) RunPost(params struct {
	DomainId int64

	Must *actions.Must
	Auth *helpers.UserShouldAuth
}) {
	if !params.Auth.IsUser() {
		this.Fail("请先登录")
		return
	}

	params.Must.
		Field("domainId", params.DomainId).
		Require("请选择要删除的域名")

	// TODO: 调用 RPC 删除域名

	this.Success()
}
