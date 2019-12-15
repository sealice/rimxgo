package redis

import (
	"github.com/go-redis/redis/v7"
	"github.com/kataras/iris/v12"
	"github.com/rimxgo/config"
)

var Client *redis.Client

func init() {
	Client = NewClient("Redis")
}

func NewClient(configKey string) *redis.Client {
	v := config.Sub(configKey)
	client := redis.NewClient(&redis.Options{
		Network:  v.GetString("Network"),
		Addr:     v.GetString("Addr"),
		Password: v.GetString("Password"),
		DB:       v.GetInt("Database"),
		PoolSize: v.GetInt("PoolSize"),
	})

	iris.RegisterOnInterrupt(func() {
		client.Close()
	})

	if _, err := client.Ping().Result(); err != nil {
		panic(err)
	}

	return client
}
