package plans

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/plans/plan"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		// 濂楅鍒楄〃
		server.Get("/plans", new(IndexAction))

		// 鎺掑簭銆佸垹闄?
		server.Post("/plans/sort", new(SortAction))
		server.Post("/plans/plan.delete", new(DeleteAction))

		// 棰勭暀璇︽儏銆佹洿鏂帮紙鏆傛湭瀹炵幇锛?
		server.Get("/plans/plan", new(plan.EmptyAction))
		server.Get("/plans/plan/update", new(plan.EmptyAction))
	})
}
