package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"board-game-tracker/internal/model"
)

type fakeRepo struct {
	games  map[int64]*model.Game
	nextID int64
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{
		games:  map[int64]*model.Game{},
		nextID: 1,
	}
}

func (r *fakeRepo) GetByID(_ context.Context, id int64) (*model.Game, error) {
	game, ok := r.games[id]
	if !ok {
		return nil, fmt.Errorf("get game %d: %w", id, model.ErrNotFound)
	}

	copied := *game

	return &copied, nil
}

func (r *fakeRepo) ListBaseGames(_ context.Context) ([]*model.Game, error) {
	res := make([]*model.Game, 0, len(r.games))
	for _, game := range r.games {
		if game.Type == model.MainGameType {
			copied := *game
			res = append(res, &copied)
		}
	}

	return res, nil
}

func (r *fakeRepo) ListExpansions(_ context.Context, baseID int64) ([]*model.Game, error) {
	res := make([]*model.Game, 0, len(r.games))
	for _, game := range r.games {
		if game.BaseGameID == baseID && game.Type != model.MainGameType {
			copied := *game
			res = append(res, &copied)
		}
	}

	return res, nil
}

func (r *fakeRepo) Save(_ context.Context, game *model.Game) (int64, error) {
	id := r.nextID
	r.nextID++

	stored := *game
	stored.ID = id
	r.games[id] = &stored

	return id, nil
}

func addBaseGame(t *testing.T, service *GameService) int64 {
	t.Helper()

	id, err := service.AddGame(context.Background(), &model.Game{
		Name:        "Древний ужас",
		Type:        model.MainGameType,
		Publisher:   "Hobby Games",
		CostKopecks: 5999_00,
	})
	if err != nil {
		t.Fatalf("добавить основную игру: %v", err)
	}

	return id
}

func TestAddGameMain(t *testing.T) {
	service := NewGameService(newFakeRepo())
	ctx := context.Background()

	id := addBaseGame(t, service)

	game, err := service.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("получить игру: %v", err)
	}

	if game.ID != id {
		t.Errorf("ID = %d, ожидался %d", game.ID, id)
	}
	if game.CostKopecks != 5999_00 {
		t.Errorf("стоимость = %d, ожидалась %d", game.CostKopecks, 5999_00)
	}
}

func TestAddExpansionComposesName(t *testing.T) {
	service := NewGameService(newFakeRepo())
	ctx := context.Background()

	baseID := addBaseGame(t, service)

	id, err := service.AddGame(ctx, &model.Game{
		Name:        "Забытые тайны",
		Type:        model.AdditionalContentGameType,
		BaseGameID:  baseID,
		CostKopecks: 2999_00,
	})
	if err != nil {
		t.Fatalf("добавить дополнение: %v", err)
	}

	game, err := service.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("получить дополнение: %v", err)
	}

	if want := "Древний ужас. Забытые тайны"; game.Name != want {
		t.Errorf("название = %q, ожидалось %q", game.Name, want)
	}
	if game.Publisher != "Hobby Games" {
		t.Errorf("издатель = %q, ожидался %q", game.Publisher, "Hobby Games")
	}
	if game.BaseGameID != baseID {
		t.Errorf("BaseGameID = %d, ожидался %d", game.BaseGameID, baseID)
	}
}

func TestAddExpansionWithoutBase(t *testing.T) {
	service := NewGameService(newFakeRepo())

	_, err := service.AddGame(context.Background(), &model.Game{
		Name: "Забытые тайны",
		Type: model.AdditionalContentGameType,
	})

	if !errors.Is(err, model.ErrInvalidGame) {
		t.Fatalf("ожидалась ErrInvalidGame, получено %v", err)
	}
}

func TestAddExpansionWithExpansionAsBase(t *testing.T) {
	service := NewGameService(newFakeRepo())
	ctx := context.Background()

	baseID := addBaseGame(t, service)
	expansionID, err := service.AddGame(ctx, &model.Game{
		Name:       "Забытые тайны",
		Type:       model.AdditionalContentGameType,
		BaseGameID: baseID,
	})
	if err != nil {
		t.Fatalf("добавить дополнение: %v", err)
	}

	_, err = service.AddGame(ctx, &model.Game{
		Name:       "Ещё одно",
		Type:       model.AdditionalContentGameType,
		BaseGameID: expansionID,
	})

	if !errors.Is(err, model.ErrNotBaseGame) {
		t.Fatalf("ожидалась ErrNotBaseGame, получено %v", err)
	}
}

func TestAddGameInvalid(t *testing.T) {
	tests := []struct {
		name string
		game *model.Game
	}{
		{
			name: "пустое название",
			game: &model.Game{Name: "   ", Type: model.MainGameType},
		},
		{
			name: "отрицательная стоимость",
			game: &model.Game{Name: "Игра", Type: model.MainGameType, CostKopecks: -1},
		},
		{
			name: "неизвестный тип",
			game: &model.Game{Name: "Игра", Type: model.GameType("Что-то")},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := NewGameService(newFakeRepo())

			_, err := service.AddGame(context.Background(), test.game)
			if !errors.Is(err, model.ErrInvalidGame) {
				t.Fatalf("ожидалась ErrInvalidGame, получено %v", err)
			}
		})
	}
}

func TestListExpansionsEmptyIsNotError(t *testing.T) {
	service := NewGameService(newFakeRepo())
	ctx := context.Background()

	baseID := addBaseGame(t, service)

	expansions, err := service.ListExpansions(ctx, baseID)
	if err != nil {
		t.Fatalf("получить дополнения: %v", err)
	}
	if len(expansions) != 0 {
		t.Fatalf("ожидался пустой список, получено %d", len(expansions))
	}
}

func TestListExpansionsForExpansionID(t *testing.T) {
	service := NewGameService(newFakeRepo())
	ctx := context.Background()

	baseID := addBaseGame(t, service)
	expansionID, err := service.AddGame(ctx, &model.Game{
		Name:       "Забытые тайны",
		Type:       model.AdditionalContentGameType,
		BaseGameID: baseID,
	})
	if err != nil {
		t.Fatalf("добавить дополнение: %v", err)
	}

	_, err = service.ListExpansions(ctx, expansionID)
	if !errors.Is(err, model.ErrNotBaseGame) {
		t.Fatalf("ожидалась ErrNotBaseGame, получено %v", err)
	}
}
