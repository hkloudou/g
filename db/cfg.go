package db

import (
	"fmt"
	"reflect"

	"github.com/hkloudou/xlib/xcolor"
	"github.com/hkloudou/xlib/xface"
	"github.com/hkloudou/xlib/xflag"
	"github.com/hkloudou/xlib/xruntime"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type cfg struct {
	_instance *gorm.DB
}

func Use() {
	var tmp xface.FlagConfig[gorm.DB] = &cfg{}
	xface.FlagRegister("base/db", tmp)
}

func Get() *gorm.DB {
	obj, err := xface.FlagGet[gorm.DB]("base/db")
	if err != nil {
		panic(err)
	}
	return obj
}

func WGet() *gorm.DB {
	obj, err := xface.FlagGet[gorm.DB]("base/db")
	if err != nil {
		panic(err)
	}
	_debug := xruntime.IsDevelopmentMode()
	if _debug {
		return obj.Debug()
	}
	return obj
}

func (m *cfg) Flags() []xflag.Flag {
	return []xflag.Flag{
		&xflag.StringFlag{
			Name:  "mysql_host",
			Value: "127.0.0.1",
		},
		&xflag.IntFlag{
			Name:  "mysql_port",
			Value: 3306,
		},
		&xflag.StringFlag{
			Name:  "mysql_user",
			Value: "root",
		},
		&xflag.StringFlag{
			Name:  "mysql_pwd",
			Value: "root",
		},
		&xflag.StringFlag{
			Name:     "mysql_db",
			Value:    "",
			Required: true,
		},
		&xflag.StringFlag{
			Name:   "mysql_parameter",
			Value:  "charset=utf8mb4&parseTime=True&loc=Local",
			Hidden: true,
		},
	}
}

func (m *cfg) Action(c *xflag.Context) error {
	tag := fmt.Sprintf(xcolor.Red("%-15s"), reflect.TypeOf(*m).PkgPath())
	addr := fmt.Sprintf("%s:%d", c.String("mysql_host"), c.Int("mysql_port"))
	fmt.Println(tag, xcolor.Yellow(fmt.Sprintf("[%s] connecting", addr)))
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		c.String("mysql_user"),
		c.String("mysql_pwd"),
		c.String("mysql_host"),
		c.Int("mysql_port"),
		c.String("mysql_db"),
	)
	if c.String("mysql_parameter") != "" {
		dns += ("?" + c.String("mysql_parameter"))
	}

	if db, err := gorm.Open(mysql.Open(dns), &gorm.Config{}); err != nil {
		return err
	} else {
		m._instance = db
	}
	fmt.Println(tag, xcolor.Green(fmt.Sprintf("[%s] connected", addr)))
	schema.RegisterSerializer("json", JSONSerializer{})
	return nil
}

func (m *cfg) Instance() *gorm.DB {
	return m._instance
}
