package teaconst

// Package-level constants used across the application.
// Keep these stable and document purpose when changing values.
const (
	// Version: EdgeAdmin release version. Bumped for releases.
	Version = "1.0.11"

	// APINodeVersion: expected EdgeAPI/ApiNode release version used when packaging.
	// Avoid changing this unless you also update EdgeAPI accordingly.
	APINodeVersion = "1.0.4"

	// ProductName: user-facing product name.
	ProductName = "LingCDN"

	// ProcessName: short process/service name used for sockets, service registration, etc.
	// We standardize on 'lingcdnadmin' for the runtime/service name.
	ProcessName = "lingcdnadmin"

	// ProductNameZH: Chinese product name (for UI/messages).
	ProductNameZH = "LingCDN"

	// Role: default role for this binary (admin).
	Role = "admin"

	// EncryptMethod: symmetric encryption algorithm used in RPC token construction.
	EncryptMethod = "aes-256-cfb"

	// ErrServer: generic user-facing server error message.
	ErrServer = "服务器出了点小问题，请联系技术人员处理。"

	// CookieSID: cookie name for session id.
	CookieSID = "sksid"

	// SystemdServiceName: name used when installing systemd service on Linux.
	SystemdServiceName = "lingcdnadmin"

	// UpdatesURL: URL used to query available updates; may be overridden in environment or CI.
	UpdatesURL = "http://dl.lingcdn.cloud/api/boot/versions?os=${os}&arch=${arch}"
)
