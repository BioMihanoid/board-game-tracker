package api

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
