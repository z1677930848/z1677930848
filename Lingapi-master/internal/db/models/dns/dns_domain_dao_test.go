package dns

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iwind/TeaGo/dbs"
)

func TestDNSDomainDAO_ExistDomainRecord(t *testing.T) {
	if _, err := dbs.Default(); err != nil {
		t.Skipf("skip: %v", err)
	}
	var tx *dbs.Tx

	dao, err := NewDNSDomainDAO()
	if err != nil {
		t.Fatal(err)
	}

	{
		b, err := dao.ExistDomainRecord(tx, 1, "mycluster", "A", "", "")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(b)
	}
	{
		b, err := dao.ExistDomainRecord(tx, 2, "mycluster", "A", "", "")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(b)
	}
	{
		b, err := dao.ExistDomainRecord(tx, 2, "mycluster", "MX", "", "")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(b)
	}
	{
		b, err := dao.ExistDomainRecord(tx, 2, "mycluster123", "A", "", "")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(b)
	}
}
