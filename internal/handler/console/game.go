package console

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"board-game-tracker/internal/model"
)

var ErrInputClosed = errors.New("ввод завершён")

const maxKopecks = 1_000_000_000_00

type GameService interface {
	GetByID(ctx context.Context, id int64) (*model.Game, error)
	ListBaseGames(ctx context.Context) ([]*model.Game, error)
	ListExpansions(ctx context.Context, baseID int64) ([]*model.Game, error)
	AddGame(ctx context.Context, game *model.Game) (int64, error)
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

func (h *Handler) PrintGameByID(ctx context.Context) error {
	game, err := h.askGame(ctx, "id")
	if err != nil {
		return err
	}

	printGame(game)

	return nil
}

func (h *Handler) PrintBaseGames(ctx context.Context) error {
	games, err := h.gameService.ListBaseGames(ctx)
	if err != nil {
		return fmt.Errorf("получить игры: %w", err)
	}

	if len(games) == 0 {
		fmt.Printf("Коллекция пуста\n\n")
		return nil
	}

	for _, game := range games {
		printGame(game)
	}

	return nil
}

func (h *Handler) PrintExpansions(ctx context.Context) error {
	base, err := h.askBaseGame(ctx)
	if err != nil {
		return err
	}

	games, err := h.gameService.ListExpansions(ctx, base.ID)
	if err != nil {
		return fmt.Errorf("получить дополнения игры %q: %w", base.Name, err)
	}

	if len(games) == 0 {
		fmt.Printf("У игры %q нет дополнений\n\n", base.Name)
		return nil
	}

	fmt.Printf("Дополнения игры %q:\n\n", base.Name)
	for _, game := range games {
		printGame(game)
	}

	return nil
}

func (h *Handler) AddGame(ctx context.Context) (int64, error) {
	gameType, err := h.askGameType()
	if err != nil {
		return 0, err
	}

	game := &model.Game{Type: gameType}

	if gameType != model.MainGameType {
		base, err := h.askBaseGame(ctx)
		if err != nil {
			return 0, err
		}

		tag := "дополнения"
		if gameType == model.PromoGameType {
			tag = "промо набора"
		}

		game.BaseGameID = base.ID
		if game.Name, err = h.askString(fmt.Sprintf("название %s", tag)); err != nil {
			return 0, err
		}
		if game.CostKopecks, err = h.askKopecks(fmt.Sprintf("стоимость %s", tag)); err != nil {
			return 0, err
		}

		return h.gameService.AddGame(ctx, game)
	}

	if game.Name, err = h.askString("название"); err != nil {
		return 0, err
	}
	if game.CostKopecks, err = h.askKopecks("стоимость"); err != nil {
		return 0, err
	}
	if game.Publisher, err = h.askString("издателя"); err != nil {
		return 0, err
	}

	return h.gameService.AddGame(ctx, game)
}

func printGame(game *model.Game) {
	fmt.Printf("ID: %d\nИгра: %s\nТип: %s\nСтоимость: %s ₽\nИздатель: %s\n",
		game.ID,
		game.Name,
		game.Type,
		model.FormatKopecks(game.CostKopecks),
		game.Publisher,
	)

	if game.BaseGameID != 0 {
		fmt.Printf("ID базовой игры: %d\n", game.BaseGameID)
	}

	fmt.Println()
}

func (h *Handler) askGameType() (model.GameType, error) {
	types := []model.GameType{
		model.MainGameType,
		model.AdditionalContentGameType,
		model.PromoGameType,
	}

	for {
		fmt.Printf("Доступные типы игры:\n")
		for i, gameType := range types {
			fmt.Printf("%d. %s\n", i+1, gameType)
		}

		choice, err := h.askInt64("тип")
		if err != nil {
			return "", err
		}

		if choice >= 1 && choice <= int64(len(types)) {
			return types[choice-1], nil
		}

		fmt.Printf("Выбран несуществующий тип игры, попробуйте ещё раз\n\n")
	}
}

func (h *Handler) askGame(ctx context.Context, label string) (*model.Game, error) {
	for {
		id, err := h.askInt64(label)
		if err != nil {
			return nil, err
		}

		game, err := h.gameService.GetByID(ctx, id)
		switch {
		case errors.Is(err, model.ErrNotFound):
			fmt.Printf("Игра с таким id отсутствует, попробуйте ещё раз\n\n")
		case err != nil:
			return nil, fmt.Errorf("получить игру %d: %w", id, err)
		default:
			return game, nil
		}
	}
}

func (h *Handler) askBaseGame(ctx context.Context) (*model.Game, error) {
	for {
		game, err := h.askGame(ctx, "id базовой игры")
		if err != nil {
			return nil, err
		}

		if game.Type != model.MainGameType {
			fmt.Printf("Игра %q не является основной, выберите основную игру\n\n", game.Name)
			continue
		}

		return game, nil
	}
}

func (h *Handler) scanLine() (string, error) {
	if !h.scanner.Scan() {
		if err := h.scanner.Err(); err != nil {
			return "", fmt.Errorf("прочитать ввод: %w", err)
		}

		return "", ErrInputClosed
	}

	return strings.TrimSpace(h.scanner.Text()), nil
}

func (h *Handler) askString(label string) (string, error) {
	for {
		fmt.Printf("Введите %s игры: ", label)

		val, err := h.scanLine()
		if err != nil {
			return "", err
		}

		if val == "" {
			fmt.Printf("Пустой ввод, попробуйте ещё раз\n\n")
			continue
		}

		return val, nil
	}
}

func (h *Handler) askInt64(label string) (int64, error) {
	for {
		fmt.Printf("Введите %s игры: ", label)

		val, err := h.scanLine()
		if err != nil {
			return 0, err
		}

		if val == "" {
			fmt.Printf("Пустой ввод, попробуйте ещё раз\n\n")
			continue
		}

		number, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			fmt.Printf("Некорректное число, попробуйте ещё раз\n\n")
			continue
		}

		if number <= 0 {
			fmt.Printf("Значение должно быть больше нуля, попробуйте ещё раз\n\n")
			continue
		}

		return number, nil
	}
}

func (h *Handler) askKopecks(label string) (int64, error) {
	for {
		fmt.Printf("Введите %s игры (например 5999 или 5999,99): ", label)

		val, err := h.scanLine()
		if err != nil {
			return 0, err
		}

		kopecks, err := parseKopecks(val)
		if err != nil {
			fmt.Printf("%v, попробуйте ещё раз\n\n", err)
			continue
		}

		return kopecks, nil
	}
}

func parseKopecks(input string) (int64, error) {
	normalized := strings.ReplaceAll(strings.TrimSpace(input), " ", "")
	normalized = strings.ReplaceAll(normalized, ",", ".")

	if normalized == "" {
		return 0, errors.New("пустой ввод")
	}

	if dot := strings.IndexByte(normalized, '.'); dot >= 0 && len(normalized)-dot-1 > 2 {
		return 0, errors.New("не больше двух знаков после запятой")
	}

	value, err := strconv.ParseFloat(normalized, 64)
	if err != nil {
		return 0, errors.New("некорректная сумма")
	}

	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, errors.New("некорректная сумма")
	}

	if value < 0 {
		return 0, errors.New("стоимость не может быть отрицательной")
	}

	if value > maxKopecks/100 {
		return 0, errors.New("слишком большая сумма")
	}

	return int64(math.Round(value * 100)), nil
}
