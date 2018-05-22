package mysql

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

const testConn = "root:toor@tcp(127.0.0.1:3306)/newsapitest?charset=utf8&parseTime=True&loc=Local"

func TestNew(t *testing.T) {
	handler, err := New(testConn)
	if err != nil {
		t.Error(err)
	}
	handler.Migrate()
}
