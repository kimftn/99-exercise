package handlers

import (
	"github.com/gofiber/fiber/v2"

	"property-api/internal/dto"
	"property-api/internal/service"
)

type ListingHandler struct {
	service *service.ListingService
}

func NewListingHandler(service *service.ListingService) *ListingHandler {
	return &ListingHandler{service: service}
}

func (h *ListingHandler) CreateListing(c *fiber.Ctx) error {
	payload := new(dto.CreateListingRequest)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	listing := h.service.CreateListing(*payload)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "listing created",
		"data":    listing,
	})
}

func (h *ListingHandler) GetListings(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "listings fetched",
		"data":    h.service.GetListings(),
	})
}
