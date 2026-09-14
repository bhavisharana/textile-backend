package models

import (
	"time"
)

type Labdip struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	LabdipNo     string     `gorm:"not null" json:"labdip_no"`
	PartyName    string     `gorm:"not null" json:"party_name"`
	Status       string     `gorm:"not null" json:"status"`
	QualityID    *uint      `json:"quality_id,omitempty"`
	QualityName  string     `gorm:"not null" json:"quality_name"`
	ColorName    string     `gorm:"not null" json:"color_name"`
	ReceivedDate *time.Time `json:"received_date"`
	SendingDate  *time.Time `json:"sending_date"`
	Remarks      string     `json:"remarks"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
