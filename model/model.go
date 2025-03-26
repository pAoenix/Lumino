package model

import (
	"gorm.io/gorm"
	"time"
)

const PgDBName = "postgres"

// Model -
type Model struct {
	ID        uint           `gorm:"primarykey" json:"id" form:"id"`
	CreatedAt time.Time      `json:"created_at" form:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" form:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"  form:"deleted_at" swaggertype:"primitive,string" example:"2025-03-26T00:00:00Z"`
}
