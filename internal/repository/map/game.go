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

func (gr *GameRepository) GetGameById(_ context.Context, id int64) (*model.Game, error) {
	if val, ok := gr.m[id]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("GetGameById: game not exist")
}

func (gr *GameRepository) GetAllGame(_ context.Context) (map[int64]*model.Game, error) {
	if gr.m != nil && len(gr.m) > 0 {
		return gr.m, nil
	}

	return nil, fmt.Errorf("GetAllGame: 0 games in db")
}

func (gr *GameRepository) SaveGame(ctx context.Context, game *model.Game) error {
	// TODO допилить тут синк мьютекст!!!
	gr.m[int64(len(gr.m))] = game
	return nil
}
