package console

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"

	"board-game-tracker/internal/model"
)

type GameService interface {
	GetGameById(ctx context.Context, id int64) (*model.Game, error)
	GetAllGame(ctx context.Context) (map[int64]*model.Game, error)
	AddGame(ctx context.Context, game *model.Game) error
}

type Handler struct {
	gameService GameService
}

func NewHandler(gameService GameService) *Handler {
	return &Handler{
		gameService: gameService,
	}
}

func (h *Handler) PrintGameByID(ctx context.Context, id int64) {
	game, err := h.gameService.GetGameById(ctx, id)
	if err != nil {
		return
	}

	fmt.Printf("Игра: %s\nТип: %s\nСтоимость: %0.2f\nИздатель: %s\n\n",
		game.Name,
		game.Type,
		game.Cost,
		game.Publisher,
	)
}

func (h *Handler) PrintAllGames(ctx context.Context) {
	games, err := h.gameService.GetAllGame(ctx)
	if err != nil {
		return
	}

	for _, game := range games {
		fmt.Printf("Игра: %s\nТип: %s\nСтоимость: %0.2f\nИздатель: %s\n\n",
			game.Name,
			game.Type,
			game.Cost,
			game.Publisher,
		)
	}
}

func (h *Handler) AddGame(ctx context.Context) error {
	choice := bufio.NewScanner(os.Stdin)
	game := &model.Game{}

	fmt.Print("Введите название игры: ")
	choice.Scan()
	game.Name = choice.Text()

	fmt.Print("Введите тип игры: ")
	choice.Scan()
	game.Type = choice.Text()

	fmt.Print("Введите стоимость игры: ")
	choice.Scan()
	val, err := strconv.ParseFloat(choice.Text(), 64)
	if err != nil {
		return err
	}
	game.Cost = val

	fmt.Print("Введите издателя игры: ")
	choice.Scan()
	game.Publisher = choice.Text()

	return h.gameService.AddGame(ctx, game)
}
