package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:generate templ generate

// for website, if you need to handle big file or stream, you must stay on stdlib or chi (just router)
// otherwise, you can use fiber or echo
// for api, you can use fiber or echo

type Msg struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func main() {
	r := chi.NewRouter()
	// r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	// r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// w.Write([]byte("hello world"))
		err := Hello("world").Render(context.Background(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		// w.Write([]byte("hello world"))
		// Hello("world").Render(context.Background(), w)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(Msg{Name: "test", Message: "test"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	err := http.ListenAndServe(":3333", r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}
