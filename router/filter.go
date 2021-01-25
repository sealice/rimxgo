package router

import (
	"github.com/kataras/iris/v12"
	"github.com/rimxgo/constant"
	"github.com/rimxgo/controllers"
	"github.com/rimxgo/middleware/session"
)

func auth(ctx iris.Context) {
	if ctx.GetCookie(session.Conf.Cookie) == "" {
		ctx.StatusCode(iris.StatusUnauthorized)

		if ctx.GetCookie(session.SK) != session.SKV {
			ctx.JSON(controllers.Result{Code: 100, Msg: "未登录，请先登录"})
		} else {
			ctx.JSON(controllers.Result{Code: 100, Msg: "登录已过期，请重新登录"})
		}

		ctx.StopExecution()
		return
	}

	s := session.Instance().Start(ctx)
	if s.Get(constant.SESSION_KEY_USER) == nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(controllers.Result{Code: 100, Msg: "登录已失效，请重新登录"})
		ctx.StopExecution()
		return
	}

	ctx.Next()
}
