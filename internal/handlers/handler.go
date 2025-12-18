package handlers

import (
	"time"

	"github.com/AKA333/URLShortner/internal/models"
	"github.com/AKA333/URLShortner/internal/service"
	"github.com/gofiber/fiber/v2"
)

type URLHandler struct {
	urlService *service.URLService
}

func NewURLHandler(urlService *service.URLService) *URLHandler {
	return &URLHandler{
		urlService: urlService,
	}
}

func (h *URLHandler) ShortenURL(c *fiber.Ctx) error {
	var req models.ShortenRequest

	if err := c.BodyParser(&req); err != nil {
		return  c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request payload",
		})
	}

	if req.LongURL == "" {
		return  c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "long_url is required",
		})
	}

	response, err := h.urlService.ShortenURL(c.Context(), &req)
	if err != nil {
		if err.Error() == "custom alias already in use" {
			return  c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func (h *URLHandler) RedirectURL(c *fiber.Ctx) error {
	shortCode := c.Params("shortCode")
	if shortCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "short code is required",
		})
	}

	url, err := h.urlService.RedirectURL(c.Context(), shortCode)
	if err != nil {
		if err.Error() == "URL not found or expired" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "URL not found or has expired",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Redirect(url, fiber.StatusFound)
}

func (h *URLHandler) GetURLStats(c *fiber.Ctx) error {
	shortCode := c.Params("shortCode")
	if shortCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "short code is required",
		})
	}

	stats, err := h.urlService.GetURLStats(c.Context(), shortCode)
	if err != nil {
		if err.Error() == "URL not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "URL not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// return c.Status(fiber.StatusOK).JSON(stats)
	return c.JSON(stats)
}

func (h *URLHandler) HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "OK",
		"service": "URL Shortener",
		"time": time.Now().UTC(),
	})
}