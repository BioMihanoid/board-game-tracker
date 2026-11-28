package api

import "board-game-tracker/internal/models"

// --- GAMES ---

type CreateGameRequest struct {
	Name               string `json:"name"`
	Publisher          string `json:"publisher"`
	MinPlayers         int    `json:"min_player"`
	MaxPlayers         int    `json:"max_player"`
	MinPlaytimeMinutes int    `json:"min_play_time_minute"`
	MaxPlaytimeMinutes int    `json:"max_play_time_minute"`
}

type UpdateGameRequest struct {
	Name      string `json:"name"`
	Publisher string `json:"publisher"`
}

// --- Collections ---

type AddGameCollectionRequest struct {
	GameID        string                  `json:"game_id"`
	Status        models.CollectionStatus `json:"status"`
	Condition     models.ItemCondition    `json:"condition"`
	PurchasePrice *int                    `json:"purchase_price"`
	TargetPrice   *int                    `json:"target_price"`
	Notes         string                  `json:"notes"`
}

type GetAllGamesRequest struct {
	Status models.CollectionStatus `json:"status"`
}

type UpdateGameInCollectionRequest struct {
	GameID        string                  `json:"game_id"`
	Status        models.CollectionStatus `json:"status"`
	Condition     models.ItemCondition    `json:"condition"`
	PurchasePrice *int                    `json:"purchase_price"`
	TargetPrice   *int                    `json:"target_price"`
	Notes         string                  `json:"notes"`
	SoldPrice     *int                    `json:"sold_price"`
}
