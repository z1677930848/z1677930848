package configloaders

import (
	_ "github.com/iwind/TeaGo/bootstrap"
	"testing"
	"time"
)

func TestLoadUIConfig(t *testing.T) {
	for i := 0; i < 10; i++ {
		before := time.Now()
		config, err := LoadAdminUIConfig()
		if err != nil {
			t.Skipf("skip because config not available: %v", err)
		}
		t.Log(time.Since(before).Seconds()*1000, "ms")
		t.Logf("%p", config)
	}
}

func TestLoadUIConfig2(t *testing.T) {
	for i := 0; i < 10; i++ {
		config, err := LoadAdminUIConfig()
		if err != nil {
			t.Skipf("skip because config not available: %v", err)
		}
		t.Log(config)
	}
}
