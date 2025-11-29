package profile

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"time"
)

type IndexAction actions.Action

// RunGet 显示个人设置
func (this *IndexAction) RunGet(params struct {
	Auth *helpers.UserShouldAuth
}) {
	if !params.Auth.IsUserPortal() {
		this.RedirectURL("/user")
		return
	}

	userId := params.Auth.UserId()
	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		this.Fail("系统错误，请稍后重试")
		return
	}
	userCtx := rpcClient.UserContext(userId)

	userResp, err := rpcClient.UserRPC().FindEnabledUser(userCtx, &pb.FindEnabledUserRequest{UserId: userId})
	if err != nil || userResp == nil || userResp.User == nil {
		this.Fail("无法读取用户信息，请稍后重试")
		return
	}
	user := userResp.User

	this.Data["title"] = "个人设置"
	this.Data["userId"] = user.Id
	this.Data["username"] = user.Username
	this.Data["email"] = user.Email
	this.Data["mobile"] = user.Mobile
	this.Data["fullname"] = user.Fullname
	this.Data["planName"] = "未分配"
	this.Data["company"] = ""
	if user.CreatedAt > 0 {
		this.Data["createdAt"] = time.Unix(user.CreatedAt, 0).Format("2006-01-02 15:04:05")
	} else {
		this.Data["createdAt"] = ""
	}

	this.Show()
}

// UpdateAction 更新基本信息
type UpdateAction actions.Action

func (this *UpdateAction) RunPost(params struct {
	Email    string
	Mobile   string
	Username string
	Company  string

	Must *actions.Must
	Auth *helpers.UserShouldAuth
}) {
	if !params.Auth.IsUserPortal() {
		this.Fail("请先登录")
		return
	}
	params.Must.
		Field("email", params.Email).
		Require("请输入邮箱")

	userId := params.Auth.UserId()
	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		this.Fail("系统错误，请稍后再试")
		return
	}
	userCtx := rpcClient.UserContext(userId)
	currentFullname := ""
	if userResp, infoErr := rpcClient.UserRPC().FindEnabledUser(userCtx, &pb.FindEnabledUserRequest{UserId: userId}); infoErr == nil && userResp != nil && userResp.User != nil {
		currentFullname = userResp.User.Fullname
	}
	if len(currentFullname) == 0 {
		currentFullname = params.Username
	}
	_, err = rpcClient.UserRPC().UpdateUserInfo(userCtx, &pb.UpdateUserInfoRequest{
		UserId:   userId,
		Fullname: currentFullname,
		Mobile:   params.Mobile,
		Email:    params.Email,
	})
	if err != nil {
		this.Fail("保存失败：" + err.Error())
		return
	}
	this.Success()
}

// ChangePasswordAction 修改密码
type ChangePasswordAction actions.Action

func (this *ChangePasswordAction) RunPost(params struct {
	NewPassword string

	Must *actions.Must
	Auth *helpers.UserShouldAuth
}) {
	if !params.Auth.IsUserPortal() {
		this.Fail("请先登录")
		return
	}
	params.Must.
		Field("newPassword", params.NewPassword).
		Require("请输入新密码")

	userId := params.Auth.UserId()
	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		this.Fail("系统错误，请稍后再试")
		return
	}
	userCtx := rpcClient.UserContext(userId)
	username := params.Auth.Username()
	if len(username) == 0 {
		if userResp, infoErr := rpcClient.UserRPC().FindEnabledUser(userCtx, &pb.FindEnabledUserRequest{UserId: userId}); infoErr == nil && userResp != nil && userResp.User != nil {
			username = userResp.User.Username
		}
	}
	_, err = rpcClient.UserRPC().UpdateUserLogin(userCtx, &pb.UpdateUserLoginRequest{
		UserId:   userId,
		Username: username,
		Password: params.NewPassword,
	})
	if err != nil {
		this.Fail("修改失败：" + err.Error())
		return
	}
	this.Success()
}
