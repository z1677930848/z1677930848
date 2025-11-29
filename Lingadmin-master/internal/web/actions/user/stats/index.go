package stats

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"time"
)

type IndexAction actions.Action

// RunGet 显示统计数据
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

	this.Data["title"] = "统计数据"

	// 域名下拉
	domainOptions := []maps.Map{}
	if serversResp, err := rpcClient.ServerRPC().FindAllUserServers(userCtx, &pb.FindAllUserServersRequest{UserId: userId}); err == nil && serversResp != nil {
		for _, server := range serversResp.Servers {
			domainOptions = append(domainOptions, maps.Map{
				"id":   server.Id,
				"name": server.Name,
			})
		}
	}
	this.Data["domains"] = domainOptions

	// 统计最近7天
	dayTo := time.Now().Format("20060102")
	dayFrom := time.Now().AddDate(0, 0, -6).Format("20060102")
	var totalRequests int64 = 0
	var totalTraffic int64 = 0
	var totalCached int64 = 0

	statsResp, err := rpcClient.ServerDailyStatRPC().FindServerDailyStatsBetweenDays(userCtx, &pb.FindServerDailyStatsBetweenDaysRequest{
		UserId:  userId,
		DayFrom: dayFrom,
		DayTo:   dayTo,
	})
	if err == nil && statsResp != nil {
		for _, stat := range statsResp.Stats {
			totalRequests += stat.CountRequests
			totalTraffic += stat.Bytes
			totalCached += stat.CountCachedRequests
		}
	}

	avgHitRate := 0
	if totalRequests > 0 {
		avgHitRate = int(totalCached * 100 / totalRequests)
	}

	this.Data["totalRequests"] = numberutils.FormatCount(totalRequests)
	this.Data["totalTraffic"] = numberutils.FormatBytes(totalTraffic)
	this.Data["avgHitRate"] = avgHitRate
	this.Data["avgResponseTime"] = 0

	// 状态码和分布数据目前接口未开放，填充0/空列表以避免模板错误
	this.Data["status2xx"] = 0
	this.Data["status3xx"] = 0
	this.Data["status4xx"] = 0
	this.Data["status5xx"] = 0
	this.Data["status2xxPercent"] = 0
	this.Data["status3xxPercent"] = 0
	this.Data["status4xxPercent"] = 0
	this.Data["status5xxPercent"] = 0
	this.Data["topRegions"] = []maps.Map{}
	this.Data["topUrls"] = []maps.Map{}

	this.Show()
}
