package register

import (
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		// 用户注册（无需认证）
		server.Prefix("/user/register").
			Data("teaMenu", "user").
			GetPost("", new(RegisterAction)).
			EndAll()
	})
}
