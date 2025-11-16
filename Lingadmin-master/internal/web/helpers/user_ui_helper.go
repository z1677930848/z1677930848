package helpers

import (
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/iwind/TeaGo/actions"
	"html/template"
)

type UserUIHelper struct {
}

func NewUserUIHelper() *UserUIHelper {
	return &UserUIHelper{}
}

func (this *UserUIHelper) BeforeAction(actionPtr actions.ActionWrapper) (goNext bool) {
	var action = actionPtr.Object()

	// 设置Plus标志
	action.Data["teaIsPlus"] = teaconst.IsPlus

	// 加载用户界面配置
	userUIConfig, err := configloaders.LoadUserUIConfig()
	if err == nil {
		action.Data["userUIConfig"] = userUIConfig
		action.Data["productName"] = userUIConfig.ProductName
		action.Data["userSystemName"] = userUIConfig.UserSystemName
		action.Data["showVersion"] = userUIConfig.ShowVersion

		// 设置版本号
		if userUIConfig.ShowVersion {
			if len(userUIConfig.Version) > 0 {
				action.Data["version"] = userUIConfig.Version
			} else {
				action.Data["version"] = teaconst.Version
			}
		} else {
			action.Data["version"] = ""
		}

		// 设置favicon和logo
		if userUIConfig.FaviconFileId > 0 {
			action.Data["faviconURL"] = fmt.Sprintf("/ui/image/%d", userUIConfig.FaviconFileId)
		}
		if userUIConfig.LogoFileId > 0 {
			action.Data["logoURL"] = fmt.Sprintf("/ui/image/%d", userUIConfig.LogoFileId)
		}

		// 页脚设置 - 将HTML标记为安全类型
		action.Data["showPageFooter"] = userUIConfig.ShowPageFooter
		action.Data["pageFooterHTML"] = template.HTML(userUIConfig.PageFooterHTML)
	} else {
		// 默认配置
		action.Data["productName"] = "LingCDN"
		action.Data["userSystemName"] = "LingCDN管理系统 用户端"
		action.Data["showVersion"] = true
		action.Data["version"] = teaconst.Version
		action.Data["showPageFooter"] = false
		action.Data["pageFooterHTML"] = template.HTML("")
	}

	return true
}
