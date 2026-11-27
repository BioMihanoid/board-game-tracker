package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"board-game-tracker/internal/models"
)

type GameService interface {
	AddGame(ctx context.Context, game *models.Game) (string, error)
	GetGameById(ctx context.Context, id string) (*models.Game, error)
	UpdateGame(ctx context.Context, game *models.Game) (*models.Game, error)
	GetAllGame(ctx context.Context) (map[string]*models.Game, error)
	DeleteGame(ctx context.Context, id string) error
}

type GameHandler struct {
	service GameService
}

func NewGameHandler(service GameService) *GameHandler {
	return &GameHandler{
		service: service,
	}
}

func (h *GameHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.Create)
	r.Get("/", h.GetAllGame)
	r.Get("/{id}", h.GetGameById)
	r.Put("/{id}", h.UpdateGame)
	r.Delete("/{id}", h.DeleteGame)

	return r
}

func (h *GameHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateGameRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.Name == "" {
		writeJSON(w, http.StatusBadRequest, "title is required")
		return
	}

	if req.Publisher == "" {
		writeJSON(w, http.StatusBadRequest, "publisher is required")
		return
	}

	game, err := h.service.AddGame(
		r.Context(),
		&models.Game{
			Name:      req.Name,
			Publisher: req.Publisher,
		},
	)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, game)
}

func (h *GameHandler) GetGameById(w http.ResponseWriter, r *http.Request) {
	game, err := h.service.GetGameById(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, game)
}

func (h *GameHandler) UpdateGame(w http.ResponseWriter, r *http.Request) {
	var req UpdateGameRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	game, err := h.service.UpdateGame(
		r.Context(),
		&models.Game{
			ID:        chi.URLParam(r, "id"),
			Name:      req.Name,
			Publisher: req.Publisher,
		},
	)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, game)
}

func (h *GameHandler) GetAllGame(w http.ResponseWriter, r *http.Request) {
	games, err := h.service.GetAllGame(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, games)
}

func (h *GameHandler) DeleteGame(w http.ResponseWriter, r *http.Request) {
	err := h.service.DeleteGame(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, true)
}
