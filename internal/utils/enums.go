package utils

import "database/sql/driver"

type CodeType string

const (
	CodeTypeForgotPassword CodeType = "forgot_password"
	CodeTypeVerifyEmail    CodeType = "verify_email"
)

func (p *CodeType) Scan(value interface{}) error {
	*p = CodeType(value.([]byte))
	return nil
}

func (p CodeType) Value() (driver.Value, error) {
	return string(p), nil
}
