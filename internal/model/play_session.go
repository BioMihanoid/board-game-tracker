package model

import "time"

type PlaySession struct {
	ID       int64
	GameID   int64
	PlayedAt time.Time
	Minutes  int
	Players  int
	Notes    string
}
