package filter

import (
	"github.com/kataras/iris/v12"
	"github.com/rimxgo/constant"
	"github.com/rimxgo/controllers"
	"github.com/rimxgo/middleware/session"
	"github.com/rimxgo/models"
)

// 身份验证（需要登录）
func Authentication(ctx iris.Context) {
	if ctx.GetCookie(session.Conf.Cookie) == "" {
		// ctx.StatusCode(iris.StatusUnauthorized) // 设置状态为401
		if ctx.GetCookie(session.SK) != session.SKV {
			ctx.JSON(controllers.RetResultError(constant.CodeNotLoggedIn))
		} else {
			ctx.JSON(controllers.RetResultError(constant.CodeLoginExpired))
		}

		ctx.StopExecution()
		return
	}

	s := session.Instance().Start(ctx)
	if s.Get(constant.SESSION_KEY_USER) == nil {
		// ctx.StatusCode(iris.StatusUnauthorized) // 设置状态为401
		ctx.JSON(controllers.RetResultError(constant.CodeLoginInvalid))
		ctx.StopExecution()
		return
	}

	ctx.Next()
}

// 权限验证（需先验证身份，即在过滤器 Authentication 后使用）
func Permission(ctx iris.Context) {
	s := session.Instance().Start(ctx)
	user, ok := s.Get(constant.SESSION_KEY_USER).(*models.User)
	if !ok {
		// ctx.StatusCode(iris.StatusUnauthorized) // 设置状态为401
		ctx.JSON(controllers.RetResultError(constant.CodeLoginInvalid))
		ctx.StopExecution()
		return
	}

	allow := false
	currentPath := ctx.GetCurrentRoute().Path()
	for _, path := range user.Routers {
		if path == currentPath {
			allow = true
			break
		}
	}

	if !allow {
		// ctx.StatusCode(iris.StatusForbidden) // 设置状态为403
		ctx.JSON(controllers.RetResultError(constant.CodeNoPermission))
		ctx.StopExecution()
		return
	}

	ctx.Next()
}
