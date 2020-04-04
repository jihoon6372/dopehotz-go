package handler

import (
	"net/http"

	"github.com/jihoon6372/dopehotz-go/serializer"
	"github.com/labstack/echo"
)

// Content 반환 (테스트용)
type Content struct {
	Message string
}

// FindUser 사용자 정보 조회
func (h *Handler) FindUser(c echo.Context) error {
	id := c.Param("id")
	user := &serializer.User{}
	h.DB.Model(&user).Where("id = ?", id).Scan(&user)
	h.DB.Model(&user.Profile).Where("user_id = ?", id).Scan(&user.Profile)
	return c.JSON(http.StatusOK, user)
}
