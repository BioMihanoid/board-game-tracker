package model

type Game struct {
	ID                 int64    `json:"id"`
	BaseGameID         int64    `json:"base_game_id,omitempty"`
	Name               string   `json:"name"`
	Publisher          string   `json:"publisher,omitempty"`
	Type               GameType `json:"type"`
	ReleaseYear        int      `json:"release_year,omitempty"`
	MinPlayers         int      `json:"min_players,omitempty"`
	MaxPlayers         int      `json:"max_players,omitempty"`
	PlaytimeMinutesMin int      `json:"playtime_minutes_min,omitempty"`
	PlaytimeMinutesMax int      `json:"playtime_minutes_max,omitempty"`
	// CostKopecks — стоимость в копейках, чтобы суммы считались без потерь.
	CostKopecks int64 `json:"cost_kopecks"`
}
