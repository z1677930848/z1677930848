// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.
//go:build !plus

package users

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"time"
)

type PlansAction struct {
	actionutils.ParentAction
}

func (this *PlansAction) RunPost(params struct {
	UserId   int64
	ServerId int64
}) {
	today := time.Now().Format("2006-01-02")

	// 闇€瑕佺鐞嗗憳鏉冮檺
	helper := helpers.NewUserMustAuth("servers")
	if !helper.BeforeAction(this, "") {
		return
	}

	// 濂楅瀹氫箟
	plans := []maps.Map{}
	plansResp, err := this.RPC().PlanRPC().FindAllAvailablePlans(this.AdminContext(), &pb.FindAllAvailablePlansRequest{})
	if err == nil && plansResp != nil {
		for _, p := range plansResp.Plans {
			plans = append(plans, maps.Map{
				"id":          p.Id,
				"name":        p.Name,
				"description": p.Description,
				"isOn":        p.IsOn,
			})
		}
	}

	// 鐢ㄦ埛宸茶喘濂楅
	userPlans := []maps.Map{}
	if params.UserId > 0 {
		userCtx := this.AdminContext()
		userPlansResp, err := this.RPC().UserPlanRPC().ListEnabledUserPlans(userCtx, &pb.ListEnabledUserPlansRequest{
			UserId: params.UserId,
			Offset: 0,
			Size:   50,
		})
		if err == nil && userPlansResp != nil {
			for _, up := range userPlansResp.UserPlans {
				userPlans = append(userPlans, maps.Map{
					"id":        up.Id,
					"planId":    up.PlanId,
					"name":      up.Name,
					"dayTo":     up.DayTo,
					"isOn":      up.IsOn,
					"isExpired": len(up.DayTo) > 0 && up.DayTo < today,
				})
			}
		}
	}

	this.Data["plans"] = plans
	this.Data["userPlans"] = userPlans
	this.Success()
}
