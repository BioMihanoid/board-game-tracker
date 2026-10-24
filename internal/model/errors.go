package model

import "errors"

var (
	// ErrNotFound — запись не найдена.
	ErrNotFound = errors.New("not found")
	// ErrNotBaseGame — операция требует основную игру, а передана другая.
	ErrNotBaseGame = errors.New("not a base game")
	// ErrInvalidGame — данные игры не проходят валидацию.
	ErrInvalidGame = errors.New("invalid game")
)
