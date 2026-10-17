package service

import (
	"context"

	"board-game-tracker/internal/model"
)

type GameRepository interface {
	GetGameById(ctx context.Context, id int64) (*model.Game, error)
	GetAllGame(ctx context.Context) (map[int64]*model.Game, error)
	SaveGame(ctx context.Context, game *model.Game) error
}

type GameService struct {
	repo GameRepository
}

func NewGameService(repo GameRepository) *GameService {
	return &GameService{
		repo: repo,
	}
}

func (gs *GameService) AddGame(ctx context.Context, game *model.Game) error {
	return gs.repo.SaveGame(ctx, game)
}

func (gs *GameService) GetGameById(ctx context.Context, id int64) (*model.Game, error) {
	game, err := gs.repo.GetGameById(ctx, id)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (gs *GameService) GetAllGame(ctx context.Context) (map[int64]*model.Game, error) {
	games, err := gs.repo.GetAllGame(ctx)
	if err != nil {
		return nil, err
	}

	return games, nil
}
