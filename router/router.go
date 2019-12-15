package router

import (
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/rimxgo/controllers"
	"github.com/rimxgo/middleware/session"
)

func Register(app *iris.Application) {
	m := mvc.Configure(app.Party("/"), configure)

	// router
	m.Handle(&controllers.DefaultController{})
	api := m.Party("/v1")
	{
		api.Party("/user").Handle(&controllers.UserController{})
	}
}

func configure(m *mvc.Application) {
	// session.Register(m, redisdb.Use)
	session.Register(m)
	m.Register(time.Now())
}

func auth(ctx iris.Context) {
	s := session.Instance().Start(ctx)
	if ok, _ := s.GetBoolean("isLogin"); !ok {
		// ctx.JSON(controllers.Result{Code: 100, Msg: "未登录或登录已过期"})
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.StopExecution()
		return
	}
	ctx.Next()
}
