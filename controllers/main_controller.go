package controllers

import "github.com/rimxgo/models"

import "github.com/rimxgo/config"

type MainController struct {
	baseController
}

func (c *MainController) PostRegister() *Result {
	v := &models.User{}

	if err := c.Ctx.ReadJSON(v); err != nil {
		logger.Debug("参数错误，", err)
		return c.RetResult(2998, "参数错误")
	}

	if err := v.Validate(); err != nil {
		logger.Debug("参数错误，", err)
		return c.RetResult(2998, err.Error())
	}

	user := &models.User{Name: v.Name}
	if has, _ := user.IsExist(); has {
		logger.Debugf("用户`%s`已存在", user.Name)
		return c.RetResult(2998, "用户已存在")
	}

	if _, err := v.Insert(); err != nil {
		logger.Error(err)
		return c.RetResult(2998, "系统繁忙")
	}

	v.Password = ""
	logger.Infof("用户`%s`注册成功", v.Name)
	return c.RetResultData("", v)
}

func (c *MainController) PostLogin() *Result {
	v := &models.User{}

	if err := c.Ctx.ReadJSON(v); err != nil {
		logger.Debug("参数错误，", err)
		return c.RetResult(2998, "参数错误")
	}

	if err := v.Validate(); err != nil {
		logger.Debug("参数错误，", err)
		return c.RetResult(2998, err.Error())
	}

	user := &models.User{Name: v.Name, Password: v.Password}
	if has, err := user.GetOne(); err != nil {
		logger.Error(err)
		return c.RetResult(2998, "系统繁忙")
	} else if !has {
		logger.Debugf("用户`%s`不存在或密码错误", user.Name)
		return c.RetResult(2998, "用户不存在或密码错误")
	}

	*v = *user
	v.Password = ""
	c.Session.Set(config.SESSION_KEY_USER, v)
	logger.Infof("用户`%s`登录成功", v.Name)
	return c.RetResultData("", v)
}

func (c *MainController) GetLogout() *Result {
	if c.User != nil {
		c.Session.Delete(config.SESSION_KEY_USER)
	}
	return c.RetResultData("")
}
