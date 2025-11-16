package logout

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction actions.Action

// Run 用户退出登录
func (this *IndexAction) Run(params struct {
	Auth *helpers.UserShouldAuth
}) {
	params.Auth.Logout()
	this.RedirectURL("/user")
}
