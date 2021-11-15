package category

import "time"

type Info struct {
	ID            string    `db:"category_id" json:"category_id"`
	Name          string    `db:"name" json:"name"`
	SupermarketID string    `db:"supermarket_id" json:"supermarket_id"`
	DateCreated   time.Time `db:"date_created" json:"date_created"`
	DateUpdated   time.Time `db:"date_updated" json:"date_updated"`
}

type NewCategory struct {
	Name          string `json:"name" validate:"required"`
	SupermarketID string `json:"supermarket_id" validate:"required"`
}

type UpdateCategory struct {
	Name          *string `json:"name"`
	SupermarketID *string `json:"supermarket_id"`
}
