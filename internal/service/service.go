package service

import (
	"context"
	"tests2/internal/config"
	"tests2/internal/models"
	"tests2/internal/repository"
)

type Service struct {
	Order Order
}

type Order interface {
	GetOrderByID(ctx context.Context, orderID string) (*models.Order, error)
	CreateOrder(ctx context.Context, insert *models.Order) error
	GetOrderList(ctx context.Context) ([]*models.OrderDB, error)
}

func NewService(repos *repository.Repositories, cfg *config.Config) *Service {
	return &Service{
		Order: NewOrderService(repos, cfg),
	}
}
