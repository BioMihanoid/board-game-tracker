package model

type Game struct {
	ID                 int64
	BaseGameID         int64
	Name               string
	Publisher          string
	Type               GameType
	ReleaseYear        int
	MinPlayers         int
	MaxPlayers         int
	PlaytimeMinutesMin int
	PlaytimeMinutesMax int
	CostKopecks        int64
}
