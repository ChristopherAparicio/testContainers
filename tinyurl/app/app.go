package app

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	config "github.com/christapa/tinyurl/config"
	tinyHttp "github.com/christapa/tinyurl/internal/tinyurl/api/http"
	infra "github.com/christapa/tinyurl/internal/tinyurl/infra"
	services "github.com/christapa/tinyurl/internal/tinyurl/services"
	logger "github.com/christapa/tinyurl/pkg/logger"
	sql "github.com/christapa/tinyurl/pkg/sql"
)

func Run(config config.Config) {
	e := echo.New()

	databaseConn, err := sql.NewConn(sql.DBConfig{
		Host:     config.Database.Host,
		Port:     config.Database.Port,
		User:     config.Database.User,
		Password: config.Database.Password,
		Database: config.Database.Database,
	})
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer databaseConn.Close()

	repository := infra.NewUrlSqlRepository(databaseConn)
	service := services.NewUrlService(repository)
	handler := tinyHttp.NewHttpHandler(service)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	tinyHttp.RegisterHandlers(e, handler)

	address := fmt.Sprintf(":%d", config.Server.Port)
	e.Logger.Fatal(e.Start(address))
}
