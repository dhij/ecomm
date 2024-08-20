package storer

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	maxAttempts = 3
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
	_, err := ms.db.NamedExecContext(ctx, "UPDATE products SET name=:name, image=:image, category=:category, description=:description, rating=:rating, num_reviews=:num_reviews, price=:price, count_in_stock=:count_in_stock, updated_at=:updated_at WHERE id=:id", p)
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
	res, err := tx.NamedExecContext(ctx, "INSERT INTO orders (payment_method, tax_price, shipping_price, total_price, user_id) VALUES (:payment_method, :tax_price, :shipping_price, :total_price, :user_id)", o)
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

func (ms *MySQLStorer) GetOrder(ctx context.Context, userID int64) (*Order, error) {
	var o Order
	err := ms.db.GetContext(ctx, &o, "SELECT * FROM orders WHERE user_id=?", userID)
	if err != nil {
		return nil, fmt.Errorf("error getting order: %w", err)
	}

	var items []OrderItem
	err = ms.db.SelectContext(ctx, &items, "SELECT * FROM order_items WHERE order_id=?", o.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting order items: %w", err)
	}
	o.Items = items

	return &o, nil
}

func (ms *MySQLStorer) GetOrderStatusByID(ctx context.Context, id int64) (*Order, error) {
	var o Order
	err := ms.db.GetContext(ctx, &o, "SELECT id, user_id, status FROM orders WHERE id=?", id)
	if err != nil {
		return nil, fmt.Errorf("error getting order: %w", err)
	}

	return &o, nil
}

