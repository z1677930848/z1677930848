package user

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/user/dashboard"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/user/domains"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/user/index"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/user/logout"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/user/plan"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/user/profile"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/user/stats"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		// 用户登录（无需认证）
		server.Prefix("/user").
			Data("teaMenu", "user").
			GetPost("", new(index.IndexAction)).
			GetPost("/index", new(index.IndexAction)).
			EndAll()

		// 用户端功能（需要认证）
		server.Prefix("/user").
			Helper(helpers.NewUserUIHelper()).
			Data("teaMenu", "user").
			GetPost("/dashboard", new(dashboard.IndexAction)).
			GetPost("/domains", new(domains.IndexAction)).
			GetPost("/domains/create", new(domains.CreateAction)).
			Post("/domains/delete", new(domains.DeleteAction)).
			GetPost("/stats", new(stats.IndexAction)).
			GetPost("/plan", new(plan.IndexAction)).
			GetPost("/profile", new(profile.IndexAction)).
			Post("/profile/update", new(profile.UpdateAction)).
			Post("/profile/changePassword", new(profile.ChangePasswordAction)).
			Get("/logout", new(logout.IndexAction)).
			EndAll()
	})
}
