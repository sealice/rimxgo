package router

import (
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/rimxgo/constant"
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
			user, _ := s.Get(constant.SESSION_KEY_USER).(*models.User)
			return user
		},
	)
}
