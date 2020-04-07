package config

import (
	"github.com/dgrijalva/jwt-go"
)

// JwtCustomClaims are custom claims extending default ones.
type JwtCustomClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"user_name"`
	Email    string `json:"email"`
	OrigIat  int    `json:"orig_iat"`
	jwt.StandardClaims
}

// Config 서버 세팅정보
type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yarm:"database"`
		SSLMode  string `yarm:"sslmode"`
	} `yaml:"database"`
	Config struct {
		SecretKey string `yaml:"secretkey"`
	} `yaml:"config"`
}
