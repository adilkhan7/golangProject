package handlers

import (
	"context"
	"fmt"
	"github.com/adilkhan7/golangSoftProject/business/auth"
	"github.com/adilkhan7/golangSoftProject/business/data/supermarket"
	"github.com/adilkhan7/golangSoftProject/business/data/user"
	"github.com/adilkhan7/golangSoftProject/foundation/web"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

type supermarketGroup struct {
	supermarket supermarket.Supermarket
	auth        *auth.Auth
}

func (sg supermarketGroup) query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	params := web.Params(r)
	pageNumber, err := strconv.Atoi(params["page"])
	if err != nil {
		return web.NewRequestError(fmt.Errorf("invalid page format: %s", params["page"]), http.StatusBadRequest)
	}

	rowsPerPage, err := strconv.Atoi(params["rows"])
	if err != nil {
		return web.NewRequestError(fmt.Errorf("invalid rows format: %s", params["rows"]), http.StatusBadRequest)
	}

	smarket, err := sg.supermarket.Query(ctx, v.TraceID, pageNumber, rowsPerPage)
	if err != nil {
		return errors.Wrap(err, "unable to query for smarket")
	}

	return web.Respond(ctx, w, smarket, http.StatusOK)
}

func (sg supermarketGroup) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	params := web.Params(r)
	smarket, err := sg.supermarket.QueryByID(ctx, v.TraceID, claims, params["id"])

	if err != nil {
		if err != nil {
			switch err {
			case user.ErrInvalidID:
				return web.NewRequestError(err, http.StatusBadRequest)
			case user.ErrNotFound:
				return web.NewRequestError(err, http.StatusNotFound)
			case user.ErrForbidden:
				return web.NewRequestError(err, http.StatusForbidden)
			default:
				return errors.Wrapf(err, "ID: %s", params["id"])
			}
		}
	}

	return web.Respond(ctx, w, smarket, http.StatusOK)
}

func (sg supermarketGroup) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	var ns supermarket.NewSuperMarket
	if err := web.Decode(r, &ns); err != nil {
		return errors.Wrapf(err, "unable to decode payload")
	}

	smarket, err := sg.supermarket.Create(ctx, v.TraceID, ns, v.Now)
	if err != nil {
		return errors.Wrapf(err, "Supermarket: %+v", &smarket)
	}

	return web.Respond(ctx, w, smarket, http.StatusCreated)
}

func (sg supermarketGroup) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	var smarket supermarket.UpdateSupermarket
	if err := web.Decode(r, &smarket); err != nil {
		return errors.Wrap(err, "unable to decode payload")
	}

	params := web.Params(r)
	err := sg.supermarket.Update(ctx, v.TraceID, claims, params["id"], smarket, v.Now)
	if err != nil {
		switch err {
		case user.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case user.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case user.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "ID: %s User: %+v", params["id"], &smarket)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (sg supermarketGroup) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	params := web.Params(r)
	err := sg.supermarket.Delete(ctx, v.TraceID, params["id"])

	if err != nil {
		if err != nil {
			switch err {
			case user.ErrInvalidID:
				return web.NewRequestError(err, http.StatusBadRequest)
			case user.ErrNotFound:
				return web.NewRequestError(err, http.StatusNotFound)
			case user.ErrForbidden:
				return web.NewRequestError(err, http.StatusForbidden)
			default:
				return errors.Wrapf(err, "ID: %s", params["id"])
			}
		}
	}
	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
