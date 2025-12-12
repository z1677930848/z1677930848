package bootstrap

import "github.com/iwind/TeaGo/dbs"

// Placeholder bootstrap package to satisfy legacy imports.
// We trigger dbs.NotifyReady to initialize Shared*DAO vars.
func init() {
	dbs.NotifyReady()
}
