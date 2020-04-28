package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/astaxie/beego/utils/pagination"
	"github.com/dgrijalva/jwt-go"
	"github.com/jihoon6372/dopehotz-go/models"
	"github.com/jihoon6372/dopehotz-go/serializer"
	"github.com/jihoon6372/dopehotz-go/utils"
	"github.com/labstack/echo"
)

var paginator = &pagination.Paginator{}

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
	h.DB.Create(eventInfo).Scan(result)
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

	h.DB.Delete(&eventInfo)

	return c.JSON(http.StatusNoContent, nil)
}

// FindEventList 공연정보 리스트
func (h *Handler) FindEventList(c echo.Context) error {
	// 최대 개수
	const perPage = 15

	// 현재시간
	now := time.Now()

	// 조회 조건
	where := "performance_time >= ?"

	// 카운트 조회
	var count int
	h.DB.Model(models.EventInfo{}).Where(where, now).Count(&count)

	// 페이징
	paginator = pagination.NewPaginator(c.Request(), perPage, count)

	// 리스트 조회
	eventInfos := []models.EventInfo{}
	h.DB.Limit(perPage).Offset(paginator.Offset()).Where(where, now).Order("performance_time", true).Find(&eventInfos)

	// 관계형 데이터 조회
	for i, eventInfo := range eventInfos {
		fmt.Println("event", i)
		user := &models.User{}
		profile := &models.Profile{}
		h.DB.First(&user, eventInfo.UserID)
		h.DB.Where("user_id = ?", eventInfo.UserID).Find(&profile)
		eventInfos[i].User = *user
		eventInfos[i].User.Profile = *profile
	}

	return c.JSON(http.StatusOK, utils.ListPagination(c, *paginator, eventInfos))
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

// FindUserEventList 사용자 공연정보 리스트
func (h *Handler) FindUserEventList(c echo.Context) error {
	userID := c.Param("userId")

	// 최대 개수
	const perPage = 15

	// 현재시간
	now := time.Now()

	// 조건
	where := "performance_time >= ? and user_id = ?"

	// 카운트 조회
	var count int
	h.DB.Model(models.EventInfo{}).Where(where, now, userID).Count(&count)

	// 페이징
	paginator = pagination.NewPaginator(c.Request(), perPage, count)

	// 리스트 조회
	eventInfos := []models.EventInfo{}
	h.DB.Limit(perPage).Offset(paginator.Offset()).Where(where, now, userID).Order("performance_time", true).Find(&eventInfos)

	// 관계형 데이터 조회
	for i, eventInfo := range eventInfos {
		user := &models.User{}
		profile := &models.Profile{}
		h.DB.First(&user, eventInfo.UserID)
		h.DB.Where("user_id = ?", eventInfo.UserID).Find(&profile)
		eventInfos[i].User = *user
		eventInfos[i].User.Profile = *profile
	}

	return c.JSON(http.StatusOK, utils.ListPagination(c, *paginator, eventInfos))
}
