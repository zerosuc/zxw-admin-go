package global

import (
	"gorm.io/gorm"
	"time"
)

// BaseModel gorm 通用基础模型
type BaseModel struct {
	ID        uint           `json:"id" gorm:"primarykey"` // 主键ID
	CreatedAt time.Time      `json:"createdAt"`            // 创建时间
	UpdatedAt time.Time      `json:"updatedAt"`            // 更新时间
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`       // 删除时间
}
