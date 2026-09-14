package models

import (
	"time"
)

type Quality struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	QualityName string    `gorm:"not null" json:"quality_name"`
	Code        string    `gorm:"not null" json:"code"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
