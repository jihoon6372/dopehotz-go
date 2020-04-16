package utils

import (
	"os"

	"github.com/astaxie/beego/utils/pagination"
	"github.com/jihoon6372/dopehotz-go/config"
	"github.com/labstack/echo"
	"gopkg.in/yaml.v2"
)

// ReadConfig 서버 정보 조회 함수
func ReadConfig(cfg *config.Config) {
	f, err := os.Open("./config/config.yml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		panic(err)
	}
}

// UniqueInt 배열 중복값 제거
func UniqueInt(intSlice []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// ListPagination 리스트 페이징 처리
func ListPagination(c echo.Context, paginator pagination.Paginator, data interface{}) map[string]interface{} {
	host := c.Scheme() + "://" + c.Request().Host
	var previous *string
	if paginator.HasPrev() {
		prev := host + paginator.PageLinkPrev()
		previous = &prev
	}

	var next *string
	if paginator.HasNext() {
		n := host + paginator.PageLinkNext()
		next = &n
	}

	return map[string]interface{}{
		"next":     next,
		"previous": previous,
		"results":  data,
	}
}
