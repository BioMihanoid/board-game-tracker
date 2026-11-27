package memory

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"board-game-tracker/internal/models"
)

type GameRepository struct {
	m  map[string]*models.Game
	rw sync.RWMutex
}

func NewGameRepository() *GameRepository {
	return &GameRepository{
		m: gamesMap,
	}
}

func (gr *GameRepository) GetGameById(_ context.Context, id string) (*models.Game, error) {
	if val, ok := gr.m[id]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("GetGameById: game not exist")
}

func (gr *GameRepository) GetAllGame(_ context.Context) (map[string]*models.Game, error) {
	if gr == nil || len(gr.m) == 0 {
		return nil, fmt.Errorf("GetAllGame: 0 games in db")
	}

	res := make(map[string]*models.Game, len(gr.m))

	for _, v := range gr.m {
		if v.Type == models.MainGameType {
			res[v.ID] = v
		}
	}

	return res, nil
}

func (gr *GameRepository) GetAdditionalContentByGameBaseID(_ context.Context, id string) (map[string]*models.Game, error) {
	if gr == nil || len(gr.m) == 0 {
		return nil, fmt.Errorf("GetAdditionalContentByGameBaseID: 0 games in db")
	}

	res := make(map[string]*models.Game, len(gr.m))

	for _, v := range gr.m {
		if v.Type == models.AdditionalContentGameType && v.ID == id {
			return nil, fmt.Errorf("GetAdditionalContentByGameBaseId: id additional content game")
		}

		if v.Type == models.AdditionalContentGameType && v.BaseGameID != nil && *v.BaseGameID == id {
			res[v.ID] = v
		}
	}

	if len(res) == 0 {
		return nil, fmt.Errorf("GetAdditionalContentByGameBaseId: game additional content not exist")
	}

	return res, nil
}

func (gr *GameRepository) SaveGame(ctx context.Context, game *models.Game) (string, error) {
	game.ID = strconv.Itoa(len(gr.m) + 1)

	gr.rw.Lock()
	gr.m[game.ID] = game
	gr.rw.Unlock()

	return game.ID, nil
}

func (gr *GameRepository) UpdateGame(ctx context.Context, game *models.Game) (*models.Game, error) {
	updateGame, err := gr.GetGameById(ctx, game.ID)
	if err != nil {
		return nil, err
	}

	gr.rw.Lock()

	updateGame.Name = game.Name
	updateGame.Publisher = game.Publisher

	gr.m[game.ID] = game

	gr.rw.Unlock()

	return game, nil
}

func (gr *GameRepository) DeleteGame(ctx context.Context, id string) error {
	gr.rw.Lock()
	delete(gr.m, id)
	gr.rw.Unlock()

	return nil
}
