package queryrepo

import (
	"context"
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

const queryGetOrderByID = `
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
    WHERE o.id = $1 AND o.deleted_at IS NULL;
`

func (r *OrderPostgreQueryRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.OrderView, error) {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryGetOrderByID)
	if errStmt != nil {
		return nil, errStmt
	}
	defer stmt.Close()

	rows, err := r.Conn.QueryContext(ctx, queryGetOrderByID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query order: %w", err)
	}
	defer rows.Close()

	var order entity.OrderView
	items := make(map[uuid.UUID]entity.OrderItemView)

	for rows.Next() {
		var (
			// order
			o entity.Order
			// address
			addressID                    uuid.UUID
			street, city, state, zipCode string
			addressNote                  string
			// item
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
			// Order fields
			&o.ID, &o.UserID, &o.Status, &o.TotalPrice, &o.PaymentID,
			&o.CreatedAt, &o.UpdatedAt, &o.DeletedAt,
			// Address fields
			&addressID, &street, &city, &state, &zipCode, &addressNote,
			// Item fields
			&itemID, &productID, &productName, &productPrice, &productQuantity,
			&productImageURL, &productDescription, &productCategoryID,
			&productCategoryName, &itemNote,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		order.Address = entity.OrderAddressView{
			ID:      addressID,
			OrderID: order.ID,
			Street:  street,
			City:    city,
			State:   state,
			ZipCode: zipCode,
			Note:    addressNote,
		}

		// add item if it's not already in the map
		if _, exists := items[itemID]; !exists {
			items[itemID] = entity.OrderItemView{
				ID:                  itemID,
				OrderID:             order.ID,
				ProductID:           productID,
				ProductName:         productName,
				ProductPrice:        productPrice,
				ProductQuantity:     productQuantity,
				ProductImageUrl:     productImageURL,
				ProductDescription:  productDescription,
				ProductCategoryID:   productCategoryID,
				ProdcutCategoryName: productCategoryName,
				Note:                itemNote,
			}
		}
	}

	return &order, nil
}

const queryGetAllOrder = `SELECT * FROM order_view WHERE deleted_at IS NULL;`

func (r *OrderPostgreQueryRepo) GetAllOrders(ctx context.Context) ([]*entity.Order, error) {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryGetAllOrder)
	if errStmt != nil {
		return nil, errStmt
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]*entity.Order, 0)
	for rows.Next() {
		order := &entity.Order{}
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Status,
			&order.TotalPrice,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			continue
		}
		orders = append(orders, order)
	}

	return orders, nil
}

const queryGetOrderByUserID = `SELECT * FROM order_view WHERE user_id = $1 AND deleted_at IS NULL;`

const queryGetOrderByPaymentID = `SELECT * FROM order_view WHERE payment_id = $1 AND deleted_at IS NULL;`

const queryGetOrderByStatus = `SELECT * FROM order_view WHERE status = $1 AND deleted_at IS NULL;`
