package rdbs

import (
	"gorm.io/gorm"
	"time"
)

type DeletedAt struct {
	gorm.DeletedAt
}

type Model struct {
	ID        uint      `gorm:"primarykey" redis:"id" json:"id"`
	CreatedAt time.Time `redis:"created_at" json:"created_at"`
	UpdatedAt time.Time `redis:"updated_at" json:"updated_at"`
	DeletedAt DeletedAt `gorm:"index" json:"deleted_at" redis:"deleted_at" json:"deleted_at"`
}

func (d *DeletedAt) Unix() int64 {
	if d.Valid {
		return d.Time.Unix()
	}
	return 0
}
