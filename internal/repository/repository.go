package repository

import (
	"context"
	"database/sql"
	"tests2/internal/models"
	"tests2/internal/repository/client"
)

type Repositories struct {
	OrderDB  OrderDB
	CommonDB CommonDB
}

type CommonDB interface {
	Get() *sql.DB
	BeginTransaction(ctx context.Context) (*sql.Tx, error)
	CommitTransaction(ctx context.Context, tx *sql.Tx) error
	RollbackTransaction(ctx context.Context, tx *sql.Tx) error
}

type OrderDB interface {
	GetOrderByID(ctx context.Context, tx *sql.Tx, orderID string) (*models.Order, error)
	GetOrderList(ctx context.Context, tx *sql.Tx) ([]*models.OrderDB, error)
	CreateOrder(ctx context.Context, tx *sql.Tx, insert *models.Order) error
}

type Cache interface {
}

func NewRepositories(psqlClient *client.PostgresClient) *Repositories {
	return &Repositories{
		OrderDB:  NewOrderDB(),
		CommonDB: NewCommonRepo(*psqlClient),
	}
}
