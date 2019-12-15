package session

import (
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"github.com/rimxgo/config"
	"github.com/rimxgo/helper/logs"
)

var (
	conf sessions.Config
	sess *sessions.Sessions
)

func init() {
	v := config.Sub("Session")
	conf = (sessions.Config{
		Cookie:                      v.GetString("Cookie"),
		Expires:                     v.GetDuration("Expires") * time.Second,
		CookieSecureTLS:             v.GetBool("CookieSecureTLS"),
		DisableSubdomainPersistence: v.GetBool("DisableSubdomainPersistence"),
		AllowReclaim:                true,
	}).Validate()
	sess = sessions.New(conf)
}

func Register(m *mvc.Application, db ...func(sess *sessions.Sessions)) {
	if len(db) > 0 {
		db[0](sess)
	}

	m.Register(sess.Start)

	m.Router.Use(func(ctx iris.Context) {
		if ctx.GetCookie(conf.Cookie) != "" && conf.Expires > 0 {
			s := sess.Start(ctx)
			if !s.Lifetime.HasExpired() &&
				s.Lifetime.DurationUntilExpiration() < conf.Expires/5 {
				// 快过期前更新Session
				if err := sess.ShiftExpiration(ctx); err != nil {
					logs.Logger.Warn("Update Session: ", err)
				}
			}
		}

		ctx.Next()
	})
}

func Instance() *sessions.Sessions {
	return sess
}
