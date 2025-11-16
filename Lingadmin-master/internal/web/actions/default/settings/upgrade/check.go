package upgrade

import (
	"encoding/json"
	"os"

	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/maps"
)

// CheckAction 检查更新
type CheckAction struct {
	actionutils.ParentAction
}

func (this *CheckAction) RunPost(params struct{}) {
	// 读取更新信息文件
	updateInfoPath := Tea.ConfigFile("update_info.json")
	data, err := os.ReadFile(updateInfoPath)
	if err != nil {
		// 没有检测到新版本
		this.Data["hasUpdate"] = false
		this.Success()
		return
	}

	var updateInfo map[string]interface{}
	err = json.Unmarshal(data, &updateInfo)
	if err != nil {
		this.Fail("解析更新信息失败: " + err.Error())
		return
	}

	this.Data["hasUpdate"] = true
	this.Data["updateInfo"] = maps.Map{
		"version":        updateInfo["version"],
		"currentVersion": updateInfo["currentVersion"],
		"changelog":      updateInfo["changelog"],
		"description":    updateInfo["description"],
		"checkTime":      updateInfo["checkTime"],
		"downloadURL":    updateInfo["downloadURL"],
	}

	this.Success()
}
