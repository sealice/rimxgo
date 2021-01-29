package models

import (
	"strings"

	"github.com/go-xorm/xorm"
	"github.com/rimxgo/models/vos"
)

type PerElement struct {
	Id         int          `json:"id" xorm:"not null pk autoincr INT(11)"`
	RouterIds  string       `json:"routerIds,omitempty" xorm:"VARCHAR(255)"`
	Name       string       `json:"name" xorm:"VARCHAR(255)"`
	Code       string       `json:"code" xorm:"VARCHAR(255)"`
	Remark     string       `json:"remark,omitempty" xorm:"VARCHAR(255)"`
	CreateTime vos.JsonTime `json:"createTime" xorm:"DATETIME created"`
	UpdateTime vos.JsonTime `json:"updateTime" xorm:"DATETIME updated"`
	DeleteTime vos.JsonTime `json:"-" xorm:"DATETIME deleted"`
	Query      vos.Query    `json:"-" xorm:"-"` // 查询条件
	Routers    []*PerRouter `json:"routers,omitempty" xorm:"-"`
}

func (m *PerElement) GetList(oCols ...string) ([]*PerElement, error) {
	ls := make([]*PerElement, 0)
	sess := m.listSess(oCols...)

	err := sess.Iterate(new(PerElement), func(i int, bean interface{}) error {
		item := bean.(*PerElement)
		item.Query = vos.Query{vos.DeepKey: m.Query[vos.DeepKey]}
		item.inRouters()
		ls = append(ls, item)
		return nil
	})

	return ls, err
}

func (m *PerElement) GetPage(p *vos.Page) ([]*PerElement, error) {
	ls := make([]*PerElement, 0)
	sess := m.listSess()

	var err error
	p.Total, err = sess.Limit(p.PageSize, p.Start).FindAndCount(&ls)
	return ls, err
}

func (m *PerElement) GetOne() (bool, error) {
	has, err := engine.Get(m)
	if has && err == nil {
		m.Query = vos.Query{vos.DeepKey: true}
		err = m.inRouters()
	}

	return has, err
}

func (m *PerElement) Insert() (int64, error) {
	return engine.AllCols().InsertOne(m)
}

func (m *PerElement) Update() (int64, error) {
	return engine.ID(m.Id).MustCols("remark").Update(m)
}

func (m *PerElement) Delete() (int64, error) {
	return engine.ID(m.Id).Delete(new(PerElement))
}

func (m *PerElement) IsExist() (bool, error) {
	return engine.Exist(m)
}

func (m *PerElement) Validate() error {
	return nil
}

func (m *PerElement) listSess(oCols ...string) *xorm.Session {
	return engine.Table(m).Omit(oCols...).Where(vos.SetQueryByKeys(m.Query,
		"id", "name", "code", "create_time:date",
	))
}

func (m *PerElement) inRouters() error {
	var err error
	if strings.TrimSpace(m.RouterIds) != "" && m.Query[vos.DeepKey] == true {
		item := &PerRouter{
			Query: vos.Query{
				vos.DeepKey: true,
				"id":        vos.SplitIds(m.RouterIds),
			},
		}

		m.Routers, err = item.GetList("remark", "create_time", "update_time")
	}

	return err
}
