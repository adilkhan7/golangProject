package supermarket

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

type Supermarket struct {
	log *log.Logger
	db  *sqlx.DB
}

func New(log *log.Logger, db *sqlx.DB) Supermarket {
	return Supermarket{
		log: log,
		db:  db,
	}
}

func (s Supermarket) Create(ctx context.Context, traceID string, ns NewSuperMarket, now time.Time) (Info, error) {
	supermarket := Info{
		ID:          uuid.New().String(),
		Name:        ns.Name,
		Address:     ns.Address,
		UserID:      ns.UserID,
		DateCreated: now.UTC(),
		DateUpdated: now.UTC(),
	}

	const q = `INSERT INTO supermarkets (supermarket_id, name, address, user_id, date_created, date_updated)
			   VAlUES ($1, $2, $3, $4, $5, $6)`

	s.log.Printf("%s : %s query : %s", traceID, "supermarket.Create",
		database.Log(q, supermarket.ID, supermarket.Name, supermarket.Address, supermarket.UserID, supermarket.DateCreated, supermarket.DateUpdated),
	)

	if _, err := s.db.ExecContext(ctx, q, supermarket.ID, supermarket.Name, supermarket.Address, supermarket.UserID, supermarket.DateCreated, supermarket.DateUpdated); err != nil {
		return Info{}, errors.Wrap(err, "inserting supermarket")
	}
	return supermarket, nil
}

func (s Supermarket) Update(ctx context.Context, traceID string, claims auth.Claims, supermarketID string, us UpdateSupermarket, now time.Time) error {
	supermarket, err := s.QueryByID(ctx, traceID, claims, supermarketID)
	if err != nil {
		return err
	}

	if us.Name != nil {
		supermarket.Name = *us.Name
	}
	if us.Address != nil {
		supermarket.Address = *us.Address
	}
	supermarket.DateUpdated = now

	const q = `
	UPDATE
		supermarkets
	SET 
		"name" = $2,
		"address" = $3,
		"date_updated" = $4
	WHERE
		supermarket_id = $1`

	s.log.Printf("%s: %s: %s", traceID, "supermarket.Update",
		database.Log(q, supermarket.ID, supermarket.Name, supermarket.Address, supermarket.UserID, supermarket.DateCreated, supermarket.DateUpdated),
	)

	if _, err = s.db.ExecContext(ctx, q, supermarket.ID, supermarket.Name, supermarket.Address, supermarket.DateUpdated); err != nil {
		return errors.Wrap(err, "updating supermarket")
	}

	return nil
}

func (s Supermarket) Delete(ctx context.Context, traceID string, supermarketID string) error {
	if _, err := uuid.Parse(supermarketID); err != nil {
		return user.ErrInvalidID
	}
	const q = `DELETE FROM supermarkets where supermarket_id = $1`

	s.log.Printf("%s : %s query : %s", traceID, "supermarket.Delete",
		database.Log(q, supermarketID),
	)

	if _, err := s.db.ExecContext(ctx, q, supermarketID); err != nil {
		return errors.Wrapf(err, "deleting supermarket %s", supermarketID)
	}
	return nil
}

func (s Supermarket) Query(ctx context.Context, traceID string, pageNumber, rowsPerPage int) ([]Info, error) {
	const q = `SELECT *FROM supermarkets ORDER BY supermarket_id OFFSET $1 ROWS FETCH NEXT $2 ROWS ONLY`
	offset := (pageNumber - 1) * rowsPerPage

	s.log.Printf("%s : %s query : %s", traceID, "supermarket.Query", database.Log(q, offset, rowsPerPage))

	supermarket := []Info{}

	if err := s.db.SelectContext(ctx, &supermarket, q, offset, rowsPerPage); err != nil {
		return nil, errors.Wrap(err, "selecting supermarket")
	}

	return supermarket, nil
}

func (s Supermarket) QueryByID(ctx context.Context, traceID string, claims auth.Claims, supermarketID string) (Info, error) {
	if _, err := uuid.Parse(supermarketID); err != nil {
		return Info{}, user.ErrInvalidID
	}

	if !claims.Authorized(auth.RoleAdmin) {
		return Info{}, user.ErrForbidden
	}

	const q = `SELECT * FROM supermarkets WHERE supermarket_id = $1`

	s.log.Printf("%s : %s query : %s", traceID, "supermarket.QueryByID",
		database.Log(q, supermarketID),
	)

	var supermarket Info
	if err := s.db.GetContext(ctx, &supermarket, q, supermarketID); err != nil {
		if err == sql.ErrNoRows {
			return Info{}, user.ErrNotFound
		}
		return Info{}, errors.Wrapf(err, "selecting supermarket %q", supermarketID)
	}

	return supermarket, nil
}
