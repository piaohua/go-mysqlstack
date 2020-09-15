package main

import (
	"fmt"

	"github.com/xelabs/go-mysqlstack/driver"
	"github.com/xelabs/go-mysqlstack/xlog"
)

func main() {
	log := xlog.NewStdLog(xlog.Level(xlog.INFO))
	address := fmt.Sprintf(":4408")
	client, err := driver.NewConn("proxy", "proxy", address, "", "")
	if err != nil {
		log.Panic("client.new.connection.error:%+v", err)
	}
	defer client.Close()

	qr, err := client.FetchAll("SELECT * FROM PROXY", -1)
	if err != nil {
		log.Panic("client.query.error:%+v", err)
	}
	log.Info("results:[%+v]", qr.Rows)

	qr, err = client.FetchAll("SELECT id FROM PROXY", -1)
	if err != nil {
		log.Panic("client.query.error:%+v", err)
	}
	log.Info("results:[%+v]", qr.Rows)
}
