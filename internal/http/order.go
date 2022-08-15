package http

import (
	"database/sql"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"tests2/internal/models"
)

func (h *Handler) getOrderByID(ctx *fiber.Ctx) error {
	orderIdStr := ctx.Params("orderID")
	order, err := h.services.Order.GetOrderByID(ctx.UserContext(), orderIdStr)
	if err != nil {
		return h.Response(ctx, http.StatusBadRequest, nil, err)
	}
	return h.Response(ctx, http.StatusOK, order, nil)
}

func (h *Handler) getOrderList(ctx *fiber.Ctx) error {
	orderList, err := h.services.Order.GetOrderList(ctx.UserContext())
	if err != nil {
		if err == sql.ErrNoRows {
			return h.Response(ctx, http.StatusOK, orderList, nil)
		}
		return h.Response(ctx, http.StatusInternalServerError, nil, err)
	}
	return h.Response(ctx, http.StatusOK, orderList, nil)
}

var Validator = validator.New()

func (h *Handler) validateOrder(ctx *fiber.Ctx) error {
	var errors []*models.IError
	body := new(models.Order)
	ctx.BodyParser(&body)

	err := Validator.Struct(body)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var el models.IError
			el.Field = err.Field()
			el.Tag = err.Tag()
			el.Value = err.Param()
			errors = append(errors, &el)
		}
		return h.Response(ctx, http.StatusBadRequest, errors, err)
	}
	return ctx.Next()
}
