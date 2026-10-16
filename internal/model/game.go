package model

type Game struct {
	ID                 int64
	Name               string
	Publisher          string
	Type               string
	ReleaseYear        int
	MinPlayers         int
	MaxPlayers         int
	PlaytimeMinutesMin int
	PlaytimeMinutesMax int
	Cost               float64
}
