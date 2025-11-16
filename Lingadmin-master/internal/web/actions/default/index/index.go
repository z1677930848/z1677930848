package index

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/setup"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/types"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"net/url"
	"strings"
	"time"
)

type IndexAction struct {
	actionutils.ParentAction
}

// TokenKey 加密用的密钥
var TokenKey = stringutil.Rand(32)

func (this *IndexAction) RunGet(params struct {
	From string

	Auth *helpers.UserShouldAuth
}) {
	// ⭐ 端口判断：如果是8080端口访问根路径，重定向到用户端
	requestHost := this.Request.Host
	if strings.Contains(requestHost, ":8080") {
		this.RedirectURL("/user")
		return
	}

	// DEMO模式
	this.Data["isDemo"] = teaconst.IsDemoMode

	// 检查系统是否已经配置过
	if !setup.IsConfigured() {
		this.RedirectURL("/setup")
		return
	}

	// 是否新安装
	if setup.IsNewInstalled() {
		this.RedirectURL("/setup/confirm")
		return
	}

	// 已登录跳转到dashboard
	if params.Auth.IsUser() {
		this.RedirectURL("/dashboard")
		return
	}

	this.Data["isUser"] = false
	this.Data["menu"] = "signIn"

	var timestamp = types.String(time.Now().Unix())
	this.Data["token"] = stringutil.Md5(TokenKey+timestamp) + timestamp
	this.Data["from"] = params.From

	uiConfig, err := configloaders.LoadAdminUIConfig()
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["systemName"] = uiConfig.AdminSystemName
	this.Data["showVersion"] = uiConfig.ShowVersion
	if len(uiConfig.Version) > 0 {
		this.Data["version"] = uiConfig.Version
	} else {
		this.Data["version"] = teaconst.Version
	}
	this.Data["faviconFileId"] = uiConfig.FaviconFileId
	this.Data["logoFileId"] = uiConfig.LogoFileId

	_, err = configloaders.LoadSecurityConfig()
	if err != nil {
		this.Data["rememberLogin"] = false
	} else {
		this.Data["rememberLogin"] = true // 允许记住登录
	}

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	Username string
	Password string
	Token    string
	Remember bool
	From     string

	Must *actions.Must
	Auth *helpers.UserShouldAuth
}) {
	params.Must.
		Field("username", params.Username).
		Require("请输入用户名").
		Field("password", params.Password).
		Require("请输入密码")

	// check token
	if len(params.Token) < 32 {
		this.Fail("登录验证失败，请刷新后重试")
		return
	}
	if params.Token != stringutil.Md5(TokenKey+types.String(time.Now().Unix()))+types.String(time.Now().Unix()) {
		timestamp := params.Token[32:]
		if types.Int64(timestamp) < time.Now().Unix()-1800 {
			this.Fail("登录验证已过期，请刷新后重试")
			return
		}
		if params.Token != stringutil.Md5(TokenKey+timestamp)+timestamp {
			this.Fail("登录验证失败，请刷新后重试")
			return
		}
	}

	// 询问API节点
	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 检查登录是否允许
	resp, err := rpcClient.AdminRPC().LoginAdmin(rpcClient.Context(0), &pb.LoginAdminRequest{
		Username: params.Username,
		Password: params.Password,
	})
	if err != nil {
		this.FailField("username", "用户名或密码错误")
		return
	}

	var adminId = resp.AdminId
	if adminId <= 0 {
		this.FailField("username", "用户名或密码错误")
		return
	}

	params.Auth.StoreAdmin(adminId, params.Remember)

	// 检查from参数
	if len(params.From) > 0 {
		fromURL, err := url.Parse(params.From)
		if err == nil {
			if len(fromURL.Scheme) == 0 {
				this.Data["url"] = params.From
				this.Success()
				return
			}
		}
	}

	this.Data["url"] = "/dashboard"
	this.Success()
}
