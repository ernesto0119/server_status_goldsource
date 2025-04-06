package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Base struct {
	Uuid      uuid.UUID      `gorm:"type:uuid;primary_key;" sql:"index"`
	CreatedAt time.Time      `json:"created_at" sql:"index"`
	UpdatedAt time.Time      `json:"update_at" sql:"index"`
	DeletedAt gorm.DeletedAt `sql:"index" json:"deleted_at"`
}

func (base *Base) BeforeCreate(tx *gorm.DB) (err error) {
	// UUID version 4
	base.Uuid = uuid.New()
	return
}
