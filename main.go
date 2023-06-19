package main

import (
	"github.com/ar-siddiqui/mcat-ras/config"
	"github.com/ar-siddiqui/mcat-ras/handlers"
	"github.com/ar-siddiqui/mcat-ras/pgdb"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	// echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {

	// Connect to backend services
	appConfig := config.Init()
	dbConfig := pgdb.DBInit()

	// Instantiate echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())

	// HealthCheck
	e.GET("/ping", handlers.Ping(appConfig.FileStore))

	// Swagger
	// e.GET("/swagger/*", echoSwagger.WrapHandler)

	// ras endpoints
	// these endpoints create a Ras Model struct from files
	// and then apply receiver functions to the struct to answer desired question
	e.GET("/isamodel", handlers.IsAModel(appConfig.FileStore))
	e.GET("/modeltype", handlers.ModelType(appConfig.FileStore))
	e.GET("/modelversion", handlers.ModelVersion(appConfig.FileStore))
	e.GET("/index", handlers.Index(appConfig.FileStore))
	e.GET("/isgeospatial", handlers.IsGeospatial(appConfig.FileStore))
	e.GET("/geospatialdata", handlers.GeospatialData(appConfig))
	e.GET("/forcingdata", handlers.ForcingData(appConfig))

	// pgdb endpoints
	e.POST("/upsert/model", pgdb.UpsertRasModel(appConfig, dbConfig))
	e.POST("/upsert/geometry", pgdb.UpsertRasGeometry(appConfig, dbConfig))
	e.POST("/refresh", pgdb.RefreshRasViews(dbConfig))
	e.POST("/vacuum", pgdb.VacuumRasViews(dbConfig))

	e.Logger.Fatal(e.Start(appConfig.Address()))
}
