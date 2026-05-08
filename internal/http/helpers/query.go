package helpers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func ParsePositiveQueryInt(value string, fallback int) (int, error) {
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return 0, fiber.ErrBadRequest
	}

	return parsed, nil
}
