package index

import (
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/settings/license"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"time"
)

type IndexAction actions.Action

// TokenKey 加密用的密钥
var TokenKey = stringutil.Rand(32)

// RunGet 显示用户登录页面
func (this *IndexAction) RunGet(params struct {
	Auth *helpers.UserShouldAuth
}) {
	// 已登录跳转到dashboard
	if params.Auth.IsUserPortal() {
		this.RedirectURL("/user/dashboard")
		return
	}

	// 检查授权（仅显示警告，不阻止登录）
	lic, err := license.LoadLicense()
	if err == nil && lic != nil {
		if lic.Code != "" && !lic.IsValid {
			this.Data["licenseError"] = "系统授权已过期，请联系管理员更新授权"
			this.Data["showLicenseError"] = true
		} else {
			this.Data["showLicenseError"] = false
		}
	} else {
		this.Data["licenseError"] = "系统未授权，请联系管理员配置授权"
		this.Data["showLicenseError"] = true
	}

	this.Data["isUser"] = true
	this.Data["menu"] = "userSignIn"

	var timestamp = fmt.Sprintf("%d", time.Now().Unix())
	this.Data["token"] = stringutil.Md5(TokenKey+timestamp) + timestamp

	// 加载用户界面配置
	userUIConfig, err := configloaders.LoadUserUIConfig()
	if err != nil {
		this.Data["systemName"] = "LingCDN"
	} else {
		if len(userUIConfig.UserSystemName) > 0 {
			this.Data["systemName"] = userUIConfig.UserSystemName
		} else {
			this.Data["systemName"] = userUIConfig.ProductName
		}

		// 设置favicon和logo
		if userUIConfig.FaviconFileId > 0 {
			this.Data["faviconURL"] = fmt.Sprintf("/ui/image/%d", userUIConfig.FaviconFileId)
		}
		if userUIConfig.LogoFileId > 0 {
			this.Data["logoURL"] = fmt.Sprintf("/ui/image/%d", userUIConfig.LogoFileId)
		}
	}
	this.Data["version"] = teaconst.Version
	this.Data["rememberLogin"] = true

	this.Show()
}

// RunPost 处理用户登录
func (this *IndexAction) RunPost(params struct {
	Username string
	Password string
	Token    string
	Remember bool

	Must *actions.Must
	Auth *helpers.UserShouldAuth
}) {
	params.Must.
		Field("username", params.Username).
		Require("请输入用户名").
		Field("password", params.Password).
		Require("请输入密码")

	// 检查授权（仅记录，不阻止登录）
	_, _ = license.LoadLicense()

	// check token
	if len(params.Token) < 32 {
		this.Fail("登录验证失败，请刷新后重试")
		return
	}
	var timestamp = params.Token[32:]
	if len(timestamp) == 0 {
		this.Fail("登录验证失败，请刷新后重试")
		return
	}
	if stringutil.Md5(TokenKey+timestamp) != params.Token[:32] {
		this.Fail("登录验证失败，请刷新后重试")
		return
	}

	// 通过RPC验证用户账号
	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		this.Fail("系统错误，请稍后重试")
		return
	}

	// 调用UserRPC登录验证
	resp, err := rpcClient.UserRPC().LoginUser(rpcClient.Context(0), &pb.LoginUserRequest{
		Username: params.Username,
		Password: params.Password,
	})
	if err != nil {
		this.FailField("username", "用户名或密码错误")
		return
	}

	// 检查登录是否成功
	if !resp.IsOk {
		if len(resp.Message) > 0 {
			this.Fail(resp.Message)
		} else {
			this.FailField("username", "用户名或密码错误")
		}
		return
	}

	var userId = resp.UserId
	if userId <= 0 {
		this.FailField("username", "用户名或密码错误")
		return
	}

	// 使用真实的用户ID创建会话
	params.Auth.StoreUser(userId, params.Username, params.Remember)

	this.Success()
}
