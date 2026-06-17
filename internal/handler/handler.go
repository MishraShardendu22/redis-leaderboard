package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	appredis "redis-leaderboard/internal/redis"
	"redis-leaderboard/internal/view"
)

type Handler struct {
	service *appredis.Service
}

func New(service *appredis.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	app.Get("/", h.Home)
	app.Get("/leaderboard", h.Leaderboard)
	app.Post("/player/:name/increment", h.Increment)
	app.Post("/player/:name/decrement", h.Decrement)
	app.Get("/health", h.Health)
}

func (h *Handler) Home(c *fiber.Ctx) error {
	players, movements, err := h.service.LeaderboardSnapshot(c.UserContext())
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	return view.Render(c, view.Layout("Redis Leaderboard", players, movements))
}

func (h *Handler) Leaderboard(c *fiber.Ctx) error {
	players, movements, err := h.service.LeaderboardSnapshot(c.UserContext())
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	return view.Render(c, view.Leaderboard(players, movements))
}

func (h *Handler) Increment(c *fiber.Ctx) error {
	if err := h.service.IncreaseScore(c.Params("name"), 10); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	return h.Leaderboard(c)
}

func (h *Handler) Decrement(c *fiber.Ctx) error {
	if err := h.service.DecreaseScore(c.Params("name"), 10); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	return h.Leaderboard(c)
}

func (h *Handler) Health(c *fiber.Ctx) error {
	if err := h.service.Ping(c.UserContext()); err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{"status": "unhealthy", "error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "ok"})
}
