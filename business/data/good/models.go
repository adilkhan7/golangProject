package good

import "time"

type Info struct {
	ID          string    `db:"good_id" json:"good_id"`
	Name        string    `db:"name" json:"name"`
	Price       int       `db:"price" json:"price"`
	CategoryID  string    `db:"category_id" json:"category_id"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
	DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}

type NewGood struct {
	Name       string `json:"name" validate:"required"`
	Price      int    `json:"price" validate:"required"`
	CategoryID string `json:"category_id" validate:"required"`
}

type UpdateGood struct {
	Name       *string `json:"name"`
	Price      *int    `json:"price"`
	CategoryID *string `json:"category_id"`
}
