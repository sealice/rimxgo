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
	keyPrefix := "session"
	conf = sessions.Config{
		Cookie:                      config.GetStringDefault(keyPrefix+".cookie", sessions.DefaultCookieName),
		Expires:                     config.GetDuration(keyPrefix+".expires") * time.Second,
		CookieSecureTLS:             config.GetBool(keyPrefix + ".cookieSecureTLS"),
		DisableSubdomainPersistence: config.GetBool(keyPrefix + ".disableSubdomainPersistence"),
		AllowReclaim:                config.GetBoolDefault(keyPrefix+".allowReclaim", true),
	}
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
