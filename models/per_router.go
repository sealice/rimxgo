package models

import (
	"github.com/go-xorm/xorm"
	"github.com/rimxgo/models/vos"
)

type PerRouter struct {
	Id         int          `json:"id" xorm:"not null pk autoincr INT(11)"`
	Name       string       `json:"name" xorm:"VARCHAR(255)"`
	Path       string       `json:"path" xorm:"VARCHAR(255)"`
	Remark     string       `json:"remark,omitempty" xorm:"VARCHAR(255)"`
	CreateTime vos.JsonTime `json:"createTime" xorm:"DATETIME created"`
	UpdateTime vos.JsonTime `json:"updateTime" xorm:"DATETIME updated"`
	DeleteTime vos.JsonTime `json:"-" xorm:"DATETIME deleted"`
	Query      vos.Query    `json:"-" xorm:"-"` // 查询条件
}

func (m *PerRouter) GetList(oCols ...string) ([]*PerRouter, error) {
	ls := make([]*PerRouter, 0)
	sess := m.listSess(oCols...)

	err := sess.Find(&ls)
	return ls, err
}

func (m *PerRouter) GetPage(p *vos.Page) ([]*PerRouter, error) {
	ls := make([]*PerRouter, 0)
	sess := m.listSess()

	var err error
	p.Total, err = sess.Limit(p.PageSize, p.Start).FindAndCount(&ls)
	return ls, err
}

func (m *PerRouter) GetOne() (bool, error) {
	return engine.Get(m)
}

func (m *PerRouter) Insert() (int64, error) {
	return engine.AllCols().InsertOne(m)
}

func (m *PerRouter) Update() (int64, error) {
	return engine.ID(m.Id).MustCols("remark").Update(m)
}

func (m *PerRouter) Delete() (int64, error) {
	return engine.ID(m.Id).Delete(new(PerRouter))
}

func (m *PerRouter) IsExist() (bool, error) {
	return engine.Exist(m)
}

func (m *PerRouter) Validate() error {
	return nil
}

func (m *PerRouter) listSess(oCols ...string) *xorm.Session {
	return engine.Table(m).Omit(oCols...).Where(vos.SetQueryByKeys(m.Query,
		"id", "name", "path", "create_time:date",
	))
}
