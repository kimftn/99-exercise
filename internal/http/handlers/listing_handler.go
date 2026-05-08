package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"property-api/internal/dto"
	httphelpers "property-api/internal/http/helpers"
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

	if payload.UserID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid user_id",
		})
	}

	listing, err := h.service.CreateListing(*payload)
	if err != nil {
		if err == service.ErrInvalidListingType || err == service.ErrListingUserNotFound {
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

func (h *ListingHandler) GetListings(c *fiber.Ctx) error {
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
		parsedUserID, err := strconv.Atoi(c.Query("user_id"))
		if err != nil || parsedUserID <= 0 {
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
