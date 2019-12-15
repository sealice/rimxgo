package controllers

import (
	"encoding/json"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"github.com/rimxgo/helper/logs"
	"github.com/rimxgo/models/vos"
)

var logger = logs.Logger

type baseController struct {
	Ctx       iris.Context
	Session   *sessions.Session
	StartTime time.Time
}

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
	List interface{} `json:"list,omitempty"`
	*vos.Page
}

func (c *baseController) RetResult(code int, msg string, data ...interface{}) *Result {
	if len(data) == 0 {
		data = append(data, nil)
	}
	return &Result{Code: code, Msg: msg, Data: data[0]}
}

func (c *baseController) RetResultData(msg string, data ...interface{}) *Result {
	return c.RetResult(0, msg, data...)
}

func (c *baseController) RetResultList(list interface{}, page ...*vos.Page) *Result {
	if len(page) == 0 {
		page = append(page, nil)
	}
	return &Result{Code: 0, List: list, Page: page[0]}
}

func (c *baseController) ReadPage() (*vos.Page, error) {
	page := &vos.Page{PageNum: 1, PageSize: 10}
	err := c.Ctx.ReadForm(page)
	page.Start = (page.PageNum - 1) * page.PageSize
	return page, err
}

func (c *baseController) ReadQueryJSON(query interface{}) error {
	return json.Unmarshal([]byte(c.Ctx.FormValueDefault("query", "{}")), query)
}
