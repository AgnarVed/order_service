package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"tests2/internal/models"
)

type order struct {
}

func NewOrderDB() OrderDB {
	return &order{}
}

func (od *order) GetOrderByID(ctx context.Context, tx *sql.Tx, orderID string) (*models.Order, error) {
	query := `SELECT order_uid, data
	FROM orders
	WHERE order_uid=$1;`

	row := tx.QueryRowContext(ctx, query, orderID)
	orderDB := models.OrderDB{}
	err := row.Scan(
		&orderDB.OrderUID,
		&orderDB.Data,
	)
	order, err := convertOrderFromDB(&orderDB)

	if err != nil {
		return nil, err
	}
	return order, nil
}

func (od *order) CreateOrder(ctx context.Context, tx *sql.Tx, insert *models.Order) error {
	query := `INSERT INTO orders VALUES ($1,$2)`
	orderData, err := json.Marshal(insert)
	_, err = tx.ExecContext(ctx, query, insert.OrderUID, orderData)
	if err != nil {
		return err
	}
	return nil
}

// GetOrderList gets all orders from Database
//func (od *order) GetOrderList(ctx context.Context, tx *sql.Tx) ([]*models.OrderDB, error) {
//	querySelect := sq.Select(`order_uid`, `data`).From(`orders`)
//	query, args, err := querySelect.PlaceholderFormat(sq.Dollar).ToSql()
//	if err != nil {
//		return nil, err
//	}
//
//	rows, err := tx.QueryContext(ctx, query, args...)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	orderList := make([]*models.OrderDB, 0)
//
//	for rows.Next() {
//		tmpOrder := models.OrderDB{}
//		err = rows.Scan(
//			&tmpOrder.OrderUID,
//			&tmpOrder.Data,
//		)
//		if err != nil {
//			return orderList, err
//		}
//		orderList = append(orderList, &tmpOrder)
//	}
//	if err = rows.Err(); err != nil {
//		return orderList, err
//	}
//	return orderList, nil
//}

func (od *order) GetOrderList(ctx context.Context, tx *sql.Tx) ([]*models.OrderDB, error) {
	query := `SELECT * FROM orders`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.New("или вот тут")
	}
	defer rows.Close()

	orderList := make([]*models.OrderDB, 0)

	for rows.Next() {
		tmpOrder := models.OrderDB{}
		err = rows.Scan(
			&tmpOrder.OrderUID,
			&tmpOrder.Data,
		)
		if err != nil {
			return orderList, errors.New("Здесь говно")
		}
		orderList = append(orderList, &tmpOrder)
	}
	if err = rows.Err(); err != nil {
		return orderList, errors.New("а не, вот тут")
	}
	return orderList, nil
}

func convertOrderForDB(or *models.Order) (*models.OrderDB, error) {
	order, err := json.Marshal(or)
	if err != nil {
		return nil, err
	}
	orderDB := &models.OrderDB{
		OrderUID: or.OrderUID,
		Data:     order,
	}
	return orderDB, nil
}

func convertOrderFromDB(orDB *models.OrderDB) (*models.Order, error) {

	var orderDbData models.Order
	err := json.Unmarshal(orDB.Data, &orderDbData)
	if err != nil {
		return nil, err
	}
	return &orderDbData, nil
}
