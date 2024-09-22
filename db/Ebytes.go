package db

import (
	"database/sql/driver"
	"fmt"
)

type EnBytes []byte

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *EnBytes) Scan(value interface{}) error {
	// log.Println("Scan")
	encrypted := []byte{}
	if _bytes, ok := value.([]byte); ok {
		encrypted = _bytes
	} else {
		return fmt.Errorf("err data")
	}

	if len(encrypted) == 0 {
		*j = EnBytes{}
		return nil
	}
	if k, err := aes.Decode(encrypted); err != nil {
		return err
	} else {
		*j = EnBytes(k)
	}
	return nil
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (j EnBytes) Value() (driver.Value, error) {
	if len(j) == 0 {
		return "", nil
	}
	return aes.Encode([]byte(j))
}

func (j EnBytes) GormDataType() string {
	return "bytes"
}

func (j EnBytes) Data() []byte {
	return []byte(j)
}
