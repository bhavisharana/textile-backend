package models

import (
	"time"
)

type Order struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	LabdipID    *uint     `json:"labdip_id,omitempty"`
	LabdipNo    string    `gorm:"not null" json:"labdip_no"`
	PartyName   string    `gorm:"not null" json:"party_name"`
	Quantity    float64   `gorm:"not null" json:"quantity"`
	Rate        float64   `gorm:"not null" json:"rate"`
	TotalAmount float64   `gorm:"not null" json:"total_amount"`
	Remarks     string    `json:"remarks"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
