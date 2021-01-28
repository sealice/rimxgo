package controllers

import (
	"github.com/rimxgo/constant"
	"github.com/rimxgo/middleware/session"
	"github.com/rimxgo/models"
)

type MainController struct {
	baseController
}

func (c *MainController) PostRegister() *Result {
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

func (c *MainController) PostLogin() *Result {
	v := &models.User{}

	if err := c.Ctx.ReadJSON(v); err != nil {
		logger.Debug("解析参数错误，", err)
		return RetResult(constant.CodeBusinessError, "解析参数错误")
	}

	if err := v.Validate(); err != nil {
		logger.Debug("校验参数错误，", err)
		return RetResult(constant.CodeBusinessError, err.Error())
	}

	user := &models.User{Name: v.Name, Password: v.Password}
	if has, err := user.GetOne(); err != nil {
		logger.Error(err)
		return RetResult(constant.CodeBusinessError, "系统繁忙")
	} else if !has {
		logger.Debugf("用户`%s`不存在或密码错误", user.Name)
		return RetResult(constant.CodeBusinessError, "用户不存在或密码错误")
	}

	*v = *user
	v.Password = ""

	sess := session.Instance().Start(c.Ctx)
	sess.Set(constant.SESSION_KEY_USER, user)
	logger.Infof("用户`%s`登录成功", v.Name)

	return RetResultData(v)
}

func (c *MainController) GetLogout() *Result {
	if c.User != nil {
		c.Session.Delete(constant.SESSION_KEY_USER)
	}
	return RetResultData(nil)
}
