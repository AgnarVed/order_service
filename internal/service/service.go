package service

import (
	"context"
	"tests2/internal/config"
	"tests2/internal/minio"
	"tests2/internal/models"
	"tests2/internal/repository"
)

type Service struct {
	Order Order
	Doc   Doc
}

type Order interface {
	GetOrderByID(ctx context.Context, orderID string) (*models.Order, error)
	CreateOrder(ctx context.Context, insert *models.Order) error
	GetOrderList(ctx context.Context) ([]*models.OrderDB, error)
}

type Doc interface {
	DownloadDoc(ctx context.Context, filename string) ([]byte, error)
	UploadDoc(ctx context.Context, fileInfo *models.FileInfo) error
}

func NewService(repos *repository.Repositories, cfg *config.Config, minio *minio.MinIOClient) *Service {
	return &Service{
		Order: NewOrderService(repos, cfg),
		Doc:   NewDocService(minio),
	}
}
