// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package serverconfigs

import timeutil "github.com/iwind/TeaGo/utils/time"

// DefaultPlanExpireNoticePageBody 濂楅杩囨湡鏃舵彁绀?
const DefaultPlanExpireNoticePageBody = `<!DOCTYPE html>
<html>
<head>
<title>濂楅宸茶繃鏈?/title>
<body>

<h1>濂楅宸茶繃鏈燂紝璇峰強鏃剁画璐广€?/h1>
<p>Your server plan has been expired, please renew the plan.</p>
<address>Request ID: ${requestId}.</address>

</body>
</html>`

// UserPlanConfig 鐢ㄦ埛濂楅閰嶇疆
type UserPlanConfig struct {
	Id     int64  `yaml:"id" json:"id"`         // 鐢ㄦ埛濂楅ID
	DayTo  string `yaml:"dayTo" json:"dayTo"`   // 鏈夋晥鏈?
	PlanId int64  `yaml:"planId" json:"planId"` // 濂楅瀹氫箟ID
}

// Init 鍒濆鍖?
func (this *UserPlanConfig) Init() error {
	return nil
}

// IsAvailable 鏄惁鏈夋晥
func (this *UserPlanConfig) IsAvailable() bool {
	return this.DayTo >= timeutil.Format("Y-m-d")
}
