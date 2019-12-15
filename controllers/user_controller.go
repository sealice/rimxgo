package controllers

import (
	"github.com/rimxgo/models"
)

type UserController struct {
	baseController
}

// router get /list
func (c *UserController) GetList() *Result {
	v := &models.User{}
	p, _ := c.ReadPage()

	if err := c.ReadQueryJSON(&v.Query); err != nil {
		logger.Debug("参数错误，", err)
		return c.RetResult(2998, "参数错误")
	}

	ls, err := v.GetList(p)
	if err != nil {
		logger.Error(err)
		return c.RetResult(2998, "系统繁忙")
	}

	return c.RetResultList(ls, p)
}

// router get /get
func (c *UserController) GetGet() *Result {
	id, err := c.Ctx.URLParamInt("id")
	if err != nil {
		logger.Debug("参数错误，", err)
		return c.RetResult(2998, "参数错误")
	}

	v := &models.User{Id: id}
	if has, err := v.GetOne(); err != nil {
		logger.Error(err)
		return c.RetResult(2998, "系统繁忙")
	} else if !has {
		return c.RetResult(2998, "用户不存在")
	}

	return c.RetResultData("", v)
}

// router post /add
func (c *UserController) PostAdd() *Result {
	v := &models.User{}
	if err := c.Ctx.ReadJSON(v); err != nil {
		logger.Debug("参数错误，", err)
		return c.RetResult(2998, "参数错误")
	}

	if _, err := v.Insert(); err != nil {
		logger.Error(err)
		return c.RetResult(2998, "系统繁忙")
	}

	return c.RetResultData("ok")
}