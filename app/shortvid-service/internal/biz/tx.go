package biz

import "gorm.io/gorm"

// 事务聚合器
type TxRepo interface {
	ExecFunc(fn func(*gorm.DB) error) error
}
