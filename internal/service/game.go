package service

import (
	"context"
	"fmt"

	"board-game-tracker/internal/model"
)

type GameRepository interface {
	GetGameById(ctx context.Context, id int64) (*model.Game, error)
	GetAllGame(ctx context.Context) (map[int64]*model.Game, error)
}

type GameService struct {
	repo GameRepository
}

func NewGameService(repo GameRepository) *GameService {
	return &GameService{
		repo: repo,
	}
}

func (gs *GameService) PrintGames(ctx context.Context) {
	games, err := gs.repo.GetAllGame(ctx)
	if err != nil {
		fmt.Printf("Игры отсуствуют\n")
		return
	}

	for _, g := range games {
		fmt.Printf("Название игры - %s\nТип - %s\nЦена - %0.2f\nИздатель - %s\n\n",
			g.Name,
			g.Type,
			g.Cost,
			g.Publisher,
		)
	}
}
