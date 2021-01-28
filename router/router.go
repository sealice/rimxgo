package router

import (
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/rimxgo/constant"
	"github.com/rimxgo/controllers"
	"github.com/rimxgo/middleware/session"
	"github.com/rimxgo/models"
	"github.com/rimxgo/router/filter"
)

func Register(app *iris.Application) {
	m := mvc.Configure(app.Party("/"), configure)

	// router
	m.Handle(&controllers.DefaultController{})
	api := m.Party("/v1")
	{
		api.Handle(&controllers.MainController{})
		api.Party("/user", filter.Authentication).Handle(&controllers.UserController{})
	}
}

func configure(m *mvc.Application) {
	// session.Register(m, redisdb.Use)
	session.Register(m)

	m.Register(
		time.Now(),
		func(ctx iris.Context) *models.User {
			if ctx.GetCookie(session.Conf.Cookie) != "" {
				s := session.Instance().Start(ctx)
				user, _ := s.Get(constant.SESSION_KEY_USER).(*models.User)
				return user
			}

			return nil
		},
	)
}
