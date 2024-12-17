package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/liza/labwork_45/docs"
	configuration "github.com/liza/labwork_45/internal/config"
	"github.com/liza/labwork_45/internal/handlers"
	"github.com/liza/labwork_45/internal/middleware"
	"github.com/liza/labwork_45/internal/repository"
	"github.com/liza/labwork_45/internal/service"
)

// NewDBPsql function provides Connection with PostgreSQL database
func NewDBPsql(env string) (*pgxpool.Pool, error) {
	// Initialization a connect configuration for a PostgreSQL using pgx driver
	config, err := pgxpool.ParseConfig(env)
	if err != nil {
		return nil, fmt.Errorf("error connection to PostgreSQL: %v", err)
	}

	// Establishing a new connection to a PostgreSQL database using the pgx driver
	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("error connection to PostgreSQL: %v", err)
	}
	// Output to console
	fmt.Println("Connected to PostgreSQL!")

	return pool, nil
}

// @title Lab
// @version 1.0
// @description Diploma Documentation.

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	e := echo.New()

	cfg, err := configuration.NewConfig()
	if err != nil {
		fmt.Printf("Error extracting env variables: %v", err)
		return
	}

	pool, err := NewDBPsql(cfg.PgxDBAddr)
	if err != nil {
		err := echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("error creating database connection with PostgreSQL: %w", err))
		e.Logger.Fatal(err)
	}

	rps := repository.NewPsqlConnection(pool)

	auth := e.Group("/auth")
	{

		srv := service.NewAuthApiService(rps)
		handler := handlers.NewAuthApiHandler(srv)

		auth.GET("/getall", handler.GetAll, middleware.UserIdentity())
		auth.POST("/login", handler.Login)
		auth.POST("/signup", handler.SignUp)
		auth.POST("/refreshtokenpair", handler.RefreshTokenPair)
		auth.GET("/getpersonalinfo", handler.GetPersonalInfo, middleware.UserIdentity())
		auth.DELETE("/delete", handler.DeleteUser, middleware.UserIdentity())
	}
	courier := e.Group("/courier")
	{
		srv := service.NewCourierService(rps)
		handler := handlers.NewCourierHandler(srv)

		courier.PATCH("/updatecourier", handler.UpdateCourier, middleware.CourierIdentity())
		courier.GET("/getalldeliveries", handler.GetAlldeliveries, middleware.CourierIdentity())
		courier.PATCH("/choose_availible_delivery", handler.ChooseAvailibleDelivery, middleware.CourierIdentity())
		courier.PATCH("/update_delivery_status", handler.UpdateDeliveryStatus, middleware.CourierIdentity())
	}

	delivery := e.Group("/delivery")
	{
		srv := service.NewCourierService(rps)
		handler := handlers.NewCourierHandler(srv)

		delivery.POST("/create_delivary", handler.CreateDelivery, middleware.AdminIdentity())
	}
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":8080"))

}
