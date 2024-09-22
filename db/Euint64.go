package db

import (
	"database/sql/driver"
	"encoding/binary"
	"fmt"
)

type EnUint64 uint64

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *EnUint64) Scan(value interface{}) error {
	// log.Println("Scan")
	encrypted := make([]byte, 0)
	if _bytes, ok := value.([]byte); ok {
		encrypted = _bytes
	} else if _bytes, ok := value.(string); ok {
		encrypted = []byte(_bytes)
	} else {
		return fmt.Errorf("err data")
	}

	if len(encrypted) == 0 {
		*j = EnUint64(0)
		return nil
	}
	// binary.BigEndian.Uint64(encrypted)
	if k, err := aes.Decode(encrypted); err != nil {
		return err
	} else if len(k) != 8 {
		return fmt.Errorf("err data")
	} else {
		*j = EnUint64(binary.BigEndian.Uint64(k))
	}
	return nil
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (j EnUint64) Value() (driver.Value, error) {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(j))
	return aes.Encode(buf)
}

func (j EnUint64) GormDataType() string {
	return "bytes"
}
