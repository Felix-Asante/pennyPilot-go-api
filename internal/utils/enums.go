package utils

import (
	"database/sql/driver"
	"fmt"
)

type CodeType string

const (
	CodeTypeForgotPassword CodeType = "forgot_password"
	CodeTypeVerifyEmail    CodeType = "verify_email"
)

func (p *CodeType) Scan(value interface{}) error {
	if value == nil {
		*p = ""
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*p = CodeType(v)
	case string:
		*p = CodeType(v)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}
func (p CodeType) Value() (driver.Value, error) {
	return string(p), nil
}
