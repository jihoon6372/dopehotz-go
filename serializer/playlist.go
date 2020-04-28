package serializer

import (
	"github.com/jihoon6372/dopehotz-go/models"
)

// Playlist 시리얼라이저
type Playlist struct {
	models.Playlist
	PlayList []models.Track `json:"play_list"`
}
