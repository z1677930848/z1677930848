package teaconst

const (
	// 版本号：保持与管理端的节点版本一致
	Version = "1.1.0"

	ProductName = "Ling node"
	ProcessName = "ling-node"

	Role = "node"

	EncryptMethod = "aes-256-cfb"

	// SystemdServiceName systemd
	SystemdServiceName = "ling-node"

	AccessLogSockName    = "ling-node.accesslog"
	CacheGarbageSockName = "ling-node.cache.garbage"

	EnableKVCacheStore = true // determine store cache keys in KVStore or sqlite
)
