package dashboard

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"math"
)

type IndexAction actions.Action

// RunGet 显示用户仪表盘
func (this *IndexAction) RunGet(params struct {
	Auth *helpers.UserShouldAuth
}) {
	// 确认登录状态
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

	// 仪表盘汇总数据
	var todayRequests int64 = 0
	var todayTrafficBytes int64 = 0
	var totalTrafficBytes int64 = 0
	var totalDomains int64 = 0
	var activeDomains int64 = 0
	var sslCerts int64 = 0
	var expiringSoon int64 = 0

	dashboardResp, err := rpcClient.UserRPC().ComposeUserDashboard(userCtx, &pb.ComposeUserDashboardRequest{UserId: userId})
	if err == nil && dashboardResp != nil {
		totalDomains = dashboardResp.CountServers
		totalTrafficBytes = dashboardResp.MonthlyTrafficBytes
		// 今日流量直接使用聚合结果
		todayTrafficBytes = dashboardResp.DailyTrafficBytes
	}

	// 今日请求与缓存
	day := timeutil.Format("Ymd")
	todayStatResp, err := rpcClient.ServerDailyStatRPC().SumServerDailyStats(userCtx, &pb.SumServerDailyStatsRequest{
		UserId: userId,
		Day:    day,
	})
	if err == nil && todayStatResp != nil && todayStatResp.ServerDailyStat != nil {
		todayRequests = todayStatResp.ServerDailyStat.CountRequests
		todayTrafficBytes = todayStatResp.ServerDailyStat.Bytes
	}

	// 服务列表
	domains := []maps.Map{}
	serversResp, err := rpcClient.ServerRPC().FindAllUserServers(userCtx, &pb.FindAllUserServersRequest{UserId: userId})
	if err == nil && serversResp != nil {
		for _, server := range serversResp.Servers {
			// 统计每个域名当日数据
			reqs := int64(0)
			cached := int64(0)
			statResp, statErr := rpcClient.ServerDailyStatRPC().SumServerDailyStats(userCtx, &pb.SumServerDailyStatsRequest{
				ServerId: server.Id,
				Day:      day,
			})
			if statErr == nil && statResp != nil && statResp.ServerDailyStat != nil {
				reqs = statResp.ServerDailyStat.CountRequests
				cached = statResp.ServerDailyStat.CountCachedRequests
			}
			hitRate := 0
			if reqs > 0 {
				hitRate = int(math.Round(float64(cached*100) / float64(reqs)))
			}
			domains = append(domains, maps.Map{
				"id":        server.Id,
				"name":      server.Name,
				"requests":  numberutils.FormatCount(reqs),
				"status":    map[bool]string{true: "running", false: "stopped"}[server.IsOn],
				"hitRate":   hitRate,
				"createdAt": server.CreatedAt,
			})
			if server.IsOn {
				activeDomains++
			}
		}
	}
	if totalDomains == 0 {
		totalDomains = int64(len(domains))
	}

	this.Data["title"] = "用户仪表盘"
	this.Data["totalDomains"] = totalDomains
	this.Data["activeDomains"] = activeDomains
	this.Data["todayRequests"] = numberutils.FormatCount(todayRequests)
	this.Data["todayTraffic"] = numberutils.FormatBytes(todayTrafficBytes)
	this.Data["totalTraffic"] = numberutils.FormatBytes(totalTrafficBytes)
	this.Data["sslCerts"] = sslCerts
	this.Data["expiringSoon"] = expiringSoon
	this.Data["domains"] = domains

	// 用户名展示
	username := params.Auth.Username()
	if len(username) == 0 {
		if userInfoResp, infoErr := rpcClient.UserRPC().FindEnabledUser(userCtx, &pb.FindEnabledUserRequest{UserId: userId}); infoErr == nil && userInfoResp != nil && userInfoResp.User != nil {
			if len(userInfoResp.User.Fullname) > 0 {
				username = userInfoResp.User.Fullname
			} else {
				username = userInfoResp.User.Username
			}
		}
	}
	this.Data["teaUsername"] = username

	this.Show()
}
