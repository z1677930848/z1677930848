package configloaders

import (
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestLoadAdminModuleMapping(t *testing.T) {
	m, err := loadAdminModuleMapping()
	if err != nil {
		t.Skipf("skip because config not available: %v", err)
	}
	logs.PrintAsJSON(m, t)
}
