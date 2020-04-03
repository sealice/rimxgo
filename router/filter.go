package router

import (
	"github.com/kataras/iris/v12"
	"github.com/rimxgo/constant"
	"github.com/rimxgo/controllers"
	"github.com/rimxgo/middleware/session"
)

func auth(ctx iris.Context) {
	s := session.Instance().Start(ctx)
	if user := s.Get(constant.SESSION_KEY_USER); user == nil {
		ctx.JSON(controllers.Result{Code: 100, Msg: "未登录或登录已过期"})
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.StopExecution()
		return
	}
	ctx.Next()
}
