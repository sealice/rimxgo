package controllers

import (
	"github.com/rimxgo/constant"
	"github.com/rimxgo/models"
)

type PerRoleController struct {
	baseController
}

// router get /list?query={}
func (c *PerRoleController) GetList() *Result {
	v := &models.PerRole{}

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

// router get /page?pageNum=1&pageSize=20?query={}
func (c *PerRoleController) GetPage() *Result {
	v := &models.PerRole{}
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

// router get /get?id=1
func (c *PerRoleController) GetGet() *Result {
	id, err := c.Ctx.URLParamInt("id")
	if id == 0 || err != nil {
		logger.Debug("解析参数错误，", err)
		return RetResult(constant.CodeBusinessError, "解析参数错误")
	}

	v := &models.PerRole{Id: id}
	if has, err := v.GetOne(); err != nil {
		logger.Error(err)
		return RetResult(constant.CodeBusinessError, "系统繁忙")
	} else if !has {
		return RetResult(constant.CodeBusinessError, "记录不存在")
	}

	return RetResultData(v)
}

// router post /add
func (c *PerRoleController) PostAdd() *Result {
	v := &models.PerRole{}
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

	return RetResultData(v, "ok")
}

// router post /edit
func (c *PerRoleController) PostEdit() *Result {
	v := &models.PerRouter{}
	if err := c.Ctx.ReadJSON(v); err != nil {
		logger.Debug("解析参数错误，", err)
		return RetResult(constant.CodeBusinessError, "解析参数错误")
	}

	if err := v.Validate(); err != nil {
		logger.Debug("校验参数错误，", err)
		return RetResult(constant.CodeBusinessError, err.Error())
	}

	if _, err := v.Update(); err != nil {
		logger.Error(err)
		return RetResult(constant.CodeBusinessError, "系统繁忙")
	}

	return RetResultData(nil, "ok")
}

// router get /del?id=1
func (c *PerRoleController) GetDel() *Result {
	id, err := c.Ctx.URLParamInt("id")
	if err != nil {
		logger.Debug("解析参数错误，", err)
		return RetResult(constant.CodeBusinessError, "解析参数错误")
	}

	v := &models.PerRole{Id: id}
	if _, err := v.Delete(); err != nil {
		logger.Error(err)
		return RetResult(constant.CodeBusinessError, "系统繁忙")
	}

	return RetResultData(nil, "ok")
}
