package nodeconfigs

// NodeValueRange 统计查询的时间范围
type NodeValueRange = string

const (
	NodeValueRangeMinute NodeValueRange = "minute"
	NodeValueRangeHour   NodeValueRange = "hour"
	NodeValueRangeDay    NodeValueRange = "day"
)

// NodeValueSumMethod 汇总方式
type NodeValueSumMethod = string

const (
	NodeValueSumMethodAvg NodeValueSumMethod = "avg"
	NodeValueSumMethodSum NodeValueSumMethod = "sum"
)
