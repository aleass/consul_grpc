package common

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
	"time"
)

var MainDbEngine *xorm.Engine

func init() {
	var err error
	DsName := fmt.Sprintf("%v:%v@(%v)/%v?charset=%v", "root", "root", "127.0.0.1:3306", "swhc", "utf8mb4")
	fmt.Println(DsName)
	MainDbEngine, err = xorm.NewEngine("mysql", DsName)
	if err != nil {
		log.Println(err)
	}
	MainDbEngine.ShowSQL(true)
	MainDbEngine.SetMaxOpenConns(100)
	MainDbEngine.SetMaxIdleConns(100)
	MainDbEngine.SetConnMaxLifetime(240*time.Second)
	//匹配表名
	MainDbEngine.Sync2()
}
