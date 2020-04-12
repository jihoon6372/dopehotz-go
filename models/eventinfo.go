package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// EventInfo 공연 정보 모델
type EventInfo struct {
	gorm.Model
	Address         string
	PerformanceTime time.Time `json:"performance_time"`
	PerformanceName string    `json:"performance_name"`
	Link            string
	UserID          int  `json:"user_id"`
	User            User `gorm:"foreignkey:UserID" json:"user"`
}
