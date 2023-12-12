package webserver

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Webserver struct {
	srv    *http.Server
	router *chi.Mux
	db     *sql.DB
	port   int
}

func NewWebserver(db *sql.DB, port int) *Webserver {
	w := &Webserver{
		db:   db,
		port: port,
	}
	w.initRouter()
	w.initRoutes()
	return w
}

func (w *Webserver) initRouter() {
	w.router = chi.NewRouter()
	w.router.Use(middleware.RequestID)
	w.router.Use(middleware.Logger)
	w.router.Use(middleware.Recoverer)
}

func (w *Webserver) Run() error {
	w.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", w.port),
		Handler: w.router,
	}
	return w.srv.ListenAndServe()
}

func (w *Webserver) Shutdown() error {
	return w.srv.Shutdown(context.TODO())
}
