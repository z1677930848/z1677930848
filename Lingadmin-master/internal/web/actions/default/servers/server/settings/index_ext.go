// Copyright 2022 Liuxiangchao iwind.liu@gmail.com. All rights reserved. Official site: https://goedge.cn .
//go:build !plus

package settings

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

func (this *IndexAction) initUserPlan(server *pb.Server) {
	var userPlanMap = maps.Map{"id": server.UserPlanId, "dayTo": "", "plan": maps.Map{}}
	plans := []maps.Map{}

	// 查询当前用户可用套餐
	userId := int64(0)
	if server.User != nil {
		userId = server.User.Id
	}
	resp, err := this.RPC().UserPlanRPC().ListEnabledUserPlans(this.AdminContext(), &pb.ListEnabledUserPlansRequest{
		UserId: userId,
		Offset: 0,
		Size:   200,
	})
	if err == nil && resp != nil {
		for _, up := range resp.UserPlans {
			plans = append(plans, maps.Map{
				"id":    up.Id,
				"name":  up.Name,
				"dayTo": up.DayTo,
			})
			if up.Id == server.UserPlanId {
				userPlanMap["dayTo"] = up.DayTo
			}
		}
	}

	this.Data["plans"] = plans
	this.Data["userPlan"] = userPlanMap
}
