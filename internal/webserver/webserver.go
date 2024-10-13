package webserver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const defaultReadHeaderTimeout = 5 * time.Second

// Webserver is a struct representing a web server
type Webserver struct {
	srv    *http.Server
	router *chi.Mux
	db     *sql.DB
	port   int
}

// NewWebserver is a function that creates a new Webserver
func NewWebserver(db *sql.DB, port int) *Webserver {
	w := &Webserver{
		db:   db,
		port: port,
	}
	w.initRouter()
	w.initRoutes()
	return w
}

// initRouter is a method that initializes the router
func (w *Webserver) initRouter() {
	w.router = chi.NewRouter()
	w.router.Use(middleware.RequestID)
	w.router.Use(middleware.Logger)
	w.router.Use(middleware.Recoverer)
}

// Run is a method that starts the web server
func (w *Webserver) Run() error {
	w.srv = &http.Server{
		Addr:              fmt.Sprintf(":%d", w.port),
		Handler:           w.router,
		ReadHeaderTimeout: defaultReadHeaderTimeout,
	}
	err := w.srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("srv.ListenAndServe: %w", err)
	}
	return nil
}

// Shutdown is a method that shuts down the web server
func (w *Webserver) Shutdown() error {
	err := w.srv.Shutdown(context.TODO())
	if err != nil {
		return fmt.Errorf("srv.Shutdown: %w", err)
	}
	return nil
}
