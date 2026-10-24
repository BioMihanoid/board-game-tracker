package memory

import (
	"context"
	"fmt"
	"sync"

	"board-game-tracker/internal/model"
)

type GameRepository struct {
	m  map[int64]*model.Game
	rw sync.RWMutex
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
	res := make(map[int64]*model.Game, len(gr.m))
	if gr.m != nil && len(gr.m) > 0 {
		for _, v := range gr.m {
			if v.Type == model.MainGameType {
				res[v.ID] = v
			}
		}

		return res, nil
	}

	return nil, fmt.Errorf("GetAllGame: 0 games in db")
}

func (gr *GameRepository) GetAdditionalContentByGameBaseID(_ context.Context, id int64) (map[int64]*model.Game, error) {
	res := make(map[int64]*model.Game, len(gr.m))
	if gr.m != nil && len(gr.m) > 0 {
		for _, v := range gr.m {
			if v.Type == model.AdditionalContentGameType && v.ID == id {
				return nil, fmt.Errorf("GetAdditionalContentByGameBaseId: id additional content game")
			}

			if v.Type == model.AdditionalContentGameType && v.BaseGameID == id {
				res[v.ID] = v
			}
		}

		if len(res) == 0 {
			return nil, fmt.Errorf("GetAdditionalContentByGameBaseId: game additional content not exist")
		}

		return res, nil
	}

	return nil, fmt.Errorf("GetAdditionalContentByGameBaseID: 0 games in db")
}

func (gr *GameRepository) SaveGame(ctx context.Context, game *model.Game) (int64, error) {
	game.ID = int64(len(gr.m) + 1)

	gr.rw.Lock()
	gr.m[game.ID] = game
	gr.rw.Unlock()

	return game.ID, nil
}
