package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/rimxgo/config"
	"github.com/rimxgo/router"
)

func main() {
	app := iris.New()

	app.Logger().SetLevel(config.GetString("logLevel"))

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Status:            true,
		IP:                true,
		Method:            true,
		Path:              true,
		Query:             true,
		MessageHeaderKeys: []string{"User-Agent"},
	}))

	// Register router
	router.Register(app)

	app.Run(
		iris.Addr(config.GetString("prot")),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
