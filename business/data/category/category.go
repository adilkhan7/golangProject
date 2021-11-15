package category

import (
	"context"
	"database/sql"
	"github.com/adilkhan7/golangSoftProject/business/auth"
	"github.com/adilkhan7/golangSoftProject/business/data/user"
	"github.com/adilkhan7/golangSoftProject/foundation/database"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"log"
	"time"
)

type Category struct {
	log *log.Logger
	db  *sqlx.DB
}

func New(log *log.Logger, db *sqlx.DB) Category {
	return Category{
		log: log,
		db:  db,
	}
}

func (c Category) Create(ctx context.Context, traceID string, nc NewCategory, now time.Time) (Info, error) {
	category := Info{
		ID:            uuid.New().String(),
		Name:          nc.Name,
		SupermarketID: nc.SupermarketID,
		DateCreated:   now.UTC(),
		DateUpdated:   now.UTC(),
	}

	const q = `INSERT INTO categories (category_id, name, supermarket_id, date_created, date_updated)
			   VAlUES ($1, $2, $3, $4, $5)`

	c.log.Printf("%s : %s query : %s", traceID, "category.Create",
		database.Log(q, category.ID, category.Name, category.SupermarketID, category.DateCreated, category.DateUpdated),
	)

	if _, err := c.db.ExecContext(ctx, q, category.ID, category.Name, category.SupermarketID, category.DateCreated, category.DateUpdated); err != nil {
		return Info{}, errors.Wrap(err, "inserting category")
	}
	return category, nil
}

func (c Category) Update(ctx context.Context, traceID string, claims auth.Claims, categoryID string, uc UpdateCategory, now time.Time) error {
	category, err := c.QueryByID(ctx, traceID, claims, categoryID)
	if err != nil {
		return err
	}

	if uc.Name != nil {
		category.Name = *uc.Name
	}
	if uc.SupermarketID != nil {
		category.SupermarketID = *uc.SupermarketID
	}
	category.DateUpdated = now

	const q = `
	UPDATE
		categories
	SET 
		"name" = $2,
		"supermarket_id" = $3,
		"date_updated" = $4
	WHERE
		category_id = $1`

	c.log.Printf("%s: %s: %s", traceID, "category.Update",
		database.Log(q, category.ID, category.Name, category.SupermarketID, category.DateCreated, category.DateUpdated),
	)

	if _, err = c.db.ExecContext(ctx, q, category.ID, category.Name, category.SupermarketID, category.DateUpdated); err != nil {
		return errors.Wrap(err, "updating category")
	}

	return nil
}

func (c Category) Delete(ctx context.Context, traceID string, categoryID string) error {
	if _, err := uuid.Parse(categoryID); err != nil {
		return user.ErrInvalidID
	}
	const q = `DELETE FROM categories where category_id = $1`

	c.log.Printf("%s : %s query : %s", traceID, "category.Delete",
		database.Log(q, categoryID),
	)

	if _, err := c.db.ExecContext(ctx, q, categoryID); err != nil {
		return errors.Wrapf(err, "deleting category %s", categoryID)
	}
	return nil
}

func (c Category) Query(ctx context.Context, traceID string, pageNumber, rowsPerPage int) ([]Info, error) {
	const q = `SELECT * FROM categories ORDER BY category_id OFFSET $1 ROWS FETCH NEXT $2 ROWS ONLY`
	offset := (pageNumber - 1) * rowsPerPage

	c.log.Printf("%s : %s query : %s", traceID, "category.Query", database.Log(q, offset, rowsPerPage))

	category := []Info{}

	if err := c.db.SelectContext(ctx, &category, q, offset, rowsPerPage); err != nil {
		return nil, errors.Wrap(err, "selecting category")
	}

	return category, nil
}

func (c Category) QueryByID(ctx context.Context, traceID string, claims auth.Claims, categoryID string) (Info, error) {
	if _, err := uuid.Parse(categoryID); err != nil {
		return Info{}, user.ErrInvalidID
	}

	if !claims.Authorized(auth.RoleAdmin) {
		return Info{}, user.ErrForbidden
	}

	const q = `SELECT * FROM categories WHERE category_id = $1`

	c.log.Printf("%s : %s query : %s", traceID, "category.QueryByID",
		database.Log(q, categoryID),
	)

	var category Info
	if err := c.db.GetContext(ctx, &category, q, categoryID); err != nil {
		if err == sql.ErrNoRows {
			return Info{}, user.ErrNotFound
		}
		return Info{}, errors.Wrapf(err, "selecting category %q", categoryID)
	}

	return category, nil
}
