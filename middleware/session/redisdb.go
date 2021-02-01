package session

import (
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"github.com/kataras/iris/v12/sessions/sessiondb/redis"
	"github.com/rimxgo/config"
)

func useRedisdb(sess *sessions.Sessions) {
	keyPrefix := "redis"
	driver := redis.Redigo()
	db := redis.New(redis.Config{
		Network:   config.GetString(keyPrefix + ".network"),
		Addr:      config.GetString(keyPrefix + ".addr"),
		Password:  config.GetString(keyPrefix + ".password"),
		Database:  config.GetString(keyPrefix + ".database"),
		MaxActive: config.GetInt(keyPrefix + ".poolsize"),
		Prefix:    strings.ToUpper(config.GetStringDefault("name", "iris") + ":SESSION:"),
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
