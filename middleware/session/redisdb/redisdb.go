package redisdb

import (
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"github.com/kataras/iris/v12/sessions/sessiondb/redis"
	"github.com/rimxgo/config"
)

func Use(sess *sessions.Sessions) {
	v := config.Sub("Redis")
	driver := redis.Redigo()
	db := redis.New(redis.Config{
		Network:   v.GetString("Network"),
		Addr:      v.GetString("Addr"),
		Password:  v.GetString("Password"),
		Database:  v.GetString("Database"),
		MaxActive: v.GetInt("PoolSize"),
		Prefix:    strings.ToUpper(config.GetStringDefault("Name", "iris") + ":SESSION:"),
		Driver:    driver,
	})

	iris.RegisterOnInterrupt(func() {
		db.Close()
	})

	if _, err := driver.PingPong(); err != nil {
		panic(err)
	}

	sess.UseDatabase(db)
}
