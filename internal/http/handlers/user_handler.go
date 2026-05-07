package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"property-api/internal/dto"
	"property-api/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "users fetched",
		"data":    h.service.GetUsers(),
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
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "user not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to fetch user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "user fetched",
		"data":    user,
	})
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	payload := new(dto.CreateUserRequest)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	user := h.service.CreateUser(*payload)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "user created",
		"data":    user,
	})
}
