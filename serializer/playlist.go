package serializer

import (
	"github.com/jihoon6372/dopehotz-go/models"
)

// Playlist 시리얼라이저
type Playlist struct {
	models.Playlist
	TrackList []models.Track `json:"track_list"`
}
