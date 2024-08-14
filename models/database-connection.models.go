package models

import "gorm.io/gorm"

type Dialector struct {
	*Config
}

type Config struct {
	DriverName           string
	DSN                  string
	WithoutQuotingCheck  bool
	PreferSimpleProtocol bool
	WithoutReturning     bool
	Conn                 gorm.ConnPool
}
