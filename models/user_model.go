package models

import (
	"errors"
	"strings"

	"github.com/rimxgo/models/vos"
)

func init() {
	engine.Sync2(new(User))
}

type User struct {
	Id         int                    `json:"id" xorm:"not null pk autoincr INT(11)"`
	Name       string                 `json:"name" xorm:"VARCHAR(45)"`
	Password   string                 `json:"password,omitempty" xorm:"-> VARCHAR(45)"`
	Gender     int                    `json:"gender" xorm:"TINYINT(1)"`
	Age        int                    `json:"age" xorm:"TINYINT(3)"`
	CreateTime vos.JsonTime           `json:"createTime" xorm:"created DATETIME"`
	UpdateTime vos.JsonTime           `json:"updateTime" xorm:"updated DATETIME"`
	Query      map[string]interface{} `json:"-" xorm:"-"` // 查询条件
}

func (m *User) GetList(p *vos.Page) ([]*User, error) {
	ls := make([]*User, 0)
	sess := engine.Table(m).Where(vos.SetQueryByKeys(m.Query,
		"id", "age", "gender", "create_time:date",
	))

	var err error
	p.Total, err = sess.Limit(p.PageSize, p.Start).FindAndCount(&ls)
	return ls, err
}

func (m *User) GetOne() (bool, error) {
	return engine.Get(m)
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
