package models

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/kataras/iris/v12"
	"github.com/rimxgo/config"
)

var engine *xorm.Engine

func init() {
	engine = initEngine("mysql", "mysqlDefault")
}

func initEngine(driverName, configKey string) *xorm.Engine {
	uri, maxIdle, maxOpen := loadConfig(configKey)

	fmt.Println("Connect to the database", configKey)
	engine, err := xorm.NewEngine(driverName, uri)
	if err != nil {
		panic(err)
	}

	iris.RegisterOnInterrupt(func() {
		engine.Close()
	})

	if err = engine.Ping(); err != nil {
		panic(err)
	}

	if maxIdle > 0 {
		engine.SetMaxIdleConns(maxIdle)
	}

	if maxOpen > 0 {
		engine.SetMaxOpenConns(maxOpen)
	}

	if config.GetBool("debug") {
		engine.ShowSQL() // 调试模式打印SQL日志
	}

	return engine
}

func loadConfig(keyPrefix string) (uri string, maxIdle, maxOpen int) {
	uri = config.GetString(keyPrefix + ".uri")
	maxIdle = config.GetInt(keyPrefix + ".maxIdleConns")
	maxOpen = config.GetInt(keyPrefix + ".maxOpenConns")
	return
}
