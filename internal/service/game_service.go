package service

import (
	"context"

	"board-game-tracker/internal/models"
)

type GameRepository interface {
	GetGameByID(ctx context.Context, id string) (*models.Game, error)
	GetAllGames(ctx context.Context) ([]*models.Game, error)
	SaveGame(ctx context.Context, game *models.Game) (string, error)
	GetAdditionalContentByGameBaseID(ctx context.Context, id string) ([]*models.Game, error)
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

func (gs *GameService) GetGameByID(ctx context.Context, id string) (*models.Game, error) {
	return gs.repo.GetGameByID(ctx, id)
}

func (gs *GameService) GetAllGames(ctx context.Context) ([]*models.Game, error) {
	return gs.repo.GetAllGames(ctx)
}

func (gs *GameService) GetAdditionalContentByGameBaseID(ctx context.Context, id string) ([]*models.Game, error) {
	return gs.repo.GetAdditionalContentByGameBaseID(ctx, id)
}

func (gs *GameService) UpdateGame(ctx context.Context, game *models.Game) (*models.Game, error) {
	return gs.repo.UpdateGame(ctx, game)
}

func (gs *GameService) DeleteGame(ctx context.Context, id string) error {
	return gs.repo.DeleteGame(ctx, id)
}
