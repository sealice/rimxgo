package controllers

import (
	"encoding/json"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"github.com/rimxgo/constant"
	"github.com/rimxgo/helper/logs"
	"github.com/rimxgo/models"
	"github.com/rimxgo/models/vos"
)

var logger = logs.Logger

type baseController struct {
	Ctx       iris.Context
	StartTime time.Time
	Session   *sessions.Session
	User      *models.User
}

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
	List interface{} `json:"list,omitempty"`

	*vos.Page
}

// 返回自定义（正确/错误）码的数据
func RetResult(code int, msg string, data ...interface{}) *Result {
	if len(data) == 0 {
		data = append(data, nil)
	}

	return &Result{Code: code, Msg: msg, Data: data[0]}
}

// 返回成功的数据
func RetResultData(data interface{}, msg ...string) *Result {
	if len(msg) == 0 {
		msg = append(msg, "")
	}

	return &Result{Code: 0, Msg: msg[0], Data: data}
}

// 返回成功的列表/分页数据
func RetResultList(list interface{}, page ...*vos.Page) *Result {
	if len(page) == 0 {
		page = append(page, nil)
	}

	return &Result{Code: 0, List: list, Page: page[0]}
}

// 返回已定义的错误码
func RetResultError(code int, data ...interface{}) *Result {
	msg := constant.CodeText(code)
	if msg == "" {
		msg = "系统繁忙"
	}

	if len(data) == 0 {
		data = append(data, nil)
	}

	return &Result{Code: code, Msg: msg, Data: data[0]}
}

func (c *baseController) ReadPage() (*vos.Page, error) {
	page := &vos.Page{}
	err := c.Ctx.ReadForm(page)

	if page.PageNum <= 0 {
		page.PageNum = 1
	}

	if page.PageSize <= 0 {
		page.PageSize = 20
	}

	page.Start = (page.PageNum - 1) * page.PageSize

	return page, err
}

func (c *baseController) ReadQueryJSON(query interface{}) error {
	return json.Unmarshal([]byte(c.Ctx.FormValueDefault("query", "{}")), query)
}
