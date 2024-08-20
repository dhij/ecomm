package storer

import "time"

type Product struct {
	ID           int64      `db:"id"`
	Name         string     `db:"name"`
	Image        string     `db:"image"`
	Category     string     `db:"category"`
	Description  string     `db:"description"`
	Rating       int64      `db:"rating"`
	NumReviews   int64      `db:"num_reviews"`
	Price        float32    `db:"price"`
	CountInStock int64      `db:"count_in_stock"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at"`
}

type OrderStatus string

const (
	Pending   OrderStatus = "pending"
	Shipped   OrderStatus = "shipped"
	Delivered OrderStatus = "delivered"
)

type Order struct {
	ID            int64       `db:"id"`
	PaymentMethod string      `db:"payment_method"`
	TaxPrice      float32     `db:"tax_price"`
	ShippingPrice float32     `db:"shipping_price"`
	TotalPrice    float32     `db:"total_price"`
	UserID        int64       `db:"user_id"`
	Status        OrderStatus `db:"status"`
	CreatedAt     time.Time   `db:"created_at"`
	UpdatedAt     *time.Time  `db:"updated_at"`
	Items         []OrderItem
}

type OrderItem struct {
	ID        int64   `db:"id"`
	Name      string  `db:"name"`
	Quantity  int64   `db:"quantity"`
	Image     string  `db:"image"`
	Price     float32 `db:"price"`
	ProductID int64   `db:"product_id"`
	OrderID   int64   `db:"order_id"`
}

type User struct {
	ID        int64      `db:"id"`
	Name      string     `db:"name"`
	Email     string     `db:"email"`
	Password  string     `db:"password"`
	IsAdmin   bool       `db:"is_admin"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

type Session struct {
	ID           string    `db:"id"`
	UserEmail    string    `db:"user_email"`
	RefreshToken string    `db:"refresh_token"`
	IsRevoked    bool      `db:"is_revoked"`
	CreatedAt    time.Time `db:"created_at"`
	ExpiresAt    time.Time `db:"expires_at"`
}

type NotificationEventState string

const (
	NotSent NotificationEventState = "not sent"
	Sent    NotificationEventState = "sent"
	Failed  NotificationEventState = "failed"
)

type NotificationResponseType string

const (
	NotificationSucess  NotificationResponseType = "success"
	NotificationFailure NotificationResponseType = "failure"
)

type NotificationState struct {
	ID          int64                  `db:"id"`
	OrderID     int64                  `db:"order_id"`
	State       NotificationEventState `db:"state"`
	Message     string                 `db:"message"`
	RequestedAt time.Time              `db:"requested_at"`
	CompletedAt *time.Time             `db:"completed_at"`
}

type NotificationEvent struct {
	ID          int64       `db:"id"`
	UserEmail   string      `db:"user_email"`
	OrderStatus OrderStatus `db:"order_status"`
	OrderID     int64       `db:"order_id"`
	StateID     int64       `db:"state_id"`
	Attempts    int64       `db:"attempts"`
	CreatedAt   time.Time   `db:"created_at"`
	UpdatedAt   *time.Time  `db:"updated_at"`
}