func (ms *MySQLStorer) ListOrders(ctx context.Context) ([]*Order, error) {
	var orders []*Order
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

func (ms *MySQLStorer) UpdateOrderStatus(ctx context.Context, o *Order) (*Order, error) {
	_, err := ms.db.NamedExecContext(ctx, "UPDATE orders SET status=:status, updated_at=:updated_at WHERE id=:id", o)
	if err != nil {
		return nil, fmt.Errorf("error updating order status: %w", err)
	}

	return o, nil
}

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

func (ms *MySQLStorer) CreateUser(ctx context.Context, u *User) (*User, error) {
	res, err := ms.db.NamedExecContext(ctx, "INSERT INTO users (name, email, password, is_admin) VALUES (:name, :email, :password, :is_admin)", u)
	if err != nil {
		return nil, fmt.Errorf("error inserting user: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert ID: %w", err)
	}
	u.ID = id

	return u, nil
}

func (ms *MySQLStorer) GetUser(ctx context.Context, email string) (*User, error) {
	var u User
	err := ms.db.GetContext(ctx, &u, "SELECT * FROM users WHERE email=?", email)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return &u, nil
}

func (ms *MySQLStorer) ListUsers(ctx context.Context) ([]*User, error) {
	var users []*User
	err := ms.db.SelectContext(ctx, &users, "SELECT * FROM users")
	if err != nil {
		return nil, fmt.Errorf("error listing users: %w", err)
	}

	return users, nil
}

func (ms *MySQLStorer) UpdateUser(ctx context.Context, u *User) (*User, error) {
	_, err := ms.db.NamedExecContext(ctx, "UPDATE users SET name=:name, email=:email, password=:password, is_admin=:is_admin, updated_at=:updated_at WHERE id=:id", u)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return u, nil
}

func (ms *MySQLStorer) DeleteUser(ctx context.Context, id int64) error {
	_, err := ms.db.ExecContext(ctx, "DELETE FROM users WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	return nil
}

func (ms *MySQLStorer) CreateSession(ctx context.Context, s *Session) (*Session, error) {
	_, err := ms.db.NamedExecContext(ctx, "INSERT INTO sessions (id, user_email, refresh_token, is_revoked, expires_at) VALUES (:id, :user_email, :refresh_token, :is_revoked, :expires_at)", s)
	if err != nil {
		return nil, fmt.Errorf("error inserting session: %w", err)
	}

	return s, nil
}

func (ms *MySQLStorer) GetSession(ctx context.Context, id string) (*Session, error) {
	var s Session
	err := ms.db.GetContext(ctx, &s, "SELECT * FROM sessions WHERE id=?", id)
	if err != nil {
		return nil, fmt.Errorf("error getting session: %w", err)
	}

	return &s, nil
}

func (ms *MySQLStorer) RevokeSession(ctx context.Context, id string) error {
	_, err := ms.db.NamedExecContext(ctx, "UPDATE sessions SET is_revoked=1 WHERE id=:id", map[string]interface{}{"id": id})
	if err != nil {
		return fmt.Errorf("error revoking session: %w", err)
	}

	return nil
}

func (ms *MySQLStorer) DeleteSession(ctx context.Context, id string) error {
	_, err := ms.db.ExecContext(ctx, "DELETE FROM sessions WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("error deleting session: %w", err)
	}

	return nil
}

func insertNotificationState(ctx context.Context, tx *sqlx.Tx, es *NotificationState) (*NotificationState, error) {
	res, err := tx.NamedExecContext(ctx, "INSERT INTO notification_states (order_id, state, message) VALUES (:order_id, :state, :message)", es)
	if err != nil {
		return nil, fmt.Errorf("error inserting notification state: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert ID: %w", err)
	}
	es.ID = id

	return es, nil
}

func insertNotificationEvent(ctx context.Context, tx *sqlx.Tx, u *NotificationEvent) (*NotificationEvent, error) {
	res, err := tx.NamedExecContext(ctx, "INSERT INTO notification_events_queue (user_email, order_status, order_id, state_id, attempts) VALUES (:user_email, :order_status, :order_id, :state_id, :attempts)", u)
	if err != nil {
		return nil, fmt.Errorf("error inserting notification event: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert ID: %w", err)
	}
	u.ID = id

	return u, nil
}

func (ms *MySQLStorer) EnqueueNotificationEvent(ctx context.Context, ne *NotificationEvent) (*NotificationEvent, error) {
	var ev *NotificationEvent
	err := ms.execTx(ctx, func(tx *sqlx.Tx) error {
		ns, err := insertNotificationState(ctx, tx, &NotificationState{
			OrderID: ne.OrderID,
			State:   NotSent,
			Message: "",
		})
		if err != nil {
			return fmt.Errorf("error inserting notification state: %w", err)
		}
		ne.StateID = ns.ID

		ev, err = insertNotificationEvent(ctx, tx, ne)
		if err != nil {
			return fmt.Errorf("error inserting notification event: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error enqueuing notification event: %w", err)
	}

	return ev, nil
}

func (ms *MySQLStorer) ListNotificationEvents(ctx context.Context) ([]*NotificationEvent, error) {
	var events []*NotificationEvent

	q := fmt.Sprintf("SELECT * FROM notification_events_queue WHERE attempts < %d ORDER BY created_at", maxAttempts)
	err := ms.db.SelectContext(ctx, &events, q)
	if err != nil {
		return nil, fmt.Errorf("error listing notification events: %w", err)
	}

	return events, nil
}

func getNotificationEventAttempts(ctx context.Context, tx *sqlx.Tx, id int64) (*NotificationEvent, error) {
	var u NotificationEvent
	err := tx.GetContext(ctx, &u, "SELECT id, attempts FROM notification_events_queue WHERE id=?", id)
	if err != nil {
		return nil, fmt.Errorf("error getting notification event: %w", err)
	}

	return &u, nil
}

func updateNotificationEventAttempts(ctx context.Context, tx *sqlx.Tx, u *NotificationEvent) (*NotificationEvent, error) {
	_, err := tx.NamedExecContext(ctx, "UPDATE notification_events_queue SET attempts=:attempts, updated_at=:updated_at WHERE id=:id", u)
	if err != nil {
		return nil, fmt.Errorf("error updating notification event: %w", err)
	}

	return u, nil
}

func deleteNotificationEvent(ctx context.Context, tx *sqlx.Tx, id int64) error {
	_, err := tx.ExecContext(ctx, "DELETE FROM notification_events_queue WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("error deleting notification event: %w", err)
	}

	return nil
}

func updateNotificationState(ctx context.Context, tx *sqlx.Tx, es *NotificationState) error {
	q := "UPDATE notification_states SET state=:state, message=:message WHERE id=:id"
	if es.State == Sent {
		t := time.Now()
		es.CompletedAt = &t
		q = "UPDATE notification_states SET state=:state, message=:message, completed_at=:completed_at WHERE id=:id"
	}

	_, err := tx.NamedExecContext(ctx, q, es)
	if err != nil {
		return fmt.Errorf("error updating notification state: %w", err)
	}

	return nil
}

func (ms *MySQLStorer) UpdateNotificationEvent(ctx context.Context, ev *NotificationEvent, es *NotificationState, responseType NotificationResponseType) (bool, error) {
	succeeded := false
	err := ms.execTx(ctx, func(tx *sqlx.Tx) error {
		switch responseType {
		case NotificationSucess:
			err := updateNotificationState(ctx, tx, &NotificationState{
				ID:      ev.StateID,
				State:   Sent,
				Message: es.Message,
			})
			if err != nil {
				return fmt.Errorf("error updating notification state: %w", err)
			}

			err = deleteNotificationEvent(ctx, tx, ev.ID)
			if err != nil {
				return fmt.Errorf("error deleting notification event: %w", err)
			}
			succeeded = true
		case NotificationFailure:
			u, err := getNotificationEventAttempts(ctx, tx, ev.ID)
			if err != nil {
				return fmt.Errorf("error getting notification event: %w", err)
			}

			if u.Attempts+1 < maxAttempts {
				t := time.Now()
				u.UpdatedAt = &t
				u.Attempts += 1

				_, err = updateNotificationEventAttempts(ctx, tx, u)
				if err != nil {
					return fmt.Errorf("error updating notification event: %w", err)
				}
			} else {
				err = updateNotificationState(ctx, tx, &NotificationState{
					ID:      ev.StateID,
					State:   Failed,
					Message: es.Message,
				})
				if err != nil {
					return fmt.Errorf("error updating notification state: %w", err)
				}

				err = deleteNotificationEvent(ctx, tx, u.ID)
				if err != nil {
					return fmt.Errorf("error deleting notification event: %w", err)
				}
			}

		default:
			return fmt.Errorf("invalid notification response type: %v", responseType)
		}

		return nil
	})

	return succeeded, err
}
