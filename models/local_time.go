package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// LocalTime 自定义时间类型
type LocalTime time.Time

// MarshalJSON 实现json序列化接口
func (t LocalTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

// UnmarshalJSON 实现json反序列化接口
func (t *LocalTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	// 去掉引号
	str := string(data)[1 : len(data)-1]
	parsed, err := time.Parse("2006-01-02 15:04:05", str)
	if err != nil {
		return err
	}
	*t = LocalTime(parsed)
	return nil
}

// Value 实现 driver.Valuer 接口
func (t LocalTime) Value() (driver.Value, error) {
	return time.Time(t), nil
}

// Scan 实现 sql.Scanner 接口
func (t *LocalTime) Scan(v interface{}) error {
	if value, ok := v.(time.Time); ok {
		*t = LocalTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to LocalTime", v)
}

// String 实现 Stringer 接口
func (t LocalTime) String() string {
	return time.Time(t).Format("2006-01-02 15:04:05")
}
