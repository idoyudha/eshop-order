package queryrepo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-order/internal/entity"
	"github.com/idoyudha/eshop-order/pkg/postgresql/postgrequery"
)

type OrderPostgreQueryRepo struct {
	*postgrequery.PostgresQuery
}

func NewOrderPostgreCommandRepo(conn *postgrequery.PostgresQuery) *OrderPostgreQueryRepo {
	return &OrderPostgreQueryRepo{
		PostgresQuery: conn,
	}
}

const (
	queryInsertOrdersView       = `INSERT INTO orders_view (id, order_id, user_id, status, total_price, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7);`
	queryInserrOrderItemsView   = `INSERT INTO order_items_view (id, order_view_id, product_id, product_name, product_price, product_quantity, product_image_url, product_description, product_category_id, product_category_name, note, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);`
	queryInsertOrderAddressView = `INSERT INTO order_addresses_view (id, order_view_id, street, city, state, zip_code, note, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`
)

func (r *OrderPostgreQueryRepo) Insert(ctx context.Context, order *entity.OrderView) error {
	tx, err := r.Conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// order
	_, err = tx.ExecContext(ctx, queryInsertOrdersView,
		order.ID, order.OrderID, order.UserID, order.Status, order.TotalPrice,
		order.CreatedAt, order.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert order view: %w", err)
	}

	// order items
	for _, item := range order.Items {
		_, err = tx.ExecContext(ctx, queryInserrOrderItemsView,
			item.ID, order.ID,
			item.ProductID, item.ProductName, item.ProductPrice,
			item.ProductQuantity, item.ProductImageURL, item.ProductDescription,
			item.ProductCategoryID, item.ProductCategoryName, item.Note,
			item.CreatedAt, item.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to insert order item view: %w", err)
		}
	}

	// order address
	_, err = tx.ExecContext(ctx, queryInsertOrderAddressView,
		order.Address.ID, order.ID, order.Address.Street,
		order.Address.City, order.Address.State, order.Address.ZipCode,
		order.Address.Note, order.Address.CreatedAt, order.Address.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert order address view: %w", err)
	}

	return tx.Commit()
}

const baseQueryOrder = `
	SELECT 
        o.*,
        oa.id as address_id,
        oa.street as address_street,
        oa.city as address_city,
        oa.state as address_state,
        oa.zip_code as address_zip_code,
        oa.note as address_note,
        oi.id as item_id,
        oi.product_id,
        oi.product_name,
        oi.product_price,
        oi.product_quantity,
        oi.product_image_url,
        oi.product_description,
        oi.product_category_id,
        oi.product_category_name,
        oi.note as item_note
    FROM orders_view o
    LEFT JOIN order_addresses_view oa ON o.id = oa.order_id
    LEFT JOIN order_items_view oi ON o.id = oi.order_id
    o.deleted_at IS NULL
`

const queryGetOrderByID = baseQueryOrder + ` AND o.id = $1;`

func (r *OrderPostgreQueryRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.OrderView, error) {
	return r.scanSingleOrder(ctx, queryGetOrderByID, id)
}

const queryGetAllOrder = baseQueryOrder + ` ORDER BY o.created_at DESC;`

func (r *OrderPostgreQueryRepo) GetAll(ctx context.Context) ([]*entity.OrderView, error) {
	return r.scanMultipleOrders(ctx, queryGetAllOrder)
}

const queryGetOrderByUserID = baseQueryOrder + ` AND o.user_id = $1 ORDER BY o.created_at DESC;`

func (r *OrderPostgreQueryRepo) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.OrderView, error) {
	return r.scanMultipleOrders(ctx, queryGetOrderByUserID, userID)
}

const queryGetOrderByPaymentID = baseQueryOrder + ` AND o.payment_id = $1;`

func (r *OrderPostgreQueryRepo) GetByPaymentID(ctx context.Context, paymentID uuid.UUID) (*entity.OrderView, error) {
	return r.scanSingleOrder(ctx, queryGetOrderByPaymentID, paymentID)
}

const queryGetOrderByStatus = baseQueryOrder + ` AND o.status = $1 ORDER BY o.created_at DESC;`

func (r *OrderPostgreQueryRepo) GetByStatus(ctx context.Context, status string) ([]*entity.OrderView, error) {
	return r.scanMultipleOrders(ctx, queryGetOrderByStatus, status)
}

