package domains

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"math"
	"time"
)

type IndexAction actions.Action

// RunGet 显示域名列表
func (this *IndexAction) RunGet(params struct {
	Auth *helpers.UserShouldAuth
}) {
	// 检查是否登录
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

	this.Data["title"] = "域名管理"

	day := timeutil.Format("Ymd")
	domains := []maps.Map{}
	serversResp, err := rpcClient.ServerRPC().FindAllUserServers(userCtx, &pb.FindAllUserServersRequest{UserId: userId})
	if err == nil && serversResp != nil {
		for _, server := range serversResp.Servers {
			reqs := int64(0)
			traffic := int64(0)
			cached := int64(0)
			statResp, statErr := rpcClient.ServerDailyStatRPC().SumServerDailyStats(userCtx, &pb.SumServerDailyStatsRequest{
				ServerId: server.Id,
				Day:      day,
			})
			if statErr == nil && statResp != nil && statResp.ServerDailyStat != nil {
				reqs = statResp.ServerDailyStat.CountRequests
				traffic = statResp.ServerDailyStat.Bytes
				cached = statResp.ServerDailyStat.CountCachedRequests
			}
			hitRate := 0
			if reqs > 0 {
				hitRate = int(math.Round(float64(cached*100) / float64(reqs)))
			}
			createdAt := ""
			if server.CreatedAt > 0 {
				createdAt = time.Unix(server.CreatedAt, 0).Format("2006-01-02")
			}
			domains = append(domains, maps.Map{
				"id":            server.Id,
				"name":          server.Name,
				"status":        map[bool]string{true: "running", false: "stopped"}[server.IsOn],
				"createdAt":     createdAt,
				"todayRequests": numberutils.FormatCount(reqs),
				"todayTraffic":  numberutils.FormatBytes(traffic),
				"hitRate":       hitRate,
				"hasSSL":        len(server.HttpsJSON) > 0,
			})
		}
	}

	this.Data["domains"] = domains
	this.Show()
}

type CreateAction actions.Action

func (this *CreateAction) RunGet(params struct {
	Auth *helpers.UserShouldAuth
}) {
	if !params.Auth.IsUserPortal() {
		this.RedirectURL("/user")
		return
	}

	this.Data["title"] = "添加域名"
	this.Show()
}

func (this *CreateAction) RunPost(params struct {
	Domain string

	Must *actions.Must
	Auth *helpers.UserShouldAuth
}) {
	if !params.Auth.IsUserPortal() {
		this.Fail("请先登录")
		return
	}

	params.Must.
		Field("domain", params.Domain).
		Require("请输入域名")

	// 用户端暂不开放自助创建域名，提示用户联系管理员
	this.Fail("当前版本暂未开放自助创建域名，请联系管理员协助添加。")
}

type DeleteAction actions.Action

func (this *DeleteAction) RunPost(params struct {
	DomainId int64

	Must *actions.Must
	Auth *helpers.UserShouldAuth
}) {
	if !params.Auth.IsUserPortal() {
		this.Fail("请先登录")
		return
	}

	params.Must.
		Field("domainId", params.DomainId).
		Require("请选择要删除的域名")

	// 用户端暂不开放自助删除域名
	this.Fail("当前版本暂未开放自助删除域名，请联系管理员。")
}
