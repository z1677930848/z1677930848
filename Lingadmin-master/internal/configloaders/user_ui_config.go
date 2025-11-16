package configloaders

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/systemconfigs"
	"github.com/iwind/TeaGo/logs"
	"reflect"
)

var sharedUserUIConfig *systemconfigs.UserUIConfig = nil

func LoadUserUIConfig() (*systemconfigs.UserUIConfig, error) {
	locker.Lock()
	defer locker.Unlock()

	config, err := loadUserUIConfig()
	if err != nil {
		return nil, err
	}

	v := reflect.Indirect(reflect.ValueOf(config)).Interface().(systemconfigs.UserUIConfig)
	return &v, nil
}

func ReloadUserUIConfig() error {
	locker.Lock()
	defer locker.Unlock()

	sharedUserUIConfig = nil
	_, err := loadUserUIConfig()
	return err
}

func UpdateUserUIConfig(uiConfig *systemconfigs.UserUIConfig) error {
	locker.Lock()
	defer locker.Unlock()

	var rpcClient, err = rpc.SharedRPC()
	if err != nil {
		return err
	}
	valueJSON, err := json.Marshal(uiConfig)
	if err != nil {
		return err
	}
	_, err = rpcClient.SysSettingRPC().UpdateSysSetting(rpcClient.Context(0), &pb.UpdateSysSettingRequest{
		Code:      systemconfigs.SettingCodeUserUIConfig,
		ValueJSON: valueJSON,
	})
	if err != nil {
		return err
	}
	sharedUserUIConfig = uiConfig

	return nil
}

func loadUserUIConfig() (*systemconfigs.UserUIConfig, error) {
	if sharedUserUIConfig != nil {
		return sharedUserUIConfig, nil
	}
	var rpcClient, err = rpc.SharedRPC()
	if err != nil {
		return nil, err
	}
	resp, err := rpcClient.SysSettingRPC().ReadSysSetting(rpcClient.Context(0), &pb.ReadSysSettingRequest{
		Code: systemconfigs.SettingCodeUserUIConfig,
	})
	if err != nil {
		return nil, err
	}
	if len(resp.ValueJSON) == 0 {
		sharedUserUIConfig = defaultUserUIConfig()
		return sharedUserUIConfig, nil
	}

	var config = &systemconfigs.UserUIConfig{}
	err = json.Unmarshal(resp.ValueJSON, config)
	if err != nil {
		logs.Println("[USER_UI_MANAGER]" + err.Error())
		sharedUserUIConfig = defaultUserUIConfig()
		return sharedUserUIConfig, nil
	}

	sharedUserUIConfig = config
	return sharedUserUIConfig, nil
}

func defaultUserUIConfig() *systemconfigs.UserUIConfig {
	return &systemconfigs.UserUIConfig{
		ProductName:         "LingCDN",
		UserSystemName:      "LingCDN管理系统 用户端",
		ShowPageFooter:      false,
		ShowVersion:         true,
		ShowFinance:         true,
		BandwidthUnit:       systemconfigs.BandwidthUnitBit,
		ShowBandwidthCharts: true,
		ShowTrafficCharts:   true,
		TimeZone:            "Asia/Shanghai",
	}
}
