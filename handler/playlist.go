package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

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
		array_to_string(track_list, ',') as track_list_string,
		track_list
	from
		playlists
	where
		id = ?`

	type oriPlaylist struct {
		models.Playlist
		TrackListString string `json:"track_list_string"`
	}
	originPlaylist := &oriPlaylist{}
	playlist := &serializer.Playlist{}

	// 플레이리스트 조회
	h.DB.Raw(sql, id).Scan(&originPlaylist).Scan(&playlist)

	// 트랙리스트 조회
	var trackList []string
	trackList = strings.Split(originPlaylist.TrackListString, ",")
	h.DB.Where("track_id IN (?)", trackList).Order("array_position(array[" + originPlaylist.TrackListString + "], track_id)").Find(&playlist.PlayList)

	for i := range playlist.PlayList {
		profile := &models.Profile{}
		track := &playlist.PlayList[i]
		h.DB.Model(playlist.PlayList[i]).Related(&track.Genre).Related(&track.API).Related(&track.User)
		h.DB.Where("user_id = ?", playlist.PlayList[i].UserID).Find(&profile)
		t := track.CreatedAt.In(time.FixedZone("KST", 9*60*60))
		track.CreatedAt = t
		track.Order = i + 1
		playlist.PlayList[i].User.Profile = *profile
	}

	return c.JSON(http.StatusOK, &playlist)
}

// CreatePlaylist 플레이리스트 생성
func (h *Handler) CreatePlaylist(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	playlistName := c.FormValue("playlist_name")
	playlist := &models.Playlist{}

	// 생성
	now := time.Now()
	query := "INSERT INTO public.playlists (user_id, playlist_name, created_at, updated_at) VALUES($1, $2, $3, $4) RETURNING *"
	h.DB.Raw(query, claims["user_id"].(float64), playlistName, now, now).Scan(&playlist)

	return c.JSON(http.StatusCreated, &playlist)
}

// UpdatePlaylist 플레이리스트 수정
func (h *Handler) UpdatePlaylist(c echo.Context) error {
	playlistID := c.Param("id")
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	playlist := &models.Playlist{}
	h.DB.Where("id = ?", playlistID).Find(&playlist)

	// 권한 체크
	if playlist.UserID != int(claims["user_id"].(float64)) {
		return echo.ErrForbidden
	}

	// 입력받은 파라미터
	inpTrackList := c.FormValue("track_list")
	inpPlaylistName := c.FormValue("playlist_name")

	// 플레이리스트 이름
	var playlistName string
	if inpPlaylistName == "" {
		playlistName = playlist.PlaylistName
	} else {
		playlistName = inpPlaylistName
	}

	// 플레이리스트 소속 트랙 리스트
	var updateTrackList string
	if inpTrackList == "" {
		updateTrackList = string(playlist.TrackList)
	} else {
		updateTrackList = "{" + inpTrackList + "}"
	}

	// 업데이트
	h.DB.Model(&playlist).Updates(map[string]interface{}{"track_list": updateTrackList, "playlist_name": playlistName})

	return c.JSON(http.StatusOK, playlist)
}

// FindMyPlaylist 내 플레이리스트 조회
func (h *Handler) FindMyPlaylist(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	playlistType := c.QueryParam("type")
	userID := int(claims["user_id"].(float64))
	var isDefault bool

	switch playlistType {
	case "another":
		isDefault = false

		// 기본값
	default:
		isDefault = true
	}

	if isDefault {
		sql := `select
			id,
			playlist_name,
			user_id,
			created_at,
			updated_at,
			array_to_string(track_list, ',') as track_list_string,
			track_list
		from
			playlists
		where
			user_id = ?
		and is_default = ?`

		type oriPlaylist struct {
			models.Playlist
			TrackListString string `json:"track_list_string"`
		}
		originPlaylist := &oriPlaylist{}
		playlist := &serializer.Playlist{}

		// 플레이리스트 조회
		h.DB.Raw(sql, userID, true).Scan(&originPlaylist).Scan(&playlist)

		// 플레이리스트가 없을때 생성
		if playlist.ID == 0 {
			playlistName := "기본 플레이리스트"

			// 생성
			now := time.Now()
			query := "INSERT INTO public.playlists (user_id, playlist_name, created_at, updated_at, is_default) VALUES($1, $2, $3, $4, true) RETURNING *"
			h.DB.Raw(query, userID, playlistName, now, now).Scan(&playlist)
		}

		fmt.Println("playlist", len(playlist.PlayList))

		// 트랙리스트 조회
		var trackList []string
		trackList = strings.Split(originPlaylist.TrackListString, ",")
		h.DB.Where("track_id IN (?)", trackList).Order("array_position(array[" + originPlaylist.TrackListString + "], track_id)").Find(&playlist.PlayList)

		for i := range playlist.PlayList {
			profile := &models.Profile{}
			track := &playlist.PlayList[i]
			h.DB.Model(playlist.PlayList[i]).Related(&track.Genre).Related(&track.API).Related(&track.User)
			h.DB.Where("user_id = ?", playlist.PlayList[i].UserID).Find(&profile)
			t := track.CreatedAt.In(time.FixedZone("KST", 9*60*60))
			track.CreatedAt = t
			track.Order = i + 1
			playlist.PlayList[i].User.Profile = *profile
		}

		return c.JSON(http.StatusOK, &playlist)
	}

	// 플레이리스트 그룹 리스트 조회
	anotherPlayList := []PlaylistSimple{}
	query := `
	select
		*,
		coalesce(array_length(track_list, 1), 0) as "count"
	from
		public.playlists
	where
		user_id = ?
		and is_default = false
	order by
		created_at asc
	`
	h.DB.Raw(query, userID).Scan(&anotherPlayList)
	return c.JSON(http.StatusOK, anotherPlayList)
}

// DeletePlayList 플레이리스트 삭제
func (h *Handler) DeletePlayList(c echo.Context) error {
	ID := c.Param("id")
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := int(claims["user_id"].(float64))

	h.DB.Delete(&models.Playlist{}, "id = ? and user_id = ? and is_default = false", ID, userID)

	// fmt.Println(result.)
	return nil
}

// PlaylistSimple ...
type PlaylistSimple struct {
	ID           uint      `gorm:"primary_key" json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	PlaylistName string    `json:"playlist_name"`
	Count        int       `json:"count"`
}
