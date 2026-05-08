package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"property-api/internal/client"
	"property-api/internal/http/handlers"
	"property-api/internal/repository"
	"property-api/internal/service"
)

func NewServer() *fiber.App {
	panic("user repository is required; use NewServerWithPool")
}

func NewServerWithPool(pool *pgxpool.Pool, listingServiceBaseURL string) *fiber.App {
	if pool == nil {
		panic("postgres pool is required for user and listing APIs")
	}

	return newServerWithDependencies(
		repository.NewPostgresUserRepository(pool),
		repository.NewPostgresListingRepository(pool),
		client.NewHTTPListingClient(listingServiceBaseURL),
	)
}

func newServerWithDependencies(
	userRepository repository.UserRepository,
	listingRepository repository.ListingRepository,
	publicListingClient client.ListingClient,
) *fiber.App {
	if userRepository == nil || listingRepository == nil || publicListingClient == nil {
		panic("user repository, listing repository, and public listing client are required")
	}

	app := fiber.New()
	listingService := service.NewListingService(listingRepository, userRepository)
	userService := service.NewUserService(userRepository)
	publicService := service.NewPublicService(publicListingClient, userService)

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
