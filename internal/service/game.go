package service

import (
	"context"
	"fmt"
	"strings"

	"board-game-tracker/internal/model"
)

type GameRepository interface {
	GetByID(ctx context.Context, id int64) (*model.Game, error)
	ListBaseGames(ctx context.Context) ([]*model.Game, error)
	ListExpansions(ctx context.Context, baseID int64) ([]*model.Game, error)
	Save(ctx context.Context, game *model.Game) (int64, error)
}

type GameService struct {
	repo GameRepository
}

func NewGameService(repo GameRepository) *GameService {
	return &GameService{
		repo: repo,
	}
}

// AddGame проверяет данные игры и сохраняет её.
// Для дополнения и промо game.Name — это только суффикс названия:
// полное имя ("Базовая игра. Суффикс") собирает сервис.
func (gs *GameService) AddGame(ctx context.Context, game *model.Game) (int64, error) {
	if game == nil {
		return 0, fmt.Errorf("%w: пустая игра", model.ErrInvalidGame)
	}

	game.Name = strings.TrimSpace(game.Name)
	if game.Name == "" {
		return 0, fmt.Errorf("%w: название не может быть пустым", model.ErrInvalidGame)
	}

	if game.CostKopecks < 0 {
		return 0, fmt.Errorf("%w: стоимость не может быть отрицательной", model.ErrInvalidGame)
	}

	switch game.Type {
	case model.MainGameType:
		if game.BaseGameID != 0 {
			return 0, fmt.Errorf("%w: у основной игры не может быть базовой", model.ErrInvalidGame)
		}
	case model.AdditionalContentGameType, model.PromoGameType:
		base, err := gs.getBaseGame(ctx, game.BaseGameID)
		if err != nil {
			return 0, err
		}

		game.Name = base.Name + ". " + game.Name
		game.Publisher = base.Publisher
	default:
		return 0, fmt.Errorf("%w: неизвестный тип игры %q", model.ErrInvalidGame, game.Type)
	}

	return gs.repo.Save(ctx, game)
}

func (gs *GameService) GetByID(ctx context.Context, id int64) (*model.Game, error) {
	return gs.repo.GetByID(ctx, id)
}

func (gs *GameService) ListBaseGames(ctx context.Context) ([]*model.Game, error) {
	return gs.repo.ListBaseGames(ctx)
}

func (gs *GameService) ListExpansions(ctx context.Context, baseID int64) ([]*model.Game, error) {
	if _, err := gs.getBaseGame(ctx, baseID); err != nil {
		return nil, err
	}

	return gs.repo.ListExpansions(ctx, baseID)
}

func (gs *GameService) getBaseGame(ctx context.Context, id int64) (*model.Game, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: не указана базовая игра", model.ErrInvalidGame)
	}

	base, err := gs.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if base.Type != model.MainGameType {
		return nil, fmt.Errorf("%w: игра %q не является основной", model.ErrNotBaseGame, base.Name)
	}

	return base, nil
}
