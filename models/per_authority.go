package models

import (
	"strings"

	"github.com/go-xorm/xorm"
	"github.com/rimxgo/models/vos"
)

type PerAuthority struct {
	Id         int           `json:"id" xorm:"not null pk autoincr INT(11)"`
	ParentId   int           `json:"parentId" xorm:"INT(11)"`
	ElementIds string        `json:"elementIds,omitempty" xorm:"VARCHAR(255)"`
	Name       string        `json:"name" xorm:"VARCHAR(255)"`
	Remark     string        `json:"remark,omitempty" xorm:"VARCHAR(255)"`
	CreateTime vos.JsonTime  `json:"createTime" xorm:"DATETIME created"`
	UpdateTime vos.JsonTime  `json:"updateTime" xorm:"DATETIME updated"`
	DeleteTime vos.JsonTime  `json:"-" xorm:"DATETIME deleted"`
	Query      vos.Query     `json:"-" xorm:"-"` // 查询条件
	Elements   []*PerElement `json:"elements,omitempty" xorm:"-"`
}

func (m *PerAuthority) GetTree() ([]*vos.NodeData, error) {
	ls := vos.TreeData{}
	sess := m.listSess("element_ids", "remark", "create_time", "update_time")

	err := sess.Iterate(new(PerAuthority), func(i int, bean interface{}) error {
		item := bean.(*PerAuthority)
		ls = append(ls, &vos.NodeData{
			Id:       item.Id,
			ParentId: item.ParentId,
			Label:    item.Name,
			Children: make([]*vos.NodeData, 0),
		})
		return nil
	})

	if err != nil {
		return nil, err
	}

	return ls.GenerateTree(0), nil
}

func (m *PerAuthority) GetList(oCols ...string) ([]*PerAuthority, error) {
	ls := make([]*PerAuthority, 0)
	sess := m.listSess(oCols...)

	err := sess.Iterate(new(PerAuthority), func(i int, bean interface{}) error {
		item := bean.(*PerAuthority)
		item.Query = vos.Query{vos.DeepKey: m.Query[vos.DeepKey]}
		item.inElements()
		ls = append(ls, item)
		return nil
	})

	return ls, err
}

func (m *PerAuthority) GetPage(p *vos.Page) ([]*PerAuthority, error) {
	ls := make([]*PerAuthority, 0)
	sess := m.listSess()

	var err error
	p.Total, err = sess.Limit(p.PageSize, p.Start).FindAndCount(&ls)
	return ls, err
}

func (m *PerAuthority) GetOne() (bool, error) {
	has, err := engine.Get(m)
	if has && err == nil {
		m.Query = vos.Query{vos.DeepKey: true}
		err = m.inElements()
	}

	return has, err
}

func (m *PerAuthority) Insert() (int64, error) {
	return engine.AllCols().InsertOne(m)
}

func (m *PerAuthority) Update() (int64, error) {
	return engine.ID(m.Id).MustCols("parent_id", "remark").Update(m)
}

func (m *PerAuthority) Delete() (int64, error) {
	return engine.ID(m.Id).Delete(new(PerAuthority))
}

func (m *PerAuthority) IsExist() (bool, error) {
	return engine.Exist(m)
}

func (m *PerAuthority) Validate() error {
	return nil
}

func (m *PerAuthority) listSess(oCols ...string) *xorm.Session {
	return engine.Table(m).Omit(oCols...).Where(vos.SetQueryByKeys(m.Query,
		"id", "name", "create_time:date",
	))
}

func (m *PerAuthority) inElements() error {
	var err error
	m.Elements = make([]*PerElement, 0)
	if strings.TrimSpace(m.ElementIds) != "" && m.Query[vos.DeepKey] == true {
		item := &PerElement{
			Query: vos.Query{
				vos.DeepKey: true,
				"id":        vos.SplitIds(m.ElementIds),
			},
		}

		m.Elements, err = item.GetList("remark", "create_time", "update_time")
	}

	return err
}
