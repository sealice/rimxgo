package redis

import (
	"github.com/go-redis/redis/v7"
	"github.com/kataras/iris/v12"
	"github.com/rimxgo/config"
)

var Client *redis.Client

func init() {
	Client = NewClient("redis")
}

func NewClient(keyPrefix string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Network:  config.GetString(keyPrefix + ".network"),
		Addr:     config.GetString(keyPrefix + ".addr"),
		Password: config.GetString(keyPrefix + ".password"),
		DB:       config.GetInt(keyPrefix + ".database"),
		PoolSize: config.GetInt(keyPrefix + ".poolsize"),
	})

	iris.RegisterOnInterrupt(func() {
		client.Close()
	})

	if _, err := client.Ping().Result(); err != nil {
		panic(err)
	}

	return client
}
