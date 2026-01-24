package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        int            `gorm:"primaryKey;autoIncrement"` // 主键ID
	CreatedAt time.Time      `gorm:"column:created_at"`        // 创建时间
	UpdatedAt time.Time      `gorm:"column:updated_at"`        // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`        // 删除时间

	UserUID     int       `gorm:"unique;column:user_uid"` // 用户唯一ID
	Nickname    string    `gorm:"column:nickname"`        // 昵称
	Avatar      string    `gorm:"column:avatar"`          // 头像
	Email       string    `gorm:"column:email"`           // 邮箱
	Provider    string    `gorm:"column:provider"`        // 主要登录提供商: firebase
	ProviderUID string    `gorm:"column:provider_uid"`    // 第三方平台用户UID
	LastLoginAt time.Time `gorm:"column:last_login_at"`   // 最后登录时间
	LoginCount  int       `gorm:"column:login_count"`     // 登录次数
	Status      int       `gorm:"column:status"`          // 状态: 1: 正常, 2: 禁用
}

func (User) TableName() string {
	return "user"
}
