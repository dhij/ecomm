package storer

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type MySQLStorer struct {
	db *sqlx.DB
}

func NewMySQLStorer(db *sqlx.DB) *MySQLStorer {
	return &MySQLStorer{
		db: db,
	}
}

func (ms *MySQLStorer) CreateProduct(ctx context.Context, p *Product) (*Product, error) {
	res, err := ms.db.NamedExecContext(ctx, "INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (:name, :image, :category, :description, :rating, :num_reviews, :price, :count_in_stock)", p)
	if err != nil {
		return nil, fmt.Errorf("error inserting product: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert ID: %w", err)
	}
	p.ID = id

	return p, nil
}

func (ms *MySQLStorer) GetProduct(ctx context.Context, id int64) (*Product, error) {
	var p Product
	err := ms.db.GetContext(ctx, &p, "SELECT * FROM products WHERE id=?", id)
	if err != nil {
		return nil, fmt.Errorf("error getting product: %w", err)
	}

	return &p, nil
}

func (ms *MySQLStorer) ListProducts(ctx context.Context) ([]*Product, error) {
	var products []*Product
	err := ms.db.SelectContext(ctx, &products, "SELECT * FROM products")
	if err != nil {
		return nil, fmt.Errorf("error listing products: %w", err)
	}

	return products, nil
}

func (ms *MySQLStorer) UpdateProduct(ctx context.Context, p *Product) (*Product, error) {
	_, err := ms.db.NamedExecContext(ctx, "UPDATE products SET name=:name, image=:image, category=:category, description=:description, rating=:rating, num_reviews=:num_reviews, price=:price, count_in_stock=:count_in_stock WHERE id=:id", p)
	if err != nil {
		return nil, fmt.Errorf("error updating product: %w", err)
	}

	return p, nil
}

func (ms *MySQLStorer) DeleteProduct(ctx context.Context, id int64) error {
	_, err := ms.db.ExecContext(ctx, "DELETE FROM products WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("error deleting product: %w", err)
	}

	return nil
}

func (ms *MySQLStorer) CreateOrder(ctx context.Context, o *Order) (*Order, error) {
	err := ms.execTx(ctx, func(tx *sqlx.Tx) error {
		// insert into orders
		order, err := createOrder(ctx, tx, o)
		if err != nil {
			return fmt.Errorf("error creating order: %w", err)
		}

		for _, oi := range o.Items {
			oi.OrderID = order.ID
			// insert into order_items
			err = createOrderItem(ctx, tx, oi)
			if err != nil {
				return fmt.Errorf("error creating order item: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error creating order: %w", err)
	}

	return o, nil
}

func createOrder(ctx context.Context, tx *sqlx.Tx, o *Order) (*Order, error) {
	res, err := tx.NamedExecContext(ctx, "INSERT INTO orders (payment_method, tax_price, shipping_price, total_price) VALUES (:payment_method, :tax_price, :shipping_price, :total_price)", o)
	if err != nil {
		return nil, fmt.Errorf("error inserting order: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert ID: %w", err)
	}
	o.ID = id

	return o, nil
}

func createOrderItem(ctx context.Context, tx *sqlx.Tx, oi OrderItem) error {
	res, err := tx.NamedExecContext(ctx, "INSERT INTO order_items (name, quantity, image, price, product_id, order_id) VALUES (:name, :quantity, :image, :price, :product_id, :order_id)", oi)
	if err != nil {
		return fmt.Errorf("error inserting order item: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert ID: %w", err)
	}
	oi.ID = id

	return nil
}

func (ms *MySQLStorer) GetOrder(ctx context.Context, id int64) (*Order, error) {
	var o Order
	err := ms.db.GetContext(ctx, &o, "SELECT * FROM orders WHERE id=?", id)
	if err != nil {
		return nil, fmt.Errorf("error getting order: %w", err)
	}

	var items []OrderItem
	err = ms.db.SelectContext(ctx, &items, "SELECT * FROM order_items WHERE order_id=?", id)
	if err != nil {
		return nil, fmt.Errorf("error getting order items: %w", err)
	}
	o.Items = items

	return &o, nil
}

func (ms *MySQLStorer) ListOrders(ctx context.Context) ([]Order, error) {
	var orders []Order
	err := ms.db.SelectContext(ctx, &orders, "SELECT * FROM orders")
	if err != nil {
		return nil, fmt.Errorf("error listing orders: %w", err)
	}

	for i := range orders {
		var items []OrderItem
		err = ms.db.SelectContext(ctx, &items, "SELECT * FROM order_items WHERE order_id=?", orders[i].ID)
		if err != nil {
			return nil, fmt.Errorf("error getting order items: %w", err)
		}
		orders[i].Items = items
	}

	return orders, nil
}

// UpdateOrderStatus

func (ms *MySQLStorer) DeleteOrder(ctx context.Context, id int64) error {
	err := ms.execTx(ctx, func(tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM order_items WHERE order_id=?", id)
		if err != nil {
			return fmt.Errorf("error deleting order items: %w", err)
		}

		_, err = tx.ExecContext(ctx, "DELETE FROM orders WHERE id=?", id)
		if err != nil {
			return fmt.Errorf("error deleting order: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("error deleting order: %w", err)
	}

	return nil
}

func (ms *MySQLStorer) execTx(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := ms.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	err = fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("error rolling back transaction: %w", rbErr)
		}
		return fmt.Errorf("error in transaction: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}
