package models

type GameType string

const (
	GameTypeMain              GameType = "Основная игра"
	ContentGameTypeAdditional GameType = "Дополнение"
	GameTypePromo             GameType = "Промо"
)