// helper to scan single order
func (r *OrderPostgreQueryRepo) scanSingleOrder(ctx context.Context, query string, args ...interface{}) (*entity.OrderView, error) {
	rows, err := r.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query order: %w", err)
	}
	defer rows.Close()

	var order *entity.OrderView
	items := make(map[uuid.UUID]entity.OrderItemView)

	for rows.Next() {
		var (
			// order fields
			o entity.OrderView
			// address fields
			addressID                                 uuid.UUID
			street, city, state, zipCode, addressNote string
			// item fields
			itemID, productID                   uuid.UUID
			productName                         string
			productPrice                        float64
			productQuantity                     int64
			productImageURL, productDescription string
			productCategoryID                   uuid.UUID
			productCategoryName                 string
			itemNote                            string
		)

		err := rows.Scan(
			&o.ID, &o.UserID, &o.Status, &o.TotalPrice, &o.PaymentID,
			&o.PaymentStatus, &o.PaymentImageURL, &o.PaymentAdminNote,
			&o.CreatedAt, &o.UpdatedAt, &o.DeletedAt,
			&addressID, &street, &city, &state, &zipCode, &addressNote,
			&itemID, &productID, &productName, &productPrice, &productQuantity,
			&productImageURL, &productDescription, &productCategoryID,
			&productCategoryName, &itemNote,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		if order == nil {
			order = &o
			order.Address = entity.OrderAddressView{
				ID:          addressID,
				OrderViewID: order.ID,
				Street:      street,
				City:        city,
				State:       state,
				ZipCode:     zipCode,
				Note:        addressNote,
			}
		}

		items[itemID] = entity.OrderItemView{
			ID:                  itemID,
			OrderViewID:         order.ID,
			ProductID:           productID,
			ProductName:         productName,
			ProductPrice:        productPrice,
			ProductQuantity:     productQuantity,
			ProductImageURL:     productImageURL,
			ProductDescription:  productDescription,
			ProductCategoryID:   productCategoryID,
			ProductCategoryName: productCategoryName,
			Note:                itemNote,
		}

	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating order rows: %w", err)
	}

	if order == nil {
		return nil, sql.ErrNoRows
	}

	// Convert items map to slice
	order.Items = make([]entity.OrderItemView, 0, len(items))
	for _, item := range items {
		order.Items = append(order.Items, item)
	}

	return order, nil
}

// helper to scan multiple orders
func (r *OrderPostgreQueryRepo) scanMultipleOrders(ctx context.Context, query string, args ...interface{}) ([]*entity.OrderView, error) {
	rows, err := r.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}
	defer rows.Close()

	ordersMap := make(map[uuid.UUID]*entity.OrderView)
	itemsMap := make(map[uuid.UUID]map[uuid.UUID]entity.OrderItemView)

	for rows.Next() {
		var (
			// order fields
			o entity.OrderView
			// address fields
			addressID                                 uuid.UUID
			street, city, state, zipCode, addressNote string
			// item fields
			itemID, productID                   uuid.UUID
			productName                         string
			productPrice                        float64
			productQuantity                     int64
			productImageURL, productDescription string
			productCategoryID                   uuid.UUID
			productCategoryName, itemNote       string
		)

		err := rows.Scan(
			&o.ID, &o.UserID, &o.Status, &o.TotalPrice, &o.PaymentID,
			&o.PaymentStatus, &o.PaymentImageURL, &o.PaymentAdminNote,
			&o.CreatedAt, &o.UpdatedAt, &o.DeletedAt,
			&addressID, &street, &city, &state, &zipCode, &addressNote,
			&itemID, &productID, &productName, &productPrice, &productQuantity,
			&productImageURL, &productDescription, &productCategoryID,
			&productCategoryName, &itemNote,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		// if this is the first time we're seeing this order
		if _, exists := ordersMap[o.ID]; !exists {
			ordersMap[o.ID] = &o
			ordersMap[o.ID].Address = entity.OrderAddressView{
				ID:          addressID,
				OrderViewID: o.ID,
				Street:      street,
				City:        city,
				State:       state,
				ZipCode:     zipCode,
				Note:        addressNote,
			}
			itemsMap[o.ID] = make(map[uuid.UUID]entity.OrderItemView)
		}

		// add item if it exists and isn't already added
		if _, exists := itemsMap[o.ID][itemID]; !exists {
			itemsMap[o.ID][itemID] = entity.OrderItemView{
				ID:                  itemID,
				OrderViewID:         o.ID,
				ProductID:           productID,
				ProductName:         productName,
				ProductPrice:        productPrice,
				ProductQuantity:     productQuantity,
				ProductImageURL:     productImageURL,
				ProductDescription:  productDescription,
				ProductCategoryID:   productCategoryID,
				ProductCategoryName: productCategoryName,
				Note:                itemNote,
			}
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating order rows: %w", err)
	}

	// Convert maps to slices
	orders := make([]*entity.OrderView, 0, len(ordersMap))
	for orderID, order := range ordersMap {
		items := make([]entity.OrderItemView, 0, len(itemsMap[orderID]))
		for _, item := range itemsMap[orderID] {
			items = append(items, item)
		}
		order.Items = items
		orders = append(orders, order)
	}

	return orders, nil
}
