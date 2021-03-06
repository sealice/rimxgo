package models

import (
	"errors"
	"strings"

	"github.com/go-xorm/xorm"
	"github.com/rimxgo/models/vos"
)

type User struct {
	Id         int          `json:"id" xorm:"not null pk autoincr INT(11)"`
	RoleIds    string       `json:"roleIds" xorm:"VARCHAR(255)"`
	Name       string       `json:"name" xorm:"VARCHAR(45)"`
	Password   string       `json:"password,omitempty" xorm:"-> VARCHAR(45)"`
	Gender     int          `json:"gender" xorm:"TINYINT(1)"`
	Age        int          `json:"age" xorm:"TINYINT(3)"`
	CreateTime vos.JsonTime `json:"createTime" xorm:"created DATETIME"`
	UpdateTime vos.JsonTime `json:"updateTime" xorm:"updated DATETIME"`
	Query      vos.Query    `json:"-" xorm:"-"` // 查询条件
	Roles      []*PerRole   `json:"roles" xorm:"-"`
	Elements   []string     `json:"elements,omitempty" xorm:"-"`
	Routers    []string     `json:"-" xorm:"-"`
}

func (m *User) GetList(oCols ...string) ([]*User, error) {
	ls := make([]*User, 0)
	sess := m.listSess(oCols...)

	err := sess.Find(&ls)
	return ls, err
}

func (m *User) GetPage(p *vos.Page) ([]*User, error) {
	ls := make([]*User, 0)
	sess := m.listSess()
	sess2 := sess.Clone().Limit(p.PageSize, p.Start)

	err := sess2.Iterate(new(User), func(i int, bean interface{}) error {
		item := bean.(*User)
		item.inRoles()
		ls = append(ls, item)
		return nil
	})

	if err != nil {
		return ls, err
	}

	p.Total, err = sess.Count(new(User))

	return ls, err
}

func (m *User) GetOne() (bool, error) {
	has, err := engine.Get(m)
	if err != nil || !has {
		return has, err
	}

	err = m.inRoles()
	return has, err
}

func (m *User) Insert() (int64, error) {
	return engine.AllCols().InsertOne(m)
}

func (m *User) IsExist() (bool, error) {
	return engine.Exist(m)
}

func (m *User) Validate() error {
	if m.Name = strings.TrimSpace(m.Name); m.Name == "" {
		return errors.New("用户名不能为空")
	}

	if m.Password = strings.TrimSpace(m.Password); m.Password == "" {
		return errors.New("密码不能为空")
	}

	return nil
}

func (m *User) listSess(oCols ...string) *xorm.Session {
	return engine.Table(m).Omit(oCols...).Where(vos.SetQueryByKeys(m.Query,
		"id", "age", "gender", "create_time:date",
	))
}

func (m *User) inRoles() error {
	var err error
	m.Roles = make([]*PerRole, 0)
	if strings.TrimSpace(m.RoleIds) != "" {
		item := &PerRole{
			Query: vos.Query{
				"id": vos.SplitIds(m.RoleIds),
			},
		}

		m.Roles, err = item.GetList("remark", "create_time", "update_time")
	}

	return err
}

func (m *User) GetPermission() error {
	if len(m.Roles) > 0 {
		return m.permissionResolve(m.Roles)
	}

	return m.getPermission(&PerRole{}, m.RoleIds)
}

func (m *User) getPermission(st interface{}, ids string) error {
	if strings.TrimSpace(ids) == "" {
		return nil
	}

	var err error
	var ls interface{}
	var oCols = []string{`name`, `remark`, `create_time`, `update_time`}

	switch per := st.(type) {
	case *PerRole:
		per.Query = vos.Query{"id": vos.SplitIds(ids)}
		ls, err = per.GetList(oCols...)
	case *PerAuthority:
		per.Query = vos.Query{"id": vos.SplitIds(ids)}
		ls, err = per.GetList(oCols...)
	case *PerElement:
		per.Query = vos.Query{"id": vos.SplitIds(ids)}
		ls, err = per.GetList(oCols...)
	case *PerRouter:
		per.Query = vos.Query{"id": vos.SplitIds(ids)}
		ls, err = per.GetList(oCols...)
	}

	if err != nil || ls == nil {
		return err
	}

	return m.permissionResolve(ls)
}

func (m *User) permissionResolve(ls interface{}) error {
	ids := ""

	switch lst := ls.(type) {
	case []*PerRole:
		for _, s := range lst {
			ids += s.AuthorityIds + ","
		}
		return m.getPermission(&PerAuthority{}, ids)

	case []*PerAuthority:
		for _, s := range lst {
			ids += s.ElementIds + ","
		}
		return m.getPermission(&PerElement{}, ids)

	case []*PerElement:
		for _, s := range lst {
			ids += s.RouterIds + ","
			m.Elements = append(m.Elements, s.Code)
		}
		return m.getPermission(&PerRouter{}, ids)

	case []*PerRouter:
		for _, s := range lst {
			m.Routers = append(m.Routers, s.Path)
		}
	}

	return nil
}
