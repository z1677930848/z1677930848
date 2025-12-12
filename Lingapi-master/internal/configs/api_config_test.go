package configs

import (
	"testing"

	_ "github.com/iwind/TeaGo/bootstrap"
)

func TestSharedAPIConfig(t *testing.T) {
	t.Skip("skip: requires local configs/api.yaml")
	config, err := SharedAPIConfig()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(config)
}
