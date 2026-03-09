package model

import (
	"time"

	"gorm.io/gorm"
)

// UserSession 用户会话模型
type UserSession struct {
	ID        int            `gorm:"primaryKey;autoIncrement"` // 主键ID
	CreatedAt time.Time      `gorm:"column:created_at"`        // 创建时间
	UpdatedAt time.Time      `gorm:"column:updated_at"`        // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`        // 删除时间

	UserUID   int       `gorm:"column:user_uid"`   // 用户ID
	SessionID string    `gorm:"column:session_id"` // 会话ID
	IP        string    `gorm:"column:ip"`         // IP
	UserAgent string    `gorm:"column:user_agent"` // 用户代理

	ExpiresAt time.Time `gorm:"column:expires_at"` // 过期时间
}

func (UserSession) TableName() string {
	return "user_session"
}
