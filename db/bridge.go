package db

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/hkloudou/xlib/xruntime"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// var _debug = false
type baseBridge struct {
	_tbName string
}
type Bridge[t any] struct {
	*gorm.DB
	baseBridge
	// beforeQuery func()
}

func D[T any](_db *gorm.DB) *Bridge[T] {
	_debug := xruntime.IsDevelopmentMode()

	if _db == nil {
		if _debug {
			return &Bridge[T]{
				DB: Get().Debug(),
			}
		} else {
			return &Bridge[T]{
				DB: Get(),
			}
		}
	}
	if _debug {
		return &Bridge[T]{
			DB: _db.Debug(),
		}
	} else {
		return &Bridge[T]{
			DB: _db,
		}
	}
}

func DWithTableName[T any](_db *gorm.DB, name string) *Bridge[T] {
	tmp := D[T](_db)
	tmp._tbName = name
	return tmp
}

func (m *Bridge[T]) cloneWithDB(_db *gorm.DB) *Bridge[T] {
	return &Bridge[T]{
		baseBridge: m.baseBridge,
		DB:         _db,
	}
}

func (m *Bridge[T]) F(name string) string {
	var obj T
	return HardTranslateStuctNameToDbName(m.DB, &obj, name)
}

func (m *Bridge[T]) EnsureTable(del, up, upAll bool) error {
	var obj T
	// return m.Wise().Transaction(func(tx *Bridge[T]) error {
	// 	return CreateTableAndTrigger(nil, tx.DB, &obj, m.TableName(), del, up, upAll)
	// })
	return CreateTableAndTrigger(nil, m.DB.Debug(), &obj, m.TableName(), del, up, upAll)
}

//	func (m *Bridge[T]) DebugGet__tbName() string {
//		return m._tbName
//	}
func (m *Bridge[T]) TableName() string {
	if m._tbName != "" {
		return m._tbName
	}
	var obj T
	sch, _ := GetSchema(m.DB, &obj)
	return sch.Table
}

func (m *Bridge[T]) Table(tbNames ...string) *Bridge[T] {
	if len(tbNames) > 0 {
		return m.cloneWithDB(m.DB.Table(tbNames[0]))
	}
	return m.cloneWithDB(m.DB.Table(m.TableName()))
}

func (m *Bridge[T]) Wise() *Bridge[T] {
	if m.Statement.Table == "" {
		return m.Table().Wise()
	}
	var obj T
	if m.Statement.Model == nil {
		return m.Model(&obj).Wise()
	}
	return m
}

func (m *Bridge[T]) replaceStr(query interface{}, _from, _to string) interface{} {
	if reflect.TypeOf(query).String() == "string" {
		str := query.(string)
		findFrom := func(source, from, to string) string {
			i1 := strings.Index(source, from)
			if i1 == -1 {
				return ""
			}
			i2 := strings.Index(source[i1:], to)
			if i2 == -1 {
				return ""
			}
			return source[i1+1 : i1+i2]
		}
		for {
			sub := findFrom(str, _from, _to)
			if sub == "" {
				return str
			}
			str = strings.Replace(str, fmt.Sprintf("%s%s%s", _from, sub, _to), fmt.Sprintf("`%s`", m.F(sub)), 1)
		}
	}
	return query
}
func (m *Bridge[T]) Model(val interface{}) *Bridge[T] {
	return m.cloneWithDB(m.DB.Model(val))
}

// Mix find
func (m *Bridge[T]) Where(query interface{}, args ...interface{}) *Bridge[T] {
	// if m.Statement.Table == "" {
	// 	return m.Table().Where(query, args...)
	// }
	return m.cloneWithDB(m.DB.Where(m.replaceStr(m.replaceStr(query, "{", "}"), "[", "]"), args...))
}

func (m *Bridge[T]) Update(column string, value interface{}) *Bridge[T] {
	query := m.replaceStr(m.replaceStr(column, "{", "}"), "[", "]")
	return m.cloneWithDB(m.DB.Update(query.(string), value))
}
func (m *Bridge[T]) Order(query interface{}) *Bridge[T] {
	return m.cloneWithDB(m.DB.Order(m.replaceStr(m.replaceStr(query, "{", "}"), "[", "]")))
}

func (m *Bridge[T]) Limit(limit int) *Bridge[T] {
	return m.cloneWithDB(m.DB.Limit(limit))
}

func (m *Bridge[T]) Transaction(fc func(tx *Bridge[T]) error, opts ...*sql.TxOptions) error {
	return m.DB.Transaction(func(tx *gorm.DB) error {
		return fc(m.cloneWithDB(tx))
	}, opts...)
}

func (m *Bridge[T]) Clauses(conds ...clause.Expression) *Bridge[T] {
	return m.cloneWithDB(m.DB.Clauses(conds...))
}

func (m *Bridge[T]) Select(query interface{}, args ...interface{}) *Bridge[T] {
	return m.cloneWithDB(m.DB.Select(m.replaceStr(m.replaceStr(query, "{", "}"), "[", "]"), args...))
}

func (m *Bridge[T]) Take(conds ...interface{}) (*T, error) {
	var obj T
	res := m.DB.Take(&obj, conds...)
	if res.Error != nil {
		return nil, res.Error
	}
	return &obj, nil
}

func (m *Bridge[T]) First(conds ...interface{}) (*T, error) {
	var obj T
	res := m.DB.First(&obj, conds...)
	if res.Error != nil {
		return nil, res.Error
	}
	return &obj, nil
}

func (m *Bridge[T]) Last(conds ...interface{}) (*T, error) {
	var obj T
	res := m.DB.Last(&obj, conds...)
	if res.Error != nil {
		return nil, res.Error
	}
	return &obj, nil
}
