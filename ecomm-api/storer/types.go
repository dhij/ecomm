package storer

import "time"

type Product struct {
	ID           int64      `json:"id" db:"id"`
	Name         string     `json:"name" db:"name"`
	Image        string     `json:"image" db:"image"`
	Category     string     `json:"category" db:"category"`
	Description  string     `json:"description" db:"description"`
	Rating       int64      `json:"rating" db:"rating"`
	NumReviews   int64      `json:"num_reviews" db:"num_reviews"`
	Price        float64    `json:"price" db:"price"`
	CountInStock int64      `json:"count_in_stock" db:"count_in_stock"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at" db:"updated_at"`
}
