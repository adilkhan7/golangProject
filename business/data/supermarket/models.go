package supermarket

import "time"

type Info struct {
	ID          string    `db:"supermarket_id" json:"supermarket_id"`
	Name        string    `db:"name" json:"name"`
	Address     string    `db:"address" json:"address"`
	UserID      string    `db:"user_id" json:"user_id"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
	DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}

type NewSuperMarket struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
	UserID  string `json:"user_id" validate:"required"`
}

type UpdateSupermarket struct {
	Name    *string `json:"name"`
	Address *string `json:"address"`
	UserID  *string `json:"user_id"`
}
