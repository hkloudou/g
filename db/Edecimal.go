package db

import (
	"database/sql/driver"
	"fmt"

	"github.com/shopspring/decimal"
)

type EnDecimal struct {
	Data decimal.Decimal
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *EnDecimal) Scan(value interface{}) error {
	if _bytes, ok := value.([]byte); !ok {
		return fmt.Errorf("err fmt EnDecimal:%v", value)
	} else if k, err := aes.Decode(_bytes); err != nil {
		return err
	} else if len(k) == 0 {
		*j = EnDecimal{
			Data: decimal.Zero,
		}
		return nil
	} else if d, err := decimal.NewFromString(string(k)); err != nil {
		return err
	} else {
		*j = EnDecimal{
			Data: d,
		}
	}
	return nil
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (j EnDecimal) Value() (driver.Value, error) {
	data := j.Data.String()
	if len(data) == 0 {
		return []byte{}, nil
	}
	if k, err := aes.Encode([]byte(data)); err != nil {
		return []byte{}, err
	} else {
		return k, nil
	}
}

func (d EnDecimal) MarshalJSON() ([]byte, error) {
	return d.Data.MarshalJSON()
}

func (d *EnDecimal) UnmarshalJSON(decimalBytes []byte) error {
	return d.Data.UnmarshalJSON(decimalBytes)
}

func (j EnDecimal) Add(dest decimal.Decimal) EnDecimal {
	return EnDecimal{
		Data: j.Data.Add(dest),
	}
}

func (j EnDecimal) Equal(dest decimal.Decimal) bool {
	return j.Data.Equal(dest)
}

func (j EnDecimal) LessThan(dest decimal.Decimal) bool {
	return j.Data.LessThan(dest)
}

func (j EnDecimal) LessThanOrEqual(dest decimal.Decimal) bool {
	return j.Data.LessThanOrEqual(dest)
}

func (j EnDecimal) GreaterThan(dest decimal.Decimal) bool {
	return j.Data.GreaterThan(dest)
}

func (j EnDecimal) GreaterThanOrEqual(dest decimal.Decimal) bool {
	return j.Data.GreaterThanOrEqual(dest)
}
