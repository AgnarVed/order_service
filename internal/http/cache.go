package http

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"tests2/internal/models"
)

func (h *Handler) getOrderFromCacheByID(ctx *fiber.Ctx) error {
	orderIdStr := ctx.Params("orderID")
	orderByte, _, err := h.cache.Get(orderIdStr)
	order := models.Order{}
	err = json.Unmarshal(orderByte, &order)
	if err != nil {
		return h.Response(ctx, http.StatusInternalServerError, nil, err)
	}
	if err != nil {
		return h.Response(ctx, http.StatusBadRequest, nil, err)
	}
	return h.Response(ctx, http.StatusOK, order, nil)
}

func (h *Handler) putOrderInCache(ctx *fiber.Ctx) error {
	order := models.Order{}
	err := ctx.BodyParser(&order)
	if err != nil {
		return h.Response(ctx, http.StatusBadRequest, nil, err)
	}
	ok := h.cache.Add(order.OrderUID, order)
	if ok {
		return h.Response(ctx, http.StatusOK, nil, nil)
	}
	err = h.services.Order.CreateOrder(ctx.UserContext(), &order)
	if err != nil {
		return h.Response(ctx, http.StatusBadRequest, nil, err)
	}
	return h.Response(ctx, http.StatusOK, nil, nil)
}
