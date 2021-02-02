package session

import (
	"net/http"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"github.com/rimxgo/config"
	"github.com/rimxgo/constant"
	"github.com/rimxgo/helper/logs"
)

const SK = "_iniris"
const SKV = "1"

var Conf sessions.Config
var sess *sessions.Sessions
var storageType string

func init() {
	keyPrefix := "session"
	storageType = config.GetString(keyPrefix + ".storage")
	Conf = sessions.Config{
		Cookie:                      config.GetStringDefault(keyPrefix+".cookie", sessions.DefaultCookieName),
		Expires:                     config.GetDuration(keyPrefix+".expires") * time.Second,
		CookieSecureTLS:             config.GetBool(keyPrefix + ".cookieSecureTLS"),
		DisableSubdomainPersistence: config.GetBool(keyPrefix + ".disableSubdomainPersistence"),
		AllowReclaim:                config.GetBoolDefault(keyPrefix+".allowReclaim", true),
	}

	sess = sessions.New(Conf)
}

func Register(m *mvc.Application) {
	if storageType == "redis" {
		useRedisdb(sess)
	} else if storageType == "file" {
		useBoltdb(sess)
	}

	m.Register(func(ctx iris.Context) *sessions.Session {
		if ctx.GetCookie(Conf.Cookie) != "" {
			s := sess.Start(ctx)
			if s.Get(constant.SESSION_KEY_USER) != nil {
				shiftExpiration(s, ctx)
				if ctx.GetCookie(SK) != SKV {
					cookie := &http.Cookie{
						Name:     SK,
						Value:    SKV,
						Path:     "/",
						HttpOnly: true,
						SameSite: http.SameSiteLaxMode,
					}
					ctx.SetCookie(cookie)
				}
			}
			return s
		}

		return nil
	})
}

func Instance() *sessions.Sessions {
	return sess
}

// Session续期
func shiftExpiration(s *sessions.Session, ctx iris.Context) {
	if Conf.Expires > 0 &&
		!s.Lifetime.HasExpired() &&
		s.Lifetime.DurationUntilExpiration() < Conf.Expires/5 {
		// 快过期前更新Session
		if err := sess.ShiftExpiration(ctx); err != nil {
			logs.Logger.Warn("Update Session: ", err)
		}
	}
}
