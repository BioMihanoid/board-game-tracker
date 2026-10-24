package console

import (
	"bufio"
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"board-game-tracker/internal/model"
)

type GameService interface {
	GetGameById(ctx context.Context, id int64) (*model.Game, error)
	GetAllGame(ctx context.Context) (map[int64]*model.Game, error)
	AddGame(ctx context.Context, game *model.Game) (int64, error)
	GetAdditionalContentByGameBaseID(ctx context.Context, id int64) (map[int64]*model.Game, error)
}

type Handler struct {
	gameService GameService
	scanner     *bufio.Scanner
}

func NewHandler(gameService GameService, scanner *bufio.Scanner) *Handler {
	return &Handler{
		gameService: gameService,
		scanner:     scanner,
	}
}

func (h *Handler) PrintGameByID(ctx context.Context) {
	var (
		id   int64
		game *model.Game
		err  error
	)
	for {
		id = h.getInt64Value("id")
		game, err = h.gameService.GetGameById(ctx, id)
		if err != nil {
			fmt.Printf("Игра с таким id отсуствует, попробуйте ещё раз\n\n")
			continue
		}

		break
	}

	fmt.Printf("\nИгра: %s\nТип: %s\nСтоимость: %0.2f\nИздатель: %s\n\n",
		game.Name,
		game.Type,
		game.Cost,
		game.Publisher,
	)
}

func (h *Handler) PrintAllGames(ctx context.Context) {
	games, err := h.gameService.GetAllGame(ctx)
	if err != nil {
		fmt.Printf("Не получилось получить все игры\n\n")
		return
	}

	for _, game := range games {
		fmt.Printf("ID: %d\nИгра: %s\nТип: %s\nСтоимость: %0.2f\nИздатель: %s\n\n",
			game.ID,
			game.Name,
			game.Type,
			game.Cost,
			game.Publisher,
		)
	}
}

func (h *Handler) PrintAllAdditionalContentByGameID(ctx context.Context) {
	var id int64
	for {
		id = h.getInt64Value("id базовой")
		_, err := h.gameService.GetGameById(ctx, id)
		if err != nil {
			fmt.Printf("Игра с таким id отсуствует, попробуйте ещё раз\n\n")
			continue
		}

		break
	}

	games, err := h.gameService.GetAdditionalContentByGameBaseID(ctx, id)
	if err != nil {
		fmt.Printf("Не получилось получить дополнения игры\n\n")
		return
	}

	for _, game := range games {
		fmt.Printf("Название: %s\nСтоимость: %0.2f\nИздатель: %s\n\n",
			game.Name,
			game.Cost,
			game.Publisher,
		)
	}
}

func (h *Handler) AddGame(ctx context.Context) (int64, error) {
	game := &model.Game{}

	game.Type = h.getGameType()

	if game.Type == model.AdditionalContentGameType || game.Type == model.PromoGameType {
		var baseGame *model.Game
		for {
			baseGame, err := h.gameService.GetGameById(ctx, h.getInt64Value("id базовой игры"))
			if err != nil {
				fmt.Printf("Игра с таким id отсуствует, попробуйте ещё раз\n\n")
				continue
			}

			game.BaseGameID = baseGame.ID
			break
		}

		tag := "дополнения"
		if game.Type == model.PromoGameType {
			tag = "промо набора"
		}

		game.Name = baseGame.Name + ". " + h.getStringValue(fmt.Sprintf("название %s", tag))
		game.Cost = h.getFloatValue(fmt.Sprintf("стоимость %s", tag))
		game.Publisher = baseGame.Publisher

		return h.gameService.AddGame(ctx, game)
	}

	game.Name = h.getStringValue("название")
	game.Cost = h.getFloatValue("стоимость")
	game.Publisher = h.getStringValue("издателя")

	return h.gameService.AddGame(ctx, game)
}

func (h *Handler) getGameType() model.GameType {
	for {
		fmt.Printf(
			"Доступные типы игры:\n1. %s\n2. %s\n3. %s\n",
			model.MainGameType,
			model.AdditionalContentGameType,
			model.PromoGameType,
		)
		choice := h.getInt64Value("тип")
		switch choice {
		case 1:
			return model.MainGameType
		case 2:
			return model.AdditionalContentGameType
		case 3:
			return model.PromoGameType
		default:
			fmt.Printf("Выбран не существующий тип игры, попробуйте ещё раз\n\n")
		}
	}
}

func (h *Handler) getStringValue(in string) string {
	for {
		fmt.Printf("Введите %s игры: ", in)
		ok := h.scanner.Scan()
		if !ok {
			fmt.Printf("Произошла ошибка, попробуйте ещё раз\n\n")
			continue
		}

		val := h.scanner.Text()
		switch val {
		case "":
			fmt.Printf("Пустой ввод, попробуйте ещё раз\n\n")
		default:
			return val
		}
	}
}

func (h *Handler) getFloatValue(in string) float64 {
	for {
		fmt.Printf("Введите %s игры: ", in)
		ok := h.scanner.Scan()
		if !ok {
			fmt.Printf("Произошла ошибка, попробуйте ещё раз\n\n")
			continue
		}

		input := h.scanner.Text()
		if input == "" {
			fmt.Printf("Пустой ввод, попробуйте ещё раз\n\n")
			continue
		}

		val, err := strconv.ParseFloat(strings.TrimSpace(input), 64)
		if err != nil {
			fmt.Printf("Произошла ошибка конвертации, попробуйте ещё раз\n\n")
			continue
		}

		switch {
		case val < 0:
			fmt.Printf("Значение не может быть отрицательным, попробуйте ещё раз\n\n")
		case math.IsNaN(val) || math.IsInf(val, 0):
			fmt.Printf("Некорректный ввод (NaN/Infinity), попробуйте ещё раз\n\n")
		default:
			return val
		}
	}
}

func (h *Handler) getInt64Value(in string) int64 {
	for {
		fmt.Printf("Введите %s игры: ", in)
		ok := h.scanner.Scan()
		if !ok {
			fmt.Printf("Произошла ошибка, попробуйте ещё раз\n\n")
			continue
		}

		input := h.scanner.Text()
		if input == "" {
			fmt.Printf("Пустой ввод, попробуйте ещё раз\n\n")
			continue
		}

		val, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			fmt.Printf("Произошла ошибка конвертации, попробуйте ещё раз\n\n")
			continue
		}

		switch {
		case val < 0:
			fmt.Printf("Значение не может быть отрицательным, попробуйте ещё раз\n\n")
		default:
			return val
		}
	}
}
