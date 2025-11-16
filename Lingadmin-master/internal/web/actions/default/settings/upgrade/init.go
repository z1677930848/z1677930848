package upgrade

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/settings/settingutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth(configloaders.AdminModuleCodeSetting)).
			Helper(settingutils.NewHelper("upgrade")).
			Prefix("/settings/upgrade").
			Get("", new(IndexAction)).
			Post("/check", new(CheckAction)).
			Post("/install", new(InstallAction)).
			Post("/progress", new(ProgressAction)).
			Post("/history", new(HistoryAction)).
			EndAll()
	})
}
