package memory

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"board-game-tracker/internal/model"
)

type GameRepository struct {
	m      map[int64]*model.Game
	nextID int64
	rw     sync.RWMutex
}

func (gr *GameRepository) GetByID(_ context.Context, id int64) (*model.Game, error) {
	gr.rw.RLock()
	defer gr.rw.RUnlock()

	game, ok := gr.m[id]
	if !ok {
		return nil, fmt.Errorf("get game by id %d: %w", id, model.ErrNotFound)
	}

	// Отдаём копию, чтобы вызывающий не мог менять состояние репозитория в обход сервиса.
	res := *game

	return &res, nil
}

func (gr *GameRepository) ListBaseGames(_ context.Context) ([]*model.Game, error) {
	gr.rw.RLock()
	defer gr.rw.RUnlock()

	return gr.filter(func(game *model.Game) bool {
		return game.Type == model.MainGameType
	}), nil
}

func (gr *GameRepository) ListExpansions(_ context.Context, baseID int64) ([]*model.Game, error) {
	gr.rw.RLock()
	defer gr.rw.RUnlock()

	return gr.filter(func(game *model.Game) bool {
		return game.BaseGameID == baseID && game.Type != model.MainGameType
	}), nil
}

// filter вызывается под RLock, поэтому копирование делает он же.
func (gr *GameRepository) filter(match func(*model.Game) bool) []*model.Game {
	res := make([]*model.Game, 0, len(gr.m))
	for _, game := range gr.m {
		if !match(game) {
			continue
		}

		copied := *game
		res = append(res, &copied)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].ID < res[j].ID
	})

	return res
}

func (gr *GameRepository) Save(_ context.Context, game *model.Game) (int64, error) {
	if game == nil {
		return 0, fmt.Errorf("%w: пустая игра", model.ErrInvalidGame)
	}

	gr.rw.Lock()
	defer gr.rw.Unlock()

	id := gr.nextID
	gr.nextID++

	stored := *game
	stored.ID = id
	gr.m[id] = &stored

	game.ID = id

	return id, nil
}
