package db

import (
	"fmt"
	"strings"
	"sync"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var schemaCache = &sync.Map{}

// var _locker = &sync.RWMutex{}

func HardTranslateStuctNameToDbName(db *gorm.DB, obj interface{}, name string) string {
	ret, err := TranslateStuctNameToDbName(db, obj, name)
	if err != nil {
		panic(err)
	}
	return ret
}

func TranslateStuctNameToDbName(db *gorm.DB, obj interface{}, name string) (string, error) {
	s, err := schema.Parse(obj, schemaCache, db.NamingStrategy)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(s.Fields); i++ {
		if s.Fields[i].Name == name {
			return s.Fields[i].DBName, nil
		}
	}
	return "", fmt.Errorf("not found:%s", name)
}

func CreateTableAndTrigger(err error, tx *gorm.DB, obj interface{}, tbName string, delProtect bool, updateProtect bool, updateProtectAll bool, setups ...(func() error)) error {
	if err != nil {
		return err
	}
	sch, err := GetSchema(tx, obj)
	if err != nil {
		return err
	}
	// log.Println(sch.Name)
	fields := make([]string, 0)
	for i := 0; i < len(sch.Fields); i++ {
		if !sch.Fields[i].Updatable {
			fields = append(fields, fmt.Sprintf("NEW.`%s`=OLD.`%s`", sch.Fields[i].DBName, sch.Fields[i].DBName))
		}
	}

	if !tx.Migrator().HasTable(obj) {
		err = tx.Table(tbName).Migrator().CreateTable(obj)
		if err != nil {
			return err
		}
		if updateProtect {
			if updateProtectAll {
				err = tx.Exec(fmt.Sprintf(`CREATE TRIGGER TR_%s_update BEFORE UPDATE ON %s FOR EACH ROW SIGNAL SQLSTATE '42000' SET MESSAGE_TEXT = "update rejected",MYSQL_ERRNO=1143;`, tbName, tbName)).Error
				if err != nil {
					return err
				}
			} else {
				fields := make([]string, 0)
				for i := 0; i < len(sch.Fields); i++ {
					// if sch.Fields[i].IgnoreMigration
					if sch.Fields[i].Creatable || sch.Fields[i].Updatable || sch.Fields[i].Readable {
						if !sch.Fields[i].Updatable {
							fields = append(fields, fmt.Sprintf("NEW.`%s`=OLD.`%s`", sch.Fields[i].DBName, sch.Fields[i].DBName))
						}
					}
				}
				if len(fields) > 0 {
					setStr := strings.Join(fields, ",")
					err = tx.Exec(fmt.Sprintf(`CREATE TRIGGER TR_%s_update BEFORE UPDATE ON %s FOR EACH ROW SET `+setStr, tbName, tbName)).Error
					if err != nil {
						return err
					}
				}
			}
		}
		//删除保护
		if delProtect {
			err = tx.Exec(fmt.Sprintf(`CREATE TRIGGER TR_%s_delete BEFORE DELETE ON %s FOR EACH ROW SIGNAL SQLSTATE '42000' SET MESSAGE_TEXT = "delete rejected",MYSQL_ERRNO=1143;`, tbName, tbName)).Error
			if err != nil {
				return err
			}
		}
		if setups != nil {
			for i := 0; i < len(setups); i++ {
				if err := setups[i](); err != nil {
					return err
				}
			}
		}
	}
	// else {
	// 	return gorm.ErrRegistered
	// }
	return nil
}

func GetSchema(db *gorm.DB, obj interface{}) (*schema.Schema, error) {
	return schema.Parse(obj, schemaCache, db.NamingStrategy)
}
