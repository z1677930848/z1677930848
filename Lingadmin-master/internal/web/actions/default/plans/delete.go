package plans

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	PlanId int64
}) {
	if params.PlanId <= 0 {
		this.Fail("请选择要删除的套餐")
		return
	}
	_, err := this.RPC().PlanRPC().DeletePlan(this.AdminContext(), &pb.DeletePlanRequest{PlanId: params.PlanId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Success()
}
