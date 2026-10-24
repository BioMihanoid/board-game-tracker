package memory

import (
	"context"
	"errors"
	"testing"

	"board-game-tracker/internal/model"
)

func TestSaveUsesUniqueIDs(t *testing.T) {
	repo := NewGameRepository()
	ctx := context.Background()

	first, err := repo.Save(ctx, &model.Game{Name: "Новая", Type: model.MainGameType})
	if err != nil {
		t.Fatalf("сохранить игру: %v", err)
	}

	second, err := repo.Save(ctx, &model.Game{Name: "Ещё одна", Type: model.MainGameType})
	if err != nil {
		t.Fatalf("сохранить игру: %v", err)
	}

	if first == second {
		t.Fatalf("ID совпали: %d", first)
	}

	// Стартовые данные не должны затираться.
	if _, err := repo.GetByID(ctx, 6); err != nil {
		t.Fatalf("стартовая игра с id 6 потерялась: %v", err)
	}

	games, err := repo.ListBaseGames(ctx)
	if err != nil {
		t.Fatalf("получить игры: %v", err)
	}
	if len(games) != 3 {
		t.Fatalf("ожидалось 3 основные игры, получено %d", len(games))
	}
}

func TestRepositoriesAreIndependent(t *testing.T) {
	ctx := context.Background()

	first := NewGameRepository()
	second := NewGameRepository()

	if _, err := first.Save(ctx, &model.Game{Name: "Новая", Type: model.MainGameType}); err != nil {
		t.Fatalf("сохранить игру: %v", err)
	}

	games, err := second.ListBaseGames(ctx)
	if err != nil {
		t.Fatalf("получить игры: %v", err)
	}
	if len(games) != 1 {
		t.Fatalf("репозитории делят состояние: во втором %d игр вместо 1", len(games))
	}
}

func TestGetByIDNotFound(t *testing.T) {
	repo := NewGameRepository()

	_, err := repo.GetByID(context.Background(), 999)
	if !errors.Is(err, model.ErrNotFound) {
		t.Fatalf("ожидалась ErrNotFound, получено %v", err)
	}
}

func TestGetByIDReturnsCopy(t *testing.T) {
	repo := NewGameRepository()
	ctx := context.Background()

	game, err := repo.GetByID(ctx, 1)
	if err != nil {
		t.Fatalf("получить игру: %v", err)
	}

	game.Name = "Изменено снаружи"

	stored, err := repo.GetByID(ctx, 1)
	if err != nil {
		t.Fatalf("получить игру: %v", err)
	}
	if stored.Name == "Изменено снаружи" {
		t.Fatalf("репозиторий отдаёт указатель на внутреннее состояние")
	}
}

func TestListExpansionsSorted(t *testing.T) {
	repo := NewGameRepository()

	expansions, err := repo.ListExpansions(context.Background(), 1)
	if err != nil {
		t.Fatalf("получить дополнения: %v", err)
	}
	if len(expansions) != 5 {
		t.Fatalf("ожидалось 5 дополнений, получено %d", len(expansions))
	}

	for i := 1; i < len(expansions); i++ {
		if expansions[i-1].ID >= expansions[i].ID {
			t.Fatalf("дополнения не отсортированы по ID")
		}
	}
}
