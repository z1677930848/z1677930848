package register

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/settings/license"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	stringutil "github.com/iwind/TeaGo/utils/string"
)

type RegisterAction actions.Action

// RunGet 显示注册页面
func (this *RegisterAction) RunGet(params struct{}) {
	// 检查授权
	lic, err := license.LoadLicense()
	if err == nil && lic != nil {
		if lic.Code != "" && !lic.IsValid {
			this.Data["licenseError"] = "系统授权已过期,暂时无法注册新用户"
			this.Data["showLicenseError"] = true
		}
	}

	// 生成Token
	token := stringutil.Rand(32)
	this.Data["token"] = token

	// 加载UI配置
	config, err := configloaders.LoadUserUIConfig()
	if err == nil {
		this.Data["systemName"] = config.UserSystemName
	} else {
		this.Data["systemName"] = "LingCDN"
	}

	this.Show()
}

// RunPost 处理注册提交
func (this *RegisterAction) RunPost(params struct {
	Username string
	Password string
	Email    string
	Fullname string
	Mobile   string
	Token    string

	Must *actions.Must
	Auth *helpers.UserShouldAuth
}) {
	// 检查授权
	lic, licErr := license.LoadLicense()
	if licErr == nil && lic != nil {
		if lic.Code != "" && !lic.IsValid {
			this.Fail("系统授权已过期，暂时无法注册新用户")
			return
		}
	}

	// 验证输入
	params.Must.
		Field("username", params.Username).
		Require("请输入用户名").
		Match("^[a-zA-Z0-9_]{4,30}$", "用户名只能包含字母、数字和下划线，长度4-30位")

	params.Must.
		Field("password", params.Password).
		Require("请输入密码").
		MinLength(6, "密码至少需要6位")

	params.Must.
		Field("fullname", params.Fullname).
		Require("请输入真实姓名")

	params.Must.
		Field("email", params.Email).
		Require("请输入邮箱").
		Email("请输入正确的邮箱地址")

	// 连接RPC
	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		this.Fail("系统错误，请稍后重试")
		return
	}

	// 检查用户名是否已存在
	checkResp, err := rpcClient.UserRPC().CheckUserUsername(rpcClient.Context(0), &pb.CheckUserUsernameRequest{
		UserId:   0,
		Username: params.Username,
	})
	if err != nil {
		this.Fail("系统错误：" + err.Error())
		return
	}
	if checkResp.Exists {
		this.Fail("用户名已被占用，请换一个")
		return
	}

	// 创建用户（不需要集群ID，后续由管理员分配）
	createResp, err := rpcClient.UserRPC().CreateUser(rpcClient.Context(0), &pb.CreateUserRequest{
		Username:      params.Username,
		Password:      params.Password,
		Fullname:      params.Fullname,
		Mobile:        params.Mobile,
		Email:         params.Email,
		Source:        "web_register",
		NodeClusterId: 0, // 注册时不分配集群
	})
	if err != nil {
		this.Fail("注册失败：" + err.Error())
		return
	}

	// 自动登录
	params.Auth.StoreAdmin(createResp.UserId, false)

	this.Data["userId"] = createResp.UserId
	this.Success()
}
