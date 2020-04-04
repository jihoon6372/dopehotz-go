package serializer

import (
	"github.com/jihoon6372/dopehotz-go/models"
)

// User 사용자 시리얼라이저
type User struct {
	models.AuthUser
	Profile models.Profile `json:"profile"`
}
