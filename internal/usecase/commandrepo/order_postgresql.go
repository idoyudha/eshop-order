package commandrepo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-order/internal/entity"
	"github.com/idoyudha/eshop-order/pkg/postgresql/postgrecommand"
)

type OrderPostgreCommandRepo struct {
	*postgrecommand.PostgresCommand
}

func NewOrderPostgreCommandRepo(conn *postgrecommand.PostgresCommand) *OrderPostgreCommandRepo {
	return &OrderPostgreCommandRepo{
		PostgresCommand: conn,
	}
}

const (
	queryInsertOrder        = `INSERT INTO orders (id, user_id, status, total_price, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6);`
	queryInsertOrderItems   = `INSERT INTO order_items (id, order_id, product_id, product_quantity, shipping_cost, note, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`
	queryInsertOrderAddress = `INSERT INTO order_addresses (id, order_id, street, city, state, zip_code, note, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`
)

func (r *OrderPostgreCommandRepo) Insert(ctx context.Context, order *entity.Order) error {
	// begin transaction
	tx, err := r.Conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// insert order
	_, err = tx.ExecContext(ctx, queryInsertOrder,
		order.ID, order.UserID, order.Status, order.TotalPrice, order.CreatedAt, order.UpdatedAt)
	if err != nil {
		return err
	}

	// insert order items
	for _, item := range order.Items {
		_, err = tx.ExecContext(ctx, queryInsertOrderItems,
			item.ID, order.ID, item.ProductID, item.ProductQuantity, item.ShippingCost, item.Note, item.CreatedAt, item.UpdatedAt)
		if err != nil {
			return err
		}
	}

	// insert order address
	_, err = tx.ExecContext(ctx, queryInsertOrderAddress,
		order.Address.ID, order.ID, order.Address.Street, order.Address.City, order.Address.State, order.Address.ZipCode, order.Address.Note, order.Address.CreatedAt, order.Address.UpdatedAt)
	if err != nil {
		return err
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

const queryUpdateStatusOrder = `UPDATE orders SET status = $1, updated_at = $2 WHERE id = $3;`

func (r *OrderPostgreCommandRepo) UpdateStatus(ctx context.Context, order *entity.Order) error {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryUpdateStatusOrder)
	if errStmt != nil {
		return errStmt
	}
	defer stmt.Close()

	_, updateErr := stmt.ExecContext(ctx, order.Status, order.UpdatedAt, order.ID)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

const queryUpdatePaymentIDOrder = `UPDATE orders SET status = $1, payment_id = $2, updated_at = $3 WHERE id = $4;`

func (r *OrderPostgreCommandRepo) UpdatePaymentID(ctx context.Context, order *entity.Order) error {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryUpdatePaymentIDOrder)
	if errStmt != nil {
		return errStmt
	}
	defer stmt.Close()

	_, updateErr := stmt.ExecContext(ctx, order.Status, order.PaymentID, order.UpdatedAt, order.ID)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

const queryGetOrderByID = `
	SELECT 
		o.id,
		o.user_id,
		oi.product_id as item_product_id,
		oi.product_quantity as item_product_quantity
	FROM orders o
	LEFT JOIN order_items oi ON o.id = oi.order_id
	WHERE o.id = $1;
`

func (r *OrderPostgreCommandRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryGetOrderByID)
	if errStmt != nil {
		return nil, errStmt
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var order entity.Order
	for rows.Next() {
		var item entity.OrderItem
		if err := rows.Scan(&order.ID, &order.UserID, &item.ProductID, &item.ProductQuantity); err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &order, nil
}
