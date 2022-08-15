package http

import (
	"context"
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"tests2/internal/cache"
	"tests2/internal/config"
	"tests2/internal/repository"
	"tests2/internal/repository/client"
	"tests2/internal/service"
	"time"
)

func setup() *fiber.App {
	cfg := config.Config{
		Port:            70707,
		DBConnStr:       "host=localhost port=5432 user=valek password=password dbname=test_database sslmode=disable",
		DriverName:      "postgres",
		ShutdownTimeout: 1000,
	}
	db, err := sql.Open(cfg.DriverName, cfg.DBConnStr)
	if err != nil {
		logrus.Error("cannot connect to db", err)
	}
	clnt := client.NewPostgresClient(db)
	repos := repository.NewRepositories(&clnt)
	services := service.NewService(repos, &cfg)
	c, err := cache.NewCache(1024)
	if err != nil {
		logrus.Error("cannot init cache")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.ShutdownTimeout)*time.Second)
	defer cancel()
	orders, err := services.Order.GetOrderList(ctx)
	if err != nil {
		logrus.Fatal("Cannot upload cache ", err)
	}
	err = c.UploadCache(orders)
	if err != nil {
		logrus.Fatal("cannot upload cache")
	}

	h := NewHandlers(&cfg, services, c)

	app := fiber.New()

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
		Output: nil,
	}))
	h.Init(app)

	return app
}
