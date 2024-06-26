package app

import (
	"testovoe_2/config"
	"testovoe_2/internal/controller"
	"testovoe_2/internal/repository"
	"testovoe_2/internal/service"
	"testovoe_2/internal/slogger"
	"testovoe_2/pkg/postgres"
	"time"

	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
)

const (
	ReadTimeout  = 3 * time.Second
	WriteTimeout = 3 * time.Second
)

func Run() {
	slogger.SetLogger()

	slog.Info("init config")
	cfg := config.NewConfig()
	slog.Info("config ok")

	slog.Info("connecting to postgres")
	db := postgres.New(cfg.URL)
	defer db.Close()
	slog.Info("connect to postgres ok")

	slog.Info("init repositories")
	repositories := repository.NewRepositories(db)

	slog.Info("init services")
	deps := service.ServicesDeps{
		Repository:       repositories,
		SecretKeyAccess:  cfg.SecretKeyAccess,
		SecretKeyRefresh: cfg.SecretKeyRefresh,
	}

	services := service.NewServices(deps)

	fiberConfig := fiber.Config{
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
	}

	app := fiber.New(fiberConfig)
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://localhost:5173/",
	}))
	controller.NewRouter(app, services)

	slog.Info("starting fiber server")
	slog.Fatal(app.Listen(":" + cfg.Port))
}
