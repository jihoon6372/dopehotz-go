package utils

import (
	"os"

	"github.com/jihoon6372/dopehotz-go/config"
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
