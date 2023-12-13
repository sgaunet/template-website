package webserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sgaunet/template-website/internal/views"
)

func (w *Webserver) initRoutes() {
	w.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		err := views.Hello("world").Render(context.Background(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	w.router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(Msg{Name: "test", Message: "test"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	w.router.Handle("/bootstrap-5.1.3-dist/*", views.BootStrapHandler("/bootstrap-5.1.3-dist/"))
}
