package models

type GameType string

const (
	MainGameType              GameType = "Основная игра"
	AdditionalContentGameType GameType = "Дополнение"
	PromoGameType             GameType = "Промо"
)
