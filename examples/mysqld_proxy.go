package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/xelabs/go-mysqlstack/driver"
	querypb "github.com/xelabs/go-mysqlstack/sqlparser/depends/query"
	"github.com/xelabs/go-mysqlstack/sqlparser/depends/sqltypes"
	"github.com/xelabs/go-mysqlstack/xlog"
)

func main() {
	result1 := &sqltypes.Result{
		Fields: []*querypb.Field{
			{
				Name: "id",
				Type: querypb.Type_INT32,
			},
			{
				Name: "name",
				Type: querypb.Type_VARCHAR,
			},
		},
		Rows: [][]sqltypes.Value{
			{
				sqltypes.MakeTrusted(querypb.Type_INT32, []byte("10")),
				sqltypes.MakeTrusted(querypb.Type_VARCHAR, []byte("nice name")),
			},
		},
	}

	log := xlog.NewStdLog(xlog.Level(xlog.INFO))
	th := driver.NewProxyHandler(log)
	th.AddQuery("SELECT * FROM MOCK", result1)

	result2 := &sqltypes.Result{
		Fields: []*querypb.Field{
			{
				Name: "id",
				Type: querypb.Type_INT32,
			},
		},
		Rows: [][]sqltypes.Value{
			{
				sqltypes.MakeTrusted(querypb.Type_INT32, []byte("10")),
			},
		},
	}
	th.AddQueryDelay("SELECT id FROM MOCK", result2, 300)

	mysqld, err := driver.ProxyMysqlServerWithPort(log, 4408, th)
	if err != nil {
		log.Panic("mysqld.start.error:%+v", err)
	}
	defer mysqld.Close()
	log.Info("mysqld.server.start.address[%v]", mysqld.Addr())

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
