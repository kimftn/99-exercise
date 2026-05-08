package handlers

import (
	"github.com/gofiber/fiber/v2"

	"property-api/internal/dto"
	httphelpers "property-api/internal/http/helpers"
	"property-api/internal/service"
)

type PublicHandler struct {
	service *service.PublicService
}

func NewPublicHandler(service *service.PublicService) *PublicHandler {
	return &PublicHandler{service: service}
}

func (h *PublicHandler) GetListings(c *fiber.Ctx) error {
	pageNum, err := httphelpers.ParsePositiveQueryInt(c.Query("page_num"), 1)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid page_num",
		})
	}

	pageSize, err := httphelpers.ParsePositiveQueryInt(c.Query("page_size"), 10)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid page_size",
		})
	}

	var userID *int
	if c.Query("user_id") != "" {
		parsedUserID, err := httphelpers.ParsePositiveQueryInt(c.Query("user_id"), 0)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "invalid user_id",
			})
		}
		userID = &parsedUserID
	}

	listings, err := h.service.GetListings(pageNum, pageSize, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to fetch listings",
		})
	}

	return c.JSON(dto.GetListingsResponse{
		Result:   true,
		Listings: listings,
	})
}

func (h *PublicHandler) CreateListing(c *fiber.Ctx) error {
	payload := new(dto.PublicCreateListingRequest)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	if payload.UserID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid user_id",
		})
	}

	listing, err := h.service.CreateListing(*payload)
	if err != nil {
		if err == service.ErrPublicListingRejected {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(dto.ListingErrorResponse{
				Return:   false,
				Listings: nil,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to create listing",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.CreateListingResponse{
		Result:  true,
		Listing: &listing,
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
