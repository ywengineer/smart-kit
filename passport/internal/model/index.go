package model

import (
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        uint           `gorm:"primarykey" redis:"id" json:"id"`
	CreatedAt time.Time      `redis:"created_at" json:"created_at"`
	UpdatedAt time.Time      `redis:"updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at" redis:"deleted_at" json:"deleted_at"`
}
