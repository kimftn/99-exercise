package handlers

import (
	"github.com/gofiber/fiber/v2"

	"property-api/internal/dto"
	"property-api/internal/service"
)

type PublicHandler struct {
	service *service.PublicService
}

func NewPublicHandler(service *service.PublicService) *PublicHandler {
	return &PublicHandler{service: service}
}

func (h *PublicHandler) GetListings(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "listings fetched",
		"data":    h.service.GetListings(),
	})
}

func (h *PublicHandler) CreateListing(c *fiber.Ctx) error {
	payload := new(dto.PublicCreateListingRequest)
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

func (h *PublicHandler) CreateUser(c *fiber.Ctx) error {
	payload := new(dto.PublicCreateUserRequest)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	user, err := h.service.CreateUser(*payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.CreateUserResponse{
		Result: true,
		User:   &user,
	})
}
