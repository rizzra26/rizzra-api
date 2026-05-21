package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizzra/api/internal/cloudinary"
	"github.com/rizzra/api/internal/handlers"
	"github.com/rizzra/api/internal/middleware"
	"github.com/rizzra/api/internal/repository"
)

func Setup(app *fiber.App, pool *pgxpool.Pool, jwtSecret string, cld *cloudinary.Service) {
	app.Use(middleware.CORS())

	authHandler := handlers.NewAuthHandler(pool, jwtSecret)
	dashboardHandler := handlers.NewDashboardHandler(pool)
	letterHandler := handlers.NewLetterHandler(pool)
	projectHandler := handlers.NewProjectHandler(pool)
	stackHandler := handlers.NewStackHandler(pool)

	projectRepo := repository.NewProjectRepo(pool)
	uploadHandler := handlers.NewUploadHandler(projectRepo, cld)

	// Public routes (no auth)
	app.Post("/api/v1/auth/login", authHandler.Login)
	app.Post("/api/v1/auth/refresh", authHandler.Refresh)

	// Public GET endpoints (portfolio website)
	public := app.Group("/api/v1")
	public.Get("/letters", letterHandler.List)
	public.Get("/letters/:id", letterHandler.Get)
	public.Get("/projects", projectHandler.List)
	public.Get("/projects/:id", projectHandler.Get)
	public.Get("/stack", stackHandler.ListItems)

	// Protected routes (admin panel)
	auth := app.Group("/api/v1")
	auth.Use(middleware.Auth(jwtSecret))

	auth.Get("/dashboard/stats", dashboardHandler.Stats)

	// Letters (mutations)
	auth.Post("/letters", letterHandler.Create)
	auth.Put("/letters/:id", letterHandler.Update)
	auth.Delete("/letters/:id", letterHandler.Delete)

	// Projects (mutations)
	auth.Post("/projects/reorder", projectHandler.Reorder)
	auth.Post("/projects", projectHandler.Create)
	auth.Put("/projects/:id", projectHandler.Update)
	auth.Delete("/projects/:id", projectHandler.Delete)

	// Stack (mutations)
	auth.Post("/stack", stackHandler.CreateItem)
	auth.Put("/stack/:id", stackHandler.UpdateItem)
	auth.Delete("/stack/:id", stackHandler.DeleteItem)

	// Upload
	auth.Post("/upload/cover", uploadHandler.Cover)
}
