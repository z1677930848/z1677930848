package models

import (
	"fmt"
	"os"
	"testing"

	"github.com/iwind/TeaGo/dbs"
)

func TestMain(m *testing.M) {
	if _, err := dbs.Default(); err != nil {
		fmt.Println("skip db-dependent tests:", err)
		os.Exit(0)
	}
	os.Exit(m.Run())
}

