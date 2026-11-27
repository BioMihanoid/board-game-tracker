package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Dependencies struct {
	GameService GameService
}

func NewRouter(deps Dependencies) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(5 * time.Second))

	r.Get(
		"/ping",
		func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
		},
	)

	gameHandler := NewGameHandler(deps.GameService)

	r.Route(
		"/api/v1",
		func(r chi.Router) {
			r.Mount("/games", gameHandler.Routes())
		},
	)

	return r
}
