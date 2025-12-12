package teaconst

const (
	// 版本号：管理员端发布版本，每次发布需更新
	Version = "1.1.11"

	// APINodeVersion：打包时用于校验 EdgeAPI / API 节点版本号，需要与 API 同步更新
	APINodeVersion = "1.1.0"

	// 产品名称（英文）
	ProductName = "LingCDN"

	// 进程/服务名称，用于 socket、服务注册等
	ProcessName = "lingcdnadmin"

	// 产品名称（中文）
	ProductNameZH = "LingCDN"

	// 默认角色（管理员端）
	Role = "admin"

	// RPC 加密算法
	EncryptMethod = "aes-256-cfb"

	// 通用服务端错误提示
	ErrServer = "服务器出了点小问题，请联系技术人员处理。"

	// 会话 cookie 名称
	CookieSID = "sksid"

	// systemd 服务名
	SystemdServiceName = "lingcdnadmin"

	// 更新检查地址（可通过环境变量或 CI 覆盖）
	UpdatesURL = "https://dl.lingcdn.cloud/api/boot/versions?os=${os}&arch=${arch}"
)
