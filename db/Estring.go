package db

import (
	"database/sql/driver"
	"fmt"

	"github.com/hkloudou/xlib/xencrypt"
)

var aes = xencrypt.NewAesEnDecrypter(
	[]byte{131, 23, 59, 213, 2, 93, 217, 90, 16, 73, 104, 98, 21, 138, 23, 12},
	[]byte{131, 23, 59, 213, 2, 93, 217, 90, 16, 73, 104, 98, 21, 138, 23, 12},
)

type EnString string

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *EnString) Scan(value interface{}) error {
	// log.Println("Scan")
	encrypted := []byte{}
	if _bytes, ok := value.([]byte); ok {
		encrypted = _bytes
	} else if _bytes, ok := value.(string); ok {
		encrypted = []byte(_bytes)
	} else {
		return fmt.Errorf("err data")
	}

	if len(encrypted) == 0 {
		*j = EnString("")
		return nil
	}
	if k, err := aes.Decode(encrypted); err != nil {
		return err
	} else {
		*j = EnString(k)
	}
	return nil
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (j EnString) Value() (driver.Value, error) {
	if len(j) == 0 {
		return "", nil
	}
	return aes.Encode([]byte(j))
}

func (j EnString) GormDataType() string {
	return "bytes"
}

func (j EnString) String() string {
	return string(j)
}
