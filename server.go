package main

import (
	"fmt"

	"github.com/jihoon6372/dopehotz-go/config"
	"github.com/jihoon6372/dopehotz-go/handler"
	"github.com/jihoon6372/dopehotz-go/utils"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	var cfg config.Config
	utils.ReadConfig(&cfg)

	dbinfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.Username, cfg.Database.Database, cfg.Database.Password, cfg.Database.SSLMode)

	db, err := gorm.Open("postgres", dbinfo)
	defer db.Close()
	if err != nil {
		panic(err)
	}

	e := echo.New()
	// jwtConfig := getJWTConfig([]byte(cfg.Config.SecretKey))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	h := &handler.Handler{DB: db}
	e.GET("/users/:id", h.FindUser)

	e.Logger.Fatal(e.Start(":8080"))
}
