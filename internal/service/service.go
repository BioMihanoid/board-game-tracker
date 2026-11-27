package service

import (
	"context"

	"board-game-tracker/internal/models"
)

type GameRepository interface {
	GetGameById(ctx context.Context, id string) (*models.Game, error)
	GetAllGame(ctx context.Context) (map[string]*models.Game, error)
	SaveGame(ctx context.Context, game *models.Game) (string, error)
	GetAdditionalContentByGameBaseID(ctx context.Context, id string) (map[string]*models.Game, error)
	UpdateGame(ctx context.Context, game *models.Game) (*models.Game, error)
	DeleteGame(ctx context.Context, id string) error
}

type GameService struct {
	repo GameRepository
}

func NewGameService(repo GameRepository) *GameService {
	return &GameService{
		repo: repo,
	}
}

func (gs *GameService) AddGame(ctx context.Context, game *models.Game) (string, error) {
	return gs.repo.SaveGame(ctx, game)
}

func (gs *GameService) GetGameById(ctx context.Context, id string) (*models.Game, error) {
	game, err := gs.repo.GetGameById(ctx, id)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (gs *GameService) GetAllGame(ctx context.Context) (map[string]*models.Game, error) {
	games, err := gs.repo.GetAllGame(ctx)
	if err != nil {
		return nil, err
	}

	return games, nil
}

func (gs *GameService) GetAdditionalContentByGameBaseID(ctx context.Context, id string) (map[string]*models.Game, error) {
	games, err := gs.repo.GetAdditionalContentByGameBaseID(ctx, id)
	if err != nil {
		return nil, err
	}

	return games, nil
}

func (gs *GameService) UpdateGame(ctx context.Context, game *models.Game) (*models.Game, error) {
	return gs.repo.UpdateGame(ctx, game)
}

func (gs *GameService) DeleteGame(ctx context.Context, id string) error {
	return gs.repo.DeleteGame(ctx, id)
}
