package handler

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jihoon6372/dopehotz-go/models"
	"github.com/jihoon6372/dopehotz-go/serializer"
	"github.com/labstack/echo"
)

// CreateEvent 공연정보 생성
func (h *Handler) CreateEvent(c echo.Context) error {
	// 등록 데이터
	address := c.FormValue("address")
	performanceName := c.FormValue("performance_name")
	performanceTime := c.FormValue("performance_time")
	link := c.FormValue("link")

	if address == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "address is required.",
		})
	}

	if performanceName == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "performance_name is required.",
		})
	}

	if performanceTime == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "performance_time is required.",
		})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	pt, ptErr := time.Parse(
		"2006-01-02 15:04:05",
		performanceTime)

	if ptErr != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "The time format is incorrect.",
		})
	}

	eventInfo := &models.EventInfo{
		UserID:          int(claims["user_id"].(float64)),
		Address:         address,
		PerformanceName: performanceName,
		PerformanceTime: pt,
		Link:            link,
	}

	result := &serializer.EventInfoBase{}
	h.DB.Create(eventInfo).Scan(&result)

	return c.JSON(http.StatusCreated, result)
}

// DeleteEvent 공연 정보 삭제
func (h *Handler) DeleteEvent(c echo.Context) error {
	ID := c.Param("id")
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := int(claims["user_id"].(float64))

	eventInfo := &models.EventInfo{}
	h.DB.First(&eventInfo, ID)
	if eventInfo.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"message": "permission denied.",
		})
	}

	eventInfo.DeletedAt = time.Now()
	h.DB.Save(&eventInfo)

	return c.JSON(http.StatusNoContent, nil)
}

// FindEventList 공연정보 리스트
func (h *Handler) FindEventList(c echo.Context) error {
	eventInfos := []models.EventInfo{}
	h.DB.Find(&eventInfos)

	for i, eventInfo := range eventInfos {
		user := &models.User{}
		profile := &models.Profile{}
		h.DB.First(&user, eventInfo.UserID)
		h.DB.Where("user_id = ?", eventInfo.UserID).Find(&profile)
		eventInfos[i].User = *user
		eventInfos[i].User.Profile = *profile
	}

	return c.JSON(http.StatusOK, eventInfos)
}

// FindEvent 공연정보 상세
func (h *Handler) FindEvent(c echo.Context) error {
	ID := c.Param("id")
	eventInfo := &models.EventInfo{}
	h.DB.First(&eventInfo, ID).Related(&eventInfo.User)

	// ID가 0이면 없는값
	if eventInfo.ID == 0 {
		return echo.ErrNotFound
	}

	// 사용자 프로필 조회
	profile := &models.Profile{}
	h.DB.Where("user_id = ?", eventInfo.UserID).Find(&profile)
	eventInfo.User.Profile = *profile

	return c.JSON(http.StatusOK, &eventInfo)
}
