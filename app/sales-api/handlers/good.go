package handlers

import (
	"context"
	"fmt"
	"github.com/adilkhan7/golangSoftProject/business/auth"
	"github.com/adilkhan7/golangSoftProject/business/data/good"
	"github.com/adilkhan7/golangSoftProject/business/data/user"
	"github.com/adilkhan7/golangSoftProject/foundation/web"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

type goodGroup struct {
	good good.Good
	auth *auth.Auth
}

func (gg goodGroup) query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
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

	gd, err := gg.good.Query(ctx, v.TraceID, pageNumber, rowsPerPage)
	if err != nil {
		return errors.Wrap(err, "unable to query for good")
	}

	return web.Respond(ctx, w, gd, http.StatusOK)
}

func (gg goodGroup) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	params := web.Params(r)
	gd, err := gg.good.QueryByID(ctx, v.TraceID, claims, params["id"])

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

	return web.Respond(ctx, w, gd, http.StatusOK)
}

func (gg goodGroup) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	var ng good.NewGood
	if err := web.Decode(r, &ng); err != nil {
		return errors.Wrapf(err, "unable to decode payload")
	}

	gs, err := gg.good.Create(ctx, v.TraceID, ng, v.Now)
	if err != nil {
		return errors.Wrapf(err, "Good: %+v", &gs)
	}

	return web.Respond(ctx, w, gs, http.StatusCreated)
}

func (gg goodGroup) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	var gd good.UpdateGood
	if err := web.Decode(r, &gd); err != nil {
		return errors.Wrap(err, "unable to decode payload")
	}

	params := web.Params(r)
	err := gg.good.Update(ctx, v.TraceID, claims, params["id"], gd, v.Now)
	if err != nil {
		switch err {
		case user.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case user.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case user.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "ID: %s User: %+v", params["id"], &gd)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (gg goodGroup) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	params := web.Params(r)
	err := gg.good.Delete(ctx, v.TraceID, params["id"])

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
