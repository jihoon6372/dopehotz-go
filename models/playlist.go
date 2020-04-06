package models

import "time"

// Playlist 플레이리스트 모델
type Playlist struct {
	ID           uint      `gorm:"primary_key" json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	UserID       int       `json:"user_id"`
	PlaylistName string    `json:"playlist_name"`
}
