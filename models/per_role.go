package models

import (
	"github.com/go-xorm/xorm"
	"github.com/rimxgo/models/vos"
)

type PerRole struct {
	Id           int          `json:"id" xorm:"not null pk autoincr INT(11)"`
	AuthorityIds string       `json:"authorityIds" xorm:"VARCHAR(255)"`
	Name         string       `json:"name" xorm:"VARCHAR(255)"`
	Remark       string       `json:"remark" xorm:"VARCHAR(255)"`
	CreateTime   vos.JsonTime `json:"createTime" xorm:"DATETIME created"`
	UpdateTime   vos.JsonTime `json:"updateTime" xorm:"DATETIME updated"`
	DeleteTime   vos.JsonTime `json:"-" xorm:"DATETIME deleted"`
	Query        vos.Query    `json:"-" xorm:"-"` // 查询条件
}

func (m *PerRole) GetList(oCols ...string) ([]*PerRole, error) {
	ls := make([]*PerRole, 0)
	sess := m.listSess(oCols...)

	err := sess.Find(&ls)
	return ls, err
}

func (m *PerRole) GetPage(p *vos.Page) ([]*PerRole, error) {
	ls := make([]*PerRole, 0)
	sess := m.listSess()

	var err error
	p.Total, err = sess.Limit(p.PageSize, p.Start).FindAndCount(&ls)
	return ls, err
}

func (m *PerRole) GetOne() (bool, error) {
	return engine.Get(m)
}

func (m *PerRole) Insert() (int64, error) {
	return engine.AllCols().InsertOne(m)
}

func (m *PerRole) Update() (int64, error) {
	return engine.ID(m.Id).MustCols("remark").Update(m)
}

func (m *PerRole) Delete() (int64, error) {
	return engine.ID(m.Id).Delete(new(PerRole))
}

func (m *PerRole) IsExist() (bool, error) {
	return engine.Exist(m)
}

func (m *PerRole) Validate() error {
	return nil
}

func (m *PerRole) listSess(oCols ...string) *xorm.Session {
	return engine.Table(m).Omit(oCols...).Where(vos.SetQueryByKeys(m.Query,
		"id", "name", "create_time:date",
	))
}
