package _map

import (
	"context"
	"fmt"

	"board-game-tracker/internal/model"
)

type GameRepository struct {
	m map[int64]*model.Game
}

func NewGameRepository() *GameRepository {
	return &GameRepository{
		m: gamesMap,
	}
}

func (r *GameRepository) GetGameById(_ context.Context, id int64) (*model.Game, error) {
	if val, ok := r.m[id]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("GetGameById: game not exist")
}

func (r *GameRepository) GetAllGame(_ context.Context) (map[int64]*model.Game, error) {
	if r.m != nil && len(r.m) > 0 {
		return r.m, nil
	}

	return nil, fmt.Errorf("GetAllGame: 0 games in db")
}
