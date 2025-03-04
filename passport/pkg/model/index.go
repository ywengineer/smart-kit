package model

import (
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        uint           `gorm:"primarykey" redis:"id" json:"id"`
	CreatedAt time.Time      `json:"created_at" redis:"created_at" json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" redis:"updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at" redis:"deleted_at" json:"deleted_at"`
}
