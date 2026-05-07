package app

import (
	"github.com/gofiber/fiber/v2"

	"property-api/internal/http/handlers"
	"property-api/internal/repository"
	"property-api/internal/service"
	"property-api/internal/store"
)

func NewServer() *fiber.App {
	app := fiber.New()
	dataStore := store.New()

	listingRepository := repository.NewInMemoryListingRepository(dataStore)
	userRepository := repository.NewInMemoryUserRepository(dataStore)

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
