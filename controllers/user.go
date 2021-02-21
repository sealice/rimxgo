package controllers

import (
	"github.com/rimxgo/constant"
	"github.com/rimxgo/models"
)

type UserController struct {
	baseController
}

// router get /
func (c *UserController) Get() *Result {
	v := *c.User
	v.Password = ""
	return RetResultData(v)
}

// router get /list?query={}
func (c *UserController) GetList() *Result {
	v := &models.User{}

	if err := c.ReadQueryJSON(&v.Query); err != nil {
		logger.Debug("解析参数错误，", err)
		return RetResult(constant.CodeBusinessError, "解析参数错误")
	}

	ls, err := v.GetList()
	if err != nil {
		logger.Error(err)
		return RetResult(constant.CodeBusinessError, "系统繁忙")
	}

	return RetResultList(ls)
}

// router get /page
func (c *UserController) GetPage() *Result {
	v := &models.User{}
	p, _ := c.ReadPage()

	if err := c.ReadQueryJSON(&v.Query); err != nil {
		logger.Debug("解析参数错误，", err)
		return RetResult(constant.CodeBusinessError, "解析参数错误")
	}

	ls, err := v.GetPage(p)
	if err != nil {
		logger.Error(err)
		return RetResult(constant.CodeBusinessError, "系统繁忙")
	}

	return RetResultList(ls, p)
}

// router get /get
func (c *UserController) GetGet() *Result {
	id, err := c.Ctx.URLParamInt("id")
	if err != nil {
		logger.Debug("解析参数错误，", err)
		return RetResult(constant.CodeBusinessError, "解析参数错误")
	}

	v := &models.User{Id: id}
	if has, err := v.GetOne(); err != nil {
		logger.Error(err)
		return RetResult(constant.CodeBusinessError, "系统繁忙")
	} else if !has {
		return RetResult(constant.CodeBusinessError, "用户不存在")
	}

	return RetResultData(v)
}

// router post /add
func (c *UserController) PostAdd() *Result {
	v := &models.User{}
	if err := c.Ctx.ReadJSON(v); err != nil {
		logger.Debug("解析参数错误，", err)
		return RetResult(constant.CodeBusinessError, "解析参数错误")
	}

	if err := v.Validate(); err != nil {
		logger.Debug("校验参数错误，", err)
		return RetResult(constant.CodeBusinessError, err.Error())
	}

	if _, err := v.Insert(); err != nil {
		logger.Error(err)
		return RetResult(constant.CodeBusinessError, "系统繁忙")
	}

	return RetResultData(nil, "ok")
}
