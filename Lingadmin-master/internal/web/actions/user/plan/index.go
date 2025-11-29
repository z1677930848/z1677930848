package plan

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"time"
)

type IndexAction actions.Action

// RunGet 显示用户套餐信息
func (this *IndexAction) RunGet(params struct {
	Auth *helpers.UserShouldAuth
}) {
	today := time.Now().Format("2006-01-02")

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
	ctx := rpcClient.UserContext(userId)

	plans := []maps.Map{}
	resp, err := rpcClient.UserPlanRPC().ListEnabledUserPlans(ctx, &pb.ListEnabledUserPlansRequest{
		UserId: userId,
		Offset: 0,
		Size:   20,
	})
	if err == nil && resp != nil {
		for _, item := range resp.UserPlans {
			planName := ""
			planDetail, errFind := rpcClient.PlanRPC().FindEnabledPlan(ctx, &pb.FindEnabledPlanRequest{PlanId: item.PlanId})
			if errFind == nil && planDetail != nil && planDetail.Plan != nil {
				planName = planDetail.Plan.Name
			}
			isExpired := len(item.DayTo) > 0 && item.DayTo < today
			plans = append(plans, maps.Map{
				"id":         item.Id,
				"name":       planName,
				"customName": item.Name,
				"dayTo":      item.DayTo,
				"isOn":       item.IsOn,
				"isExpired":  isExpired,
			})
		}
	}

	this.Data["title"] = "套餐信息"
	this.Data["plans"] = plans

	this.Show()
}
