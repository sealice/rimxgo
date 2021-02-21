package controllers

import (
	"fmt"
	"time"

	"github.com/rimxgo/constant"
	"github.com/rimxgo/middleware/session"
	"github.com/rimxgo/models"
)

type DefaultController struct {
	baseController
}

func (c *DefaultController) Get() string {
	visits := 1
	since := time.Since(c.StartTime).Seconds()
	logger.Debugf("login user: %#v", c.User)

	// write the current, updated visits.
	if c.Session != nil {
		logger.Infof("sessionid: %s", c.Session.ID())
		visits = c.Session.Increment("visits", 1)
	}

	return fmt.Sprintf("%d visit(s) from my current session in %0.1f seconds of server's up-time",
		visits, since)
}

func (c *DefaultController) PostRegister() *Result {
	v := &models.User{}

	if err := c.Ctx.ReadJSON(v); err != nil {
		logger.Debug("解析参数错误，", err)
		return RetResult(constant.CodeBusinessError, "解析参数错误")
	}

	if err := v.Validate(); err != nil {
		logger.Debug("校验参数错误，", err)
		return RetResult(constant.CodeBusinessError, err.Error())
	}

	user := &models.User{Name: v.Name}
	if has, _ := user.IsExist(); has {
		logger.Debugf("用户`%s`已存在", user.Name)
		return RetResult(constant.CodeBusinessError, "用户已存在")
	}

	if _, err := v.Insert(); err != nil {
		logger.Error("注册入库操作失败，", err)
		return RetResult(constant.CodeBusinessError, "系统繁忙")
	}

	v.Password = ""
	logger.Infof("用户`%s`注册成功", v.Name)
	return RetResultData(v)
}

func (c *DefaultController) PostLogin() *Result {
	v := &models.User{}

	if err := c.Ctx.ReadJSON(v); err != nil {
		logger.Debug("解析参数错误，", err)
		return RetResult(constant.CodeBusinessError, "解析参数错误")
	}

	if err := v.Validate(); err != nil {
		logger.Debug("校验参数错误，", err)
		return RetResult(constant.CodeBusinessError, err.Error())
	}

	if has, err := v.GetOne(); err != nil {
		logger.Error(err)
		return RetResult(constant.CodeBusinessError, "系统繁忙")
	} else if !has {
		if has, _ = (&models.User{Name: v.Name}).IsExist(); !has {
			logger.Debugf("用户`%s`不存在", v.Name)
			return RetResult(constant.CodeBusinessError, "用户不存在")
		}

		logger.Debugf("用户`%s`密码输入错误", v.Name)
		return RetResult(constant.CodeBusinessError, "密码错误")
	}

	if err := v.GetPermission(); err != nil {
		logger.Error("获取用户权限错误，", err)
		return RetResult(constant.CodeBusinessError, "系统繁忙")
	}

	user := *v
	sess := session.Instance().Start(c.Ctx)
	sess.Set(constant.SESSION_KEY_USER, &user)
	logger.Infof("用户`%s`登录成功", v.Name)
	logger.Debugf("用户`%s`访问权限，%v", v.Name, v.Routers)

	v.Password = ""
	return RetResultData(v)
}

func (c *DefaultController) GetLogout() *Result {
	if c.User != nil {
		c.Session.Delete(constant.SESSION_KEY_USER)
		c.Ctx.RemoveCookie(session.Conf.Cookie)
		c.Ctx.RemoveCookie(session.SK)
	}

	return RetResultData(nil)
}
