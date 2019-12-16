package router

import (
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/rimxgo/config"
	"github.com/rimxgo/controllers"
	"github.com/rimxgo/middleware/session"
	"github.com/rimxgo/models"
)

func Register(app *iris.Application) {
	m := mvc.Configure(app.Party("/"), configure)

	// router
	m.Handle(&controllers.DefaultController{})
	api := m.Party("/v1")
	{
		api.Handle(&controllers.MainController{})
		api.Party("/user", auth).Handle(&controllers.UserController{})
	}
}

func configure(m *mvc.Application) {
	// session.Register(m, redisdb.Use)
	session.Register(m)
	m.Register(
		time.Now(),
		func(ctx iris.Context) *models.User {
			s := session.Instance().Start(ctx)
			user, _ := s.Get(config.SESSION_KEY_USER).(*models.User)
			return user
		},
	)
}

func auth(ctx iris.Context) {
	s := session.Instance().Start(ctx)
	if user := s.Get(config.SESSION_KEY_USER); user == nil {
		ctx.JSON(controllers.Result{Code: 100, Msg: "未登录或登录已过期"})
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.StopExecution()
		return
	}
	ctx.Next()
}
