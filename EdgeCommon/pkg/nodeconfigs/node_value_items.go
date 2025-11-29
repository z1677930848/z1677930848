package nodeconfigs

// NodeValueItem 节点上报的数值型指标代码
type NodeValueItem = string

const (
	// 网络流量
	NodeValueItemAllTraffic NodeValueItem = "allTraffic"
	NodeValueItemTrafficIn  NodeValueItem = "trafficIn"
	NodeValueItemTrafficOut NodeValueItem = "trafficOut"
	// 系统资源
	NodeValueItemCPU            NodeValueItem = "cpu"
	NodeValueItemMemory         NodeValueItem = "memory"
	NodeValueItemDisk           NodeValueItem = "disk"
	NodeValueItemCacheDir       NodeValueItem = "cacheDir"
	NodeValueItemLoad           NodeValueItem = "load"
	NodeValueItemConnections    NodeValueItem = "connections"
	NodeValueItemNetworkPackets NodeValueItem = "networkPackets"
	// 请求统计
	NodeValueItemRequests       NodeValueItem = "requests"
	NodeValueItemAttackRequests NodeValueItem = "attackRequests"
)
