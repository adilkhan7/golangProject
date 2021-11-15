package handlers

import (
	"github.com/adilkhan7/golangSoftProject/business/auth"
	"github.com/adilkhan7/golangSoftProject/business/data/category"
	"github.com/adilkhan7/golangSoftProject/business/data/good"
	"github.com/adilkhan7/golangSoftProject/business/data/supermarket"
	"github.com/adilkhan7/golangSoftProject/business/data/user"
	"github.com/adilkhan7/golangSoftProject/business/mid"
	"github.com/adilkhan7/golangSoftProject/foundation/web"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
)

func API(build string, shutdown chan os.Signal, log *log.Logger, a *auth.Auth, db *sqlx.DB) *web.App {
	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics(log))

	cg := CheckGroup{
		build: build,
		db:    db,
	}

	app.Handle(http.MethodGet, "/readiness", cg.readiness)
	app.Handle(http.MethodGet, "/liveness", cg.liveness)

	ug := userGroup{
		user: user.New(log, db),
		auth: a,
	}
	app.Handle(http.MethodGet, "/users/:page/:rows", ug.query, mid.Authenticate(a), mid.Authorize(log, auth.RoleAdmin))
	app.Handle(http.MethodGet, "/users/:id", ug.queryByID, mid.Authenticate(a))
	app.Handle(http.MethodGet, "/users/token/:kid", ug.token)
	app.Handle(http.MethodPost, "/users", ug.create, mid.Authenticate(a), mid.Authorize(log, auth.RoleAdmin))
	app.Handle(http.MethodPut, "/users/:id", ug.update, mid.Authenticate(a), mid.Authorize(log, auth.RoleAdmin))
	app.Handle(http.MethodDelete, "/users/:id", ug.delete, mid.Authenticate(a), mid.Authorize(log, auth.RoleAdmin))

	sg := supermarketGroup{
		supermarket: supermarket.New(log, db),
		auth:        a,
	}
	app.Handle(http.MethodGet, "/supermarket/:page/:rows", sg.query, mid.Authenticate(a))
	app.Handle(http.MethodGet, "/supermarket/:id", sg.queryByID, mid.Authenticate(a))
	app.Handle(http.MethodPost, "/supermarket", sg.create, mid.Authenticate(a))
	app.Handle(http.MethodPut, "/supermarket/:id", sg.update, mid.Authenticate(a))
	app.Handle(http.MethodDelete, "/supermarket/:id", sg.delete, mid.Authenticate(a))

	ctg := categoryGroup{
		category: category.New(log, db),
		auth:     a,
	}
	app.Handle(http.MethodGet, "/category/:page/:rows", ctg.query, mid.Authenticate(a))
	app.Handle(http.MethodGet, "/category/:id", ctg.queryByID, mid.Authenticate(a))
	app.Handle(http.MethodPost, "/category", ctg.create, mid.Authenticate(a))
	app.Handle(http.MethodPut, "/category/:id", ctg.update, mid.Authenticate(a))
	app.Handle(http.MethodDelete, "/category/:id", ctg.delete, mid.Authenticate(a))

	gg := goodGroup{
		good: good.New(log, db),
		auth: a,
	}
	app.Handle(http.MethodGet, "/good/:page/:rows", gg.query, mid.Authenticate(a))
	app.Handle(http.MethodGet, "/good/:id", gg.queryByID, mid.Authenticate(a))
	app.Handle(http.MethodPost, "/good", gg.create, mid.Authenticate(a))
	app.Handle(http.MethodPut, "/good/:id", gg.update, mid.Authenticate(a))
	app.Handle(http.MethodDelete, "/good/:id", gg.delete, mid.Authenticate(a))

	return app
}
