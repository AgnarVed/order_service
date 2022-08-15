package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	log "github.com/sirupsen/logrus"
	"tests2/internal/cache"
	"tests2/internal/config"
	"tests2/internal/models"
	"tests2/internal/service"
)

type Handler struct {
	cfg      *config.Config
	services *service.Service
	cache    *cache.Cache
}

func (h *Handler) Init(app *fiber.App) {
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${ip}:${port} ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "2006-01-02T15:04:05.000Z",
		Output:     log.StandardLogger().Out,
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Get("/ping", h.serviceHealth)
	app.Get("/:orderID", h.getOrderByID)
	app.Get("/orders/list", h.getOrderList)
	app.Get("/cache/:orderID", h.getOrderFromCacheByID)
	app.Post("/order/post", h.validateOrder, h.putOrderInCache)
	//app.Post("/order/create", h.createOrder)
	//v1 := app.Group("/")
	//
	//h.initRoutesItem(v1)

}

func (h *Handler) serviceHealth(ctx *fiber.Ctx) error {

	msg := "pong"
	return ctx.Status(200).SendString(msg)
}

func (h *Handler) Response(ctx *fiber.Ctx, statusCode int, data interface{}, err error) error {
	resp := models.Response{
		ErrorText: "",
		HasError:  false,
		Resp:      data,
	}

	if data == nil {
		resp.Resp = struct{}{}
	}

	if err != nil {
		resp.HasError = true
		resp.ErrorText = err.Error()
		log.Errorf("%s", err.Error())
	}
	return ctx.Status(statusCode).JSON(&resp)
}

func NewHandlers(cfg *config.Config, services *service.Service, cache *cache.Cache) *Handler {
	return &Handler{
		cfg:      cfg,
		services: services,
		cache:    cache,
	}
}
