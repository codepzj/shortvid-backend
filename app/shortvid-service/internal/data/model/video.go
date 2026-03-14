package model

import (
	"time"

	"gorm.io/gorm"
)

type Video struct {
	ID        int            `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`

	VGroup string `gorm:"column:vgroup"` // 视频底层分组
	UID    int    `gorm:"column:uid"`    // 用户UID

	Description string `gorm:"column:description"` // 视频描述
	Category    string `gorm:"column:category"`    // 视频分类
	Tags        string `gorm:"column:tags"`        // 视频标签
	CustomTags  string `gorm:"column:custom_tags"` // 自定义标签
	LikeCount   int    `gorm:"column:like_count"`  // 点赞数
	ViewCount   int    `gorm:"column:view_count"`  // 浏览数
	Status      int    `gorm:"column:status"`      // 1-上传中, 2-处理中, 3-审核中, 4-发布成功, 5-发布失败, 6-封禁
}

func (Video) TableName() string {
	return "video"
}
