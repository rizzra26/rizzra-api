package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v3"
	"github.com/rizzra/api/internal/cloudinary"
	"github.com/rizzra/api/internal/config"
	"github.com/rizzra/api/internal/database"
	"github.com/rizzra/api/internal/router"
	"github.com/rizzra/api/internal/util"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	cld, err := cloudinary.New(cfg.CloudinaryCloudName, cfg.CloudinaryAPIKey, cfg.CloudinaryAPISecret)
	if err != nil {
		log.Fatalf("failed to init cloudinary: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := database.Connect(ctx, cfg.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := database.RunMigrations(ctx, pool); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	hash, err := util.HashPassword("admin123")
	if err != nil {
		log.Fatalf("failed to hash admin password: %v", err)
	}
	if err := database.SeedAdmin(ctx, pool, "admin@rizzra.dev", "rizzra", hash); err != nil {
		log.Fatalf("failed to seed admin: %v", err)
	}

	app := fiber.New(fiber.Config{
		AppName:      "Rizzra API",
		ErrorHandler: customErrorHandler,
	})

	router.Setup(app, pool, cfg.JWTSecret, cld)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		slog.Info("shutting down server...")
		cancel()
		app.Shutdown()
	}()

	addr := fmt.Sprintf(":%d", cfg.Port)
	slog.Info("starting server", "addr", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func customErrorHandler(c fiber.Ctx, err error) error {
	code := 500
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
