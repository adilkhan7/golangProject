package good

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

type Good struct {
	log *log.Logger
	db  *sqlx.DB
}

func New(log *log.Logger, db *sqlx.DB) Good {
	return Good{
		log: log,
		db:  db,
	}
}

func (c Good) Create(ctx context.Context, traceID string, ng NewGood, now time.Time) (Info, error) {
	good := Info{
		ID:          uuid.New().String(),
		Name:        ng.Name,
		Price:       ng.Price,
		CategoryID:  ng.CategoryID,
		DateCreated: now.UTC(),
		DateUpdated: now.UTC(),
	}

	const q = `INSERT INTO goods (good_id, name, price, category_id, date_created, date_updated)
			   VAlUES ($1, $2, $3, $4, $5, $6)`

	c.log.Printf("%s : %s query : %s", traceID, "good.Create",
		database.Log(q, good.ID, good.Name, good.Price, good.CategoryID, good.DateCreated, good.DateUpdated),
	)

	if _, err := c.db.ExecContext(ctx, q, good.ID, good.Name, good.Price, good.CategoryID, good.DateCreated, good.DateUpdated); err != nil {
		return Info{}, errors.Wrap(err, "inserting good")
	}
	return good, nil
}

func (c Good) Update(ctx context.Context, traceID string, claims auth.Claims, goodID string, ug UpdateGood, now time.Time) error {
	good, err := c.QueryByID(ctx, traceID, claims, goodID)
	if err != nil {
		return err
	}

	if ug.Name != nil {
		good.Name = *ug.Name
	}
	if ug.Price != nil {
		good.Price = *ug.Price
	}
	if ug.CategoryID != nil {
		good.CategoryID = *ug.CategoryID
	}
	good.DateUpdated = now

	const q = `
	UPDATE
		goods
	SET 
		"name" = $2,
		"price" = $3,
		"category_id" = $4,
		"date_updated" = $5
	WHERE
		good_id = $1`

	c.log.Printf("%s: %s: %s", traceID, "good.Update",
		database.Log(q, good.ID, good.Name, good.Price, good.CategoryID, good.DateCreated, good.DateUpdated),
	)

	if _, err = c.db.ExecContext(ctx, q, good.ID, good.Name, good.Price, good.CategoryID, good.DateUpdated); err != nil {
		return errors.Wrap(err, "updating good")
	}

	return nil
}

func (c Good) Delete(ctx context.Context, traceID string, goodID string) error {
	if _, err := uuid.Parse(goodID); err != nil {
		return user.ErrInvalidID
	}
	const q = `DELETE FROM goods where good_id = $1`

	c.log.Printf("%s : %s query : %s", traceID, "good.Delete",
		database.Log(q, goodID),
	)

	if _, err := c.db.ExecContext(ctx, q, goodID); err != nil {
		return errors.Wrapf(err, "deleting good %s", goodID)
	}
	return nil
}

func (c Good) Query(ctx context.Context, traceID string, pageNumber, rowsPerPage int) ([]Info, error) {
	const q = `SELECT * FROM goods ORDER BY good_id OFFSET $1 ROWS FETCH NEXT $2 ROWS ONLY`
	offset := (pageNumber - 1) * rowsPerPage

	c.log.Printf("%s : %s query : %s", traceID, "good.Query", database.Log(q, offset, rowsPerPage))

	good := []Info{}

	if err := c.db.SelectContext(ctx, &good, q, offset, rowsPerPage); err != nil {
		return nil, errors.Wrap(err, "selecting good")
	}

	return good, nil
}

func (c Good) QueryByID(ctx context.Context, traceID string, claims auth.Claims, goodID string) (Info, error) {
	if _, err := uuid.Parse(goodID); err != nil {
		return Info{}, user.ErrInvalidID
	}

	if !claims.Authorized(auth.RoleAdmin) {
		return Info{}, user.ErrForbidden
	}

	const q = `SELECT * FROM goods WHERE good_id = $1`

	c.log.Printf("%s : %s query : %s", traceID, "good.QueryByID",
		database.Log(q, goodID),
	)

	var good Info
	if err := c.db.GetContext(ctx, &good, q, goodID); err != nil {
		if err == sql.ErrNoRows {
			return Info{}, user.ErrNotFound
		}
		return Info{}, errors.Wrapf(err, "selecting good %q", goodID)
	}

	return good, nil
}
