package nodeconfigs

import "encoding/json"

// FindNodeValueItemName 返回指标名称
func FindNodeValueItemName(item NodeValueItem) string {
	switch item {
	case NodeValueItemAllTraffic:
		return "全部流量"
	case NodeValueItemTrafficIn:
		return "入站流量"
	case NodeValueItemTrafficOut:
		return "出站流量"
	case NodeValueItemCPU:
		return "CPU"
	case NodeValueItemMemory:
		return "内存"
	case NodeValueItemDisk:
		return "磁盘"
	case NodeValueItemCacheDir:
		return "缓存目录"
	case NodeValueItemLoad:
		return "负载"
	case NodeValueItemConnections:
		return "连接数"
	case NodeValueItemNetworkPackets:
		return "网络包"
	case NodeValueItemRequests:
		return "请求数"
	case NodeValueItemAttackRequests:
		return "攻击请求数"
	default:
		return string(item)
	}
}

// FindNodeValueItemParamName 返回指标参数名称，目前参数直接透传
func FindNodeValueItemParamName(_ NodeValueItem, param string) string {
	return param
}

// FindNodeValueOperatorName 返回操作符名称
func FindNodeValueOperatorName(op NodeValueOperator) string {
	switch op {
	case NodeValueOperatorGt:
		return ">"
	case NodeValueOperatorGte:
		return ">="
	case NodeValueOperatorLt:
		return "<"
	case NodeValueOperatorLte:
		return "<="
	case NodeValueOperatorEq:
		return "="
	case NodeValueOperatorNeq:
		return "!="
	default:
		return string(op)
	}
}

// UnmarshalNodeValue 将JSON反序列化为通用对象
func UnmarshalNodeValue(valueJSON []byte) any {
	if len(valueJSON) == 0 {
		return nil
	}
	var v any
	if err := json.Unmarshal(valueJSON, &v); err != nil {
		return string(valueJSON)
	}
	return v
}

// FindNodeValueSumMethodName 返回聚合方式名称
func FindNodeValueSumMethodName(method NodeValueSumMethod) string {
	switch method {
	case NodeValueSumMethodAvg:
		return "平均值"
	case NodeValueSumMethodSum:
		return "求和"
	default:
		return string(method)
	}
}

// FindNodeValueDurationUnitName 返回持续时间单位名称
func FindNodeValueDurationUnitName(unit NodeValueDurationUnit) string {
	switch unit {
	case NodeValueDurationUnitMinute:
		return "分钟"
	case NodeValueDurationUnitHour:
		return "小时"
	case NodeValueDurationUnitDay:
		return "天"
	default:
		return string(unit)
	}
}
