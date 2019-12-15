package vos

import (
	"fmt"
	"strings"
	"time"

	"github.com/rimxgo/helper/logs"
	"xorm.io/builder"
)

type Page struct {
	PageNum  int   `json:"pageNum" form:"pageNum"`
	PageSize int   `json:"pageSize" form:"pageSize"`
	Total    int64 `json:"total" form:"-"`
	Start    int   `json:"-" form:"-"`
}

type JsonTime time.Time

func (j JsonTime) MarshalJSON() ([]byte, error) {
	t := time.Time(j)
	if t.IsZero() {
		return []byte(`""`), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, t.Format(time.RFC3339))), nil
}

func (j *JsonTime) UnmarshalJSON(data []byte) error {
	var err error
	str, ln, t := string(data), len(data), time.Time{}
	if str == "null" || str == `""` {
		return nil
	}

	if ln == 21 {
		t, err = time.ParseInLocation(`"2006-01-02 15:04:05"`, str, time.Local)
	} else if ln == 12 {
		t, err = time.ParseInLocation(`"2006-01-02"`, str, time.Local)
	} else {
		t, err = time.Parse(fmt.Sprintf(`"%s"`, time.RFC3339), str)
	}

	*j = JsonTime(t)
	return err
}

// 根据keys设置查询条件，
// keys规则："[表别名.]{表字段}[:条件类型][:条件方式]".
// 条件类型：|eq|neq|gt|gte|lt|lte|like|%like|like%|in|notin|range|date|datetime|
// 条件方式：|and|or|
func SetQueryByKeys(query map[string]interface{}, keys ...string) builder.Cond {
	cond := builder.NewCond()

	for _, key := range keys {
		k := strings.Split(key, ":")
		kk := k[0]
		if i := strings.Index(kk, "."); i > -1 {
			kk = kk[i+1:]
		}

		v, ok := query[snakeToCamel(kk, false)]
		if ok {
			goto end
		}

		if v, ok = query[snakeToCamel(kk, true)]; ok {
			goto end
		}

		if v, ok = query[kk]; ok {
			goto end
		}

	end:
		if v1, ok1 := v.(string); !ok || (ok1 && strings.TrimSpace(v1) == "") {
			continue
		} else if ok1 {
			v = strings.TrimSpace(v1)
		}

		if len(k) > 2 && k[2] == "or" {
			cond = cond.Or(whereCond(k, v))
		} else {
			cond = cond.And(whereCond(k, v))
		}
	}

	return cond
}

func snakeToCamel(snake string, big bool) string {
	if snake == "" {
		return ""
	}

	var key string
	for i, v := range strings.Split(snake, "_") {
		vl := []rune(v)
		if len(vl) > 0 {
			if i == 0 && !big {
				if bool(vl[0] >= 'A' && vl[0] <= 'Z') {
					vl[0] += 32 //首字母小写
				}
			} else {
				if bool(vl[0] >= 'a' && vl[0] <= 'z') {
					vl[0] -= 32 //首字母大写
				}
			}
			key += string(vl)
		}
	}

	return key
}

func whereCond(key []string, val interface{}) builder.Cond {
	if len(key) > 1 {
		switch key[1] {
		case "eq", "":
			return builder.Eq{key[0]: val}
		case "neq":
			return builder.Neq{key[0]: val}
		case "gt":
			return builder.Gt{key[0]: val}
		case "gte":
			return builder.Gte{key[0]: val}
		case "lt":
			return builder.Lt{key[0]: val}
		case "lte":
			return builder.Lte{key[0]: val}
		case "like", "%like", "like%":
			if v, ok := val.(string); ok {
				return builder.Like{key[0], strings.ReplaceAll(key[1], "like", v)}
			}
			logs.Logger.Debugf("Parameter `%s` type error, expecting a string.", key[0])
			return builder.NewCond()
		case "in":
			if v, ok := val.([]interface{}); ok {
				return builder.In(key[0], v...)
			}
			logs.Logger.Debugf("Parameter `%s` type error, expecting an array.", key[0])
			return builder.NewCond()
		case "notin":
			if v, ok := val.([]interface{}); ok {
				return builder.NotIn(key[0], v...)
			}
			logs.Logger.Debugf("Parameter `%s` type error, expecting an array.", key[0])
			return builder.NewCond()
		case "range":
			if v, ok := val.([]interface{}); ok {
				return rangeCond(key[0], v)
			}
			logs.Logger.Debugf("Parameter `%s` type error, expecting an array.", key[0])
			return builder.NewCond()
		case "datetime":
			if v, ok := val.(string); ok {
				return builder.Eq{key[0]: v}
			}

			if v, ok := val.([]interface{}); ok {
				return rangeCond(key[0], v)
			}
			logs.Logger.Debugf("Parameter `%s` type error, expecting a string or an array.", key[0])
			return builder.NewCond()
		case "date":
			if v, ok := val.(string); ok {
				vv := []interface{}{v, v + " 23:59:59"}
				return rangeCond(key[0], vv)
			}

			if v, ok := val.([]interface{}); ok {
				if len(v) > 1 {
					if v1, _ := v[1].(string); strings.TrimSpace(v1) != "" {
						v[1] = strings.TrimSpace(v1) + " 23:59:59"
					}
				}
				return rangeCond(key[0], v)
			}
			logs.Logger.Debugf("Parameter `%s` type error, expecting a string or an array.", key[0])
			return builder.NewCond()
		default:
		}
	}
	return builder.Eq{key[0]: val}
}

func rangeCond(col string, val []interface{}) builder.Cond {
	cond := builder.NewCond()

	if len(val) > 0 {
		if v0, ok := val[0].(string); !ok || strings.TrimSpace(v0) != "" {
			if ok {
				val[0] = strings.TrimSpace(v0)
			}
			cond = cond.And(builder.Gte{col: val[0]})
		}
		if len(val) > 1 {
			if v1, ok := val[1].(string); !ok || strings.TrimSpace(v1) != "" {
				if ok {
					val[1] = strings.TrimSpace(v1)
				}
				cond = cond.And(builder.Lte{col: val[1]})
			}
		}
	}

	return cond
}
