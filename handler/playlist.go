package handler

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/jihoon6372/dopehotz-go/models"
	"github.com/jihoon6372/dopehotz-go/serializer"
	"github.com/labstack/echo"
)

// FindPlaylist 플레이리스트 조회
func (h *Handler) FindPlaylist(c echo.Context) error {
	id := c.Param("id")

	sql := `select
		id,
		playlist_name,
		user_id,
		created_at,
		updated_at,
		array_to_string(track_list, ',') as track_list
	from
		playlists
	where
		id = ?`

	type oriPlaylist struct {
		models.Playlist
		TrackList string `json:"track_list"`
	}
	originPlaylist := &oriPlaylist{}
	playlist := &serializer.Playlist{}

	// 플레이리스트 조회
	h.DB.Raw(sql, id).Scan(&originPlaylist).Scan(&playlist)

	// 트랙리스트 조회
	var trackList []string
	trackList = strings.Split(originPlaylist.TrackList, ",")
	h.DB.Where("track_id IN (?)", trackList).Order("array_position(array[" + originPlaylist.TrackList + "], track_id)").Find(&playlist.TrackList)

	for i := range playlist.TrackList {
		h.DB.Model(playlist.TrackList[i]).Related(&playlist.TrackList[i].Genre).Related(&playlist.TrackList[i].API).Related(&playlist.TrackList[i].User)
	}

	return c.JSON(http.StatusOK, &playlist)
}

// CreatePlaylist 플레이리스트 생성
func (h *Handler) CreatePlaylist(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	playlistName := c.FormValue("playlist_name")

	playlist := &models.Playlist{}
	playlist.UserID = int(claims["user_id"].(float64))
	playlist.PlaylistName = playlistName
	// h.DB.Create(&playlist)

	return c.JSON(http.StatusCreated, &playlist)
}
