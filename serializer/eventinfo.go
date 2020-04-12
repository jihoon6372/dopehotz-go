package serializer

import "time"

// EventInfoBase 공연정보 기본 시리얼라이저
type EventInfoBase struct {
	ID              uint      `gorm:"primary_key" json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	Address         string    `json:"address"`
	PerformanceTime time.Time `json:"performance_time"`
	PerformanceName string    `json:"performance_name"`
	Link            string    `json:"link"`
	UserID          int       `json:"user_id"`
}
