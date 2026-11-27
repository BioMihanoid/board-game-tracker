package models

import "time"

type Game struct {
	ID                 string
	BaseGameID         *string
	Name               string
	OriginalName       string
	Publisher          string
	Genres             []string
	Type               GameType
	ReleaseYear        int
	MinPlayers         int
	MaxPlayers         int
	BestPlayers        []int
	MinPlaytimeMinutes int
	MaxPlaytimeMinutes int
	Cost               float64
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
