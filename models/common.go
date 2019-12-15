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
	engine = initEngine("mysql", "MysqlDefault")
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

	if config.GetBool("Debug") {
		engine.ShowSQL()
	}

	return engine
}

func loadConfig(configKey string) (uri string, maxIdle, maxOpen int) {
	v := config.Sub(configKey)
	uri = v.GetString("Uri")
	maxIdle = v.GetInt("MaxIdleConns")
	maxOpen = v.GetInt("MaxOpenConns")
	return
}
