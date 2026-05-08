package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"property-api/internal/dto"
	httphelpers "property-api/internal/http/helpers"
	"property-api/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
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

	users, err := h.service.GetUsers(pageNum, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to fetch users",
		})
	}

	return c.JSON(dto.GetUsersResponse{
		Result: true,
		Users:  users,
	})
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid user id",
		})
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		if err == service.ErrUserNotFound {
			return c.Status(fiber.StatusOK).JSON(dto.GetUserByIDResponse{
				Result: false,
				Users:  nil,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to fetch user",
		})
	}

	return c.JSON(dto.GetUserByIDResponse{
		Result: true,
		Users:  &user,
	})
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	payload := new(dto.CreateUserRequest)
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
