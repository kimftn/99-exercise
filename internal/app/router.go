package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"property-api/internal/http/handlers"
	"property-api/internal/repository"
	"property-api/internal/service"
)

func NewServer() *fiber.App {
	panic("user repository is required; use NewServerWithPool")
}

func NewServerWithPool(pool *pgxpool.Pool) *fiber.App {
	if pool == nil {
		panic("postgres pool is required for user APIs")
	}

	return newServerWithDependencies(repository.NewPostgresUserRepository(pool))
}

func newServerWithDependencies(userRepository repository.UserRepository) *fiber.App {
	if userRepository == nil {
		panic("user repository is required")
	}

	app := fiber.New()

	listingRepository := repository.NewInMemoryListingRepository()
	listingService := service.NewListingService(listingRepository)
	userService := service.NewUserService(userRepository)
	publicService := service.NewPublicService(listingService, userService)

	listingHandler := handlers.NewListingHandler(listingService)
	userHandler := handlers.NewUserHandler(userService)
	publicHandler := handlers.NewPublicHandler(publicService)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "service is running",
		})
	})

	app.Post("/listings", listingHandler.CreateListing)
	app.Get("/listings", listingHandler.GetListings)

	app.Get("/users", userHandler.GetUsers)
	app.Get("/users/:id", userHandler.GetUserByID)
	app.Post("/users", userHandler.CreateUser)

	publicAPI := app.Group("/public-api")
	publicAPI.Get("/listings", publicHandler.GetListings)
	publicAPI.Post("/users", publicHandler.CreateUser)
	publicAPI.Post("/listings", publicHandler.CreateListing)

	return app
}
