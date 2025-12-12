package dns_test

import (
	"testing"
	"time"

	"github.com/TeaOSLab/EdgeAPI/internal/db/models/dns"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/iwind/TeaGo/bootstrap"
	"github.com/iwind/TeaGo/dbs"
)

func TestDNSTaskDAO_CreateDNSTask(t *testing.T) {
	if _, err := dbs.Default(); err != nil {
		t.Skipf("skip: %v", err)
	}
	dbs.NotifyReady()
	dao, err := dns.NewDNSTaskDAO()
	if err != nil {
		t.Fatal(err)
	}
	err = dao.CreateDNSTask(nil, 1, 2, 3, 0, "cdn", "taskType")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestDNSTaskDAO_UpdateClusterDNSTasksDone(t *testing.T) {
	if _, err := dbs.Default(); err != nil {
		t.Skipf("skip: %v", err)
	}
	dao, err := dns.NewDNSTaskDAO()
	if err != nil {
		t.Fatal(err)
	}
	var tx *dbs.Tx
	err = dao.UpdateClusterDNSTasksDone(tx, 46, time.Now().UnixNano())
	if err != nil {
		t.Fatal(err)
	}
}
