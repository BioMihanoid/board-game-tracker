package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"board-game-tracker/internal/models"
)

type CollectionService interface {
	AddGameInCollection(ctx context.Context, item *models.CollectionItem) (*models.CollectionItem, error)
	GetAllGamesFromCollection(ctx context.Context, filter *models.CollectionItem) ([]*models.CollectionItem, error)
	UpdateGameInCollection(ctx context.Context, item *models.CollectionItem) (*models.CollectionItem, error)
	DeleteGameFromCollection(ctx context.Context, userID, itemID string) error
	DeleteCollection(ctx context.Context, userID string) error
}

type CollectionHandler struct {
	collectionService CollectionService
}

func NewCollectionHandler(collectionService CollectionService) *CollectionHandler {
	return &CollectionHandler{
		collectionService: collectionService,
	}
}

func (h *CollectionHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.AddGame)
	r.Get("/", h.GetAllGames)
	r.Delete("/", h.DeleteCollection)
	r.Patch("/{itemID}", h.UpdateGame)
	r.Delete("/{itemID}", h.DeleteGame)

	return r
}

func (h *CollectionHandler) AddGame(w http.ResponseWriter, r *http.Request) {
	userID := "test_user_id" // TODO replace jwt-auth

	var req AddGameCollectionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.collectionService.AddGameInCollection(
		r.Context(),
		&models.CollectionItem{
			UserID:        userID,
			GameID:        req.GameID,
			Status:        req.Status,
			Condition:     req.Condition,
			PurchasePrice: req.PurchasePrice,
			Notes:         req.Notes,
		},
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *CollectionHandler) GetAllGames(w http.ResponseWriter, r *http.Request) {
	userID := "test_user_id" // TODO replace jwt-auth

	statusParam := r.URL.Query().Get("status")
	status := models.CollectionStatus(statusParam)
	if status != "" && !status.IsValid() {
		writeError(w, http.StatusBadRequest, "invalid status")
		return
	}

	resp, err := h.collectionService.GetAllGamesFromCollection(
		r.Context(),
		&models.CollectionItem{
			UserID: userID,
			Status: status,
		},
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *CollectionHandler) DeleteCollection(w http.ResponseWriter, r *http.Request) {
	userID := "test_user_id" // TODO replace jwt-auth

	if err := h.collectionService.DeleteCollection(r.Context(), userID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *CollectionHandler) UpdateGame(w http.ResponseWriter, r *http.Request) {
	userID := "test_user_id" // TODO replace jwt-auth

	var req UpdateGameInCollectionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.collectionService.UpdateGameInCollection(
		r.Context(),
		&models.CollectionItem{
			ID:            chi.URLParam(r, "itemID"),
			UserID:        userID,
			Status:        req.Status,
			Condition:     req.Condition,
			PurchasePrice: req.PurchasePrice,
			Notes:         req.Notes,
			SoldPrice:     req.SoldPrice,
		},
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *CollectionHandler) DeleteGame(w http.ResponseWriter, r *http.Request) {
	userID := "test_user_id" // TODO replace jwt-auth

	if err := h.collectionService.DeleteGameFromCollection(r.Context(), userID, chi.URLParam(r, "itemID")); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
