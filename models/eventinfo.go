package models

import (
	"time"
)

// EventInfo 공연 정보 모델
type EventInfo struct {
	ID              uint       `gorm:"primary_key" json:"id"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
	Address         string     `json:"address"`
	PerformanceTime time.Time  `json:"performance_time"`
	PerformanceName string     `json:"performance_name"`
	Link            string     `json:"link"`
	UserID          int        `json:"user_id"`
	User            User       `gorm:"foreignkey:UserID" json:"user"`
}
