package model

import (
	"time"

	"gorm.io/gorm"
)

// 账户模型
type Account struct {
	ID        int            `gorm:"primaryKey;autoIncrement"` // 主键ID
	CreatedAt time.Time      `gorm:"column:created_at"`        // 创建时间
	UpdatedAt time.Time      `gorm:"column:updated_at"`        // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`        // 删除时间

	UserUID int `gorm:"column:user_uid"` // 用户唯一ID

	Email    string `gorm:"column:email"`    // 邮箱
	Password string `gorm:"column:password"` // 密码(bcrypt)

	Provider    string `gorm:"column:provider"`     // 主要登录提供商
	ProviderUID string `gorm:"column:provider_uid"` // 第三方平台用户UID
}

func (Account) TableName() string {
	return "account"
}
