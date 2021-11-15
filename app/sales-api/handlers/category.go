package handlers

import (
	"context"
	"fmt"
	"github.com/adilkhan7/golangSoftProject/business/auth"
	"github.com/adilkhan7/golangSoftProject/business/data/category"
	"github.com/adilkhan7/golangSoftProject/business/data/user"
	"github.com/adilkhan7/golangSoftProject/foundation/web"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

type categoryGroup struct {
	category category.Category
	auth     *auth.Auth
}

func (gg categoryGroup) query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
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

	ctgry, err := gg.category.Query(ctx, v.TraceID, pageNumber, rowsPerPage)
	if err != nil {
		return errors.Wrap(err, "unable to query for category")
	}

	return web.Respond(ctx, w, ctgry, http.StatusOK)
}

func (gg categoryGroup) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	params := web.Params(r)
	ctgry, err := gg.category.QueryByID(ctx, v.TraceID, claims, params["id"])

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

	return web.Respond(ctx, w, ctgry, http.StatusOK)
}

func (gg categoryGroup) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	var nc category.NewCategory
	if err := web.Decode(r, &nc); err != nil {
		return errors.Wrapf(err, "unable to decode payload")
	}

	ctgry, err := gg.category.Create(ctx, v.TraceID, nc, v.Now)
	if err != nil {
		return errors.Wrapf(err, "Category: %+v", &ctgry)
	}

	return web.Respond(ctx, w, ctgry, http.StatusCreated)
}

func (gg categoryGroup) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	var ctgry category.UpdateCategory
	if err := web.Decode(r, &ctgry); err != nil {
		return errors.Wrap(err, "unable to decode payload")
	}

	params := web.Params(r)
	err := gg.category.Update(ctx, v.TraceID, claims, params["id"], ctgry, v.Now)
	if err != nil {
		switch err {
		case user.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case user.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case user.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "ID: %s User: %+v", params["id"], &ctgry)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (gg categoryGroup) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	params := web.Params(r)
	err := gg.category.Delete(ctx, v.TraceID, params["id"])

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
