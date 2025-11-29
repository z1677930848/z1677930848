package plans

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "plan", "index")
}

func (this *IndexAction) RunGet(params struct {
	Keyword string
}) {
	// 是否开启套餐计费：暂设为可用
	this.Data["canUsePlans"] = true

	// 读取所有套餐（分页简化）
	resp, err := this.RPC().PlanRPC().ListEnabledPlans(this.AdminContext(), &pb.ListEnabledPlansRequest{
		Offset: 0,
		Size:   200,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 查询集群名称
	clusterMap := map[int64]maps.Map{}
	clustersResp, err := this.RPC().NodeClusterRPC().FindAllEnabledNodeClusters(this.AdminContext(), &pb.FindAllEnabledNodeClustersRequest{})
	if err == nil && clustersResp != nil {
		for _, c := range clustersResp.NodeClusters {
			clusterMap[c.Id] = maps.Map{
				"id":   c.Id,
				"name": c.Name,
			}
		}
	}

	plans := []maps.Map{}
	for _, plan := range resp.Plans {
		cluster := maps.Map{"id": 0, "name": ""}
		if c, ok := clusterMap[plan.ClusterId]; ok {
			cluster = c
		}
		plans = append(plans, maps.Map{
			"id":                          plan.Id,
			"isOn":                        plan.IsOn,
			"name":                        plan.Name,
			"description":                 plan.Description,
			"priceType":                   plan.PriceType,
			"cluster":                     cluster,
			"totalServers":                plan.TotalServers,
			"totalServerNamesPerServer":   plan.TotalServerNamesPerServer,
			"totalServerNames":            plan.TotalServerNames,
			"monthlyPrice":                plan.MonthlyPrice,
			"seasonallyPrice":             plan.SeasonallyPrice,
			"yearlyPrice":                 plan.YearlyPrice,
			"trafficLimitJSON":            plan.TrafficLimitJSON,
			"bandwidthLimitPerNodeJSON":   plan.BandwidthLimitPerNodeJSON,
			"trafficPriceJSON":            plan.TrafficPriceJSON,
			"bandwidthPriceJSON":          plan.BandwidthPriceJSON,
			"dailyRequests":               plan.DailyRequests,
			"monthlyRequests":             plan.MonthlyRequests,
			"dailyWebsocketConnections":   plan.DailyWebsocketConnections,
			"monthlyWebsocketConnections": plan.MonthlyWebsocketConnections,
			"maxUploadSizeJSON":           plan.MaxUploadSizeJSON,
		})
	}

	this.Data["plans"] = plans
	this.Data["page"] = ""
	this.Show()
}
