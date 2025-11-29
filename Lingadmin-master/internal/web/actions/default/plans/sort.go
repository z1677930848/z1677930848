package plans

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type SortAction struct {
	actionutils.ParentAction
}

func (this *SortAction) RunPost(params struct {
	Ids []int64
}) {
	_, err := this.RPC().PlanRPC().SortPlans(this.AdminContext(), &pb.SortPlansRequest{PlanIds: params.Ids})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Success()
}
