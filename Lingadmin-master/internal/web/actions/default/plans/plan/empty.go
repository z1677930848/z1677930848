package plan

import "github.com/iwind/TeaGo/actions"

// EmptyAction 预留占位，避免路由404
type EmptyAction actions.Action

func (this *EmptyAction) RunGet(params struct{}) {
	this.RedirectURL("/plans")
}
