package memory

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"board-game-tracker/internal/models"
)

var ErrNotFound = errors.New("not found")

type GameRepository struct {
	m      map[string]*models.Game
	nextID int
	rw     sync.RWMutex
}

func NewGameRepository() *GameRepository {
	copied := make(map[string]*models.Game, len(gamesMap))
	for k, v := range gamesMap {
		if v != nil {
			clone := *v
			copied[k] = &clone
		}
	}

	return &GameRepository{
		m:      copied,
		nextID: len(copied),
	}
}

func (gr *GameRepository) GetGameByID(_ context.Context, id string) (*models.Game, error) {
	gr.rw.RLock()
	defer gr.rw.RUnlock()

	val, ok := gr.m[id]
	if !ok {
		return nil, fmt.Errorf("GetGameByID: %w", ErrNotFound)
	}

	clone := *val
	return &clone, nil
}

func (gr *GameRepository) GetAllGames(_ context.Context) ([]*models.Game, error) {
	gr.rw.RLock()
	defer gr.rw.RUnlock()

	res := make([]*models.Game, 0, len(gr.m))
	for _, v := range gr.m {
		if v.Type == models.GameTypeMain {
			clone := *v
			res = append(res, &clone)
		}
	}

	return res, nil
}

func (gr *GameRepository) GetAdditionalContentByGameBaseID(_ context.Context, baseID string) ([]*models.Game, error) {
	gr.rw.RLock()
	defer gr.rw.RUnlock()

	if target, exists := gr.m[baseID]; exists && target.Type == models.ContentGameTypeAdditional {
		return nil, fmt.Errorf("GetAdditionalContentByGameBaseID: provided id belongs to an expansion, not a base game")
	}

	res := make([]*models.Game, 0)
	for _, v := range gr.m {
		if v.Type == models.ContentGameTypeAdditional && v.BaseGameID != nil && *v.BaseGameID == baseID {
			clone := *v
			res = append(res, &clone)
		}
	}

	return res, nil
}

func (gr *GameRepository) SaveGame(_ context.Context, game *models.Game) (string, error) {
	if game == nil {
		return "", errors.New("SaveGame: game is nil")
	}

	gr.rw.Lock()
	defer gr.rw.Unlock()

	if game.ID == "" {
		gr.nextID++
		game.ID = strconv.Itoa(gr.nextID)
	}

	clone := *game
	gr.m[clone.ID] = &clone

	return clone.ID, nil
}

func (gr *GameRepository) UpdateGame(_ context.Context, game *models.Game) (*models.Game, error) {
	if game == nil {
		return nil, errors.New("UpdateGame: game is nil")
	}

	gr.rw.Lock()
	defer gr.rw.Unlock()

	if _, ok := gr.m[game.ID]; !ok {
		return nil, fmt.Errorf("UpdateGame: %w", ErrNotFound)
	}

	clone := *game
	gr.m[game.ID] = &clone

	return &clone, nil
}

func (gr *GameRepository) DeleteGame(_ context.Context, id string) error {
	gr.rw.Lock()
	defer gr.rw.Unlock()

	if _, ok := gr.m[id]; !ok {
		return fmt.Errorf("DeleteGame: %w", ErrNotFound)
	}

	delete(gr.m, id)
	return nil
}
