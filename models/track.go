package models

import "time"

// Track 트랙 모델
type Track struct {
	ID               int       `json:"id"`
	TrackID          int       `json:"track_id"`
	Title            string    `gorm:"type:varchar(100);" json:"title"`
	Slug             string    `gorm:"type:varchar(100);" json:"slug"`
	TapeInfo         string    `json:"tape_info"`
	Duration         int       `json:"duration"`
	Lyrics           string    `json:"lyrics"`
	Tag              string    `json:"tag"`
	ImageURL         string    `json:"image_url"`
	DownloadURL      string    `json:"download_url"`
	WaveformURL      string    `json:"waveform_url"`
	ViewCount        int       `json:"view_count"`
	LikeCount        int       `json:"like_count"`
	PlayCount        int       `json:"play_count"`
	TempPlayCount    int       `json:"temp_play_count"`
	TrackScore       int       `json:"track_score"`
	OnStage          int       `json:"on_stage"`
	CommentCount     int       `json:"comment_count"`
	IsPublic         bool      `json:"is_public"`
	IsDistribute     bool      `json:"is_distribute"`
	IsDeleted        bool      `json:"is_deleted"`
	IsLike           bool      `json:"is_like"`
	IsAgreeMarketing string    `json:"is_agree_marketing"`
	BuyLink          string    `json:"buy_link"`
	CreatedAt        time.Time `json:"created_at"`
	GenreID          int       `json:"genre_id"`
	APIID            int       `json:"api_id"`
	UserID           int       `json:"user_id"`
	Genre            Genre     `json:"genre"`
	API              API       `json:"api"`
	User             User      `gorm:"foreignkey:UserID" json:"user"`
	Order            int       `json:"order"`
}

// TableName 유저 테이블명 반환
func (Track) TableName() string {
	return "tracks_track"
}

// Genre 장르 모델
type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// TableName 장르 테이블명
func (Genre) TableName() string {
	return "tracks_genre"
}

// API 모델
type API struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// TableName API 테이블명
func (API) TableName() string {
	return "tracks_trackapilist"
}
