package file

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"board-game-tracker/internal/model"
)

func newRepo(t *testing.T) (*GameRepository, string) {
	t.Helper()

	path := filepath.Join(t.TempDir(), "games.json")

	repo, err := New(path)
	if err != nil {
		t.Fatalf("открыть репозиторий: %v", err)
	}

	return repo, path
}

func TestNewCreatesFileWithSeed(t *testing.T) {
	repo, path := newRepo(t)
	ctx := context.Background()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("прочитать файл: %v", err)
	}

	var stored snapshot
	if err := json.Unmarshal(data, &stored); err != nil {
		t.Fatalf("разобрать файл: %v", err)
	}

	if len(stored.Games) != 6 {
		t.Fatalf("в файле %d игр, ожидалось 6", len(stored.Games))
	}
	if stored.NextID != 7 {
		t.Fatalf("next_id = %d, ожидалось 7", stored.NextID)
	}
	if stored.Games[0].ID != 1 {
		t.Fatalf("игры не отсортированы по id: первая %d", stored.Games[0].ID)
	}

	game, err := repo.GetByID(ctx, 1)
	if err != nil {
		t.Fatalf("получить игру: %v", err)
	}
	if game.Name != "Древний ужас" {
		t.Fatalf("название = %q", game.Name)
	}
}

func TestMoneyIsStoredAsIntegerKopecks(t *testing.T) {
	_, path := newRepo(t)

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("прочитать файл: %v", err)
	}

	want := `"cost_kopecks": 599900`
	if !strings.Contains(string(data), want) {
		t.Fatalf("в файле нет %q", want)
	}
}

func TestSavePersistsAfterReopen(t *testing.T) {
	repo, path := newRepo(t)
	ctx := context.Background()

	id, err := repo.Save(ctx, &model.Game{
		Name:        "Каркассон",
		Type:        model.MainGameType,
		Publisher:   "Hobby Games",
		CostKopecks: 2499_50,
	})
	if err != nil {
		t.Fatalf("сохранить игру: %v", err)
	}

	reopened, err := New(path)
	if err != nil {
		t.Fatalf("открыть файл заново: %v", err)
	}

	game, err := reopened.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("игра не сохранилась: %v", err)
	}

	if game.Name != "Каркассон" {
		t.Errorf("название = %q", game.Name)
	}
	if game.CostKopecks != 2499_50 {
		t.Errorf("стоимость = %d, ожидалась %d", game.CostKopecks, 2499_50)
	}
}

func TestIDsAreNotReusedAfterReopen(t *testing.T) {
	repo, path := newRepo(t)
	ctx := context.Background()

	first, err := repo.Save(ctx, &model.Game{Name: "Первая", Type: model.MainGameType})
	if err != nil {
		t.Fatalf("сохранить игру: %v", err)
	}

	reopened, err := New(path)
	if err != nil {
		t.Fatalf("открыть файл заново: %v", err)
	}

	second, err := reopened.Save(ctx, &model.Game{Name: "Вторая", Type: model.MainGameType})
	if err != nil {
		t.Fatalf("сохранить игру: %v", err)
	}

	if first == second {
		t.Fatalf("ID переиспользован: %d", first)
	}
	if second != first+1 {
		t.Fatalf("второй ID = %d, ожидался %d", second, first+1)
	}
}

func TestSaveRollsBackWhenWriteFails(t *testing.T) {
	repo, path := newRepo(t)
	ctx := context.Background()

	// Подменяем файл каталогом: rename временного файла на каталог не сработает.
	if err := os.Remove(path); err != nil {
		t.Fatalf("удалить файл: %v", err)
	}
	if err := os.Mkdir(path, 0o755); err != nil {
		t.Fatalf("создать каталог: %v", err)
	}

	id, err := repo.Save(ctx, &model.Game{Name: "Не сохранится", Type: model.MainGameType})
	if err == nil {
		t.Fatalf("ожидалась ошибка записи, игра сохранена с id %d", id)
	}

	if _, err := repo.GetByID(ctx, 7); !errors.Is(err, model.ErrNotFound) {
		t.Fatalf("состояние в памяти не откатилось: %v", err)
	}
}

func TestNewFailsOnBrokenJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "games.json")

	if err := os.WriteFile(path, []byte("это не json"), 0o644); err != nil {
		t.Fatalf("подготовить файл: %v", err)
	}

	if _, err := New(path); err == nil {
		t.Fatalf("ожидалась ошибка разбора файла")
	}
}

func TestNewFailsOnDuplicateIDs(t *testing.T) {
	path := filepath.Join(t.TempDir(), "games.json")

	content := `{"next_id": 3, "games": [
		{"id": 1, "name": "Первая", "type": "Основная игра", "cost_kopecks": 0},
		{"id": 1, "name": "Дубль", "type": "Основная игра", "cost_kopecks": 0}
	]}`

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("подготовить файл: %v", err)
	}

	if _, err := New(path); err == nil {
		t.Fatalf("ожидалась ошибка про повторяющийся id")
	}
}

func TestNextIDIgnoresTooSmallValueInFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "games.json")

	// next_id меньше максимального id — как будто файл правили руками.
	content := `{"next_id": 1, "games": [
		{"id": 1, "name": "Первая", "type": "Основная игра", "cost_kopecks": 0},
		{"id": 5, "name": "Пятая", "type": "Основная игра", "cost_kopecks": 0}
	]}`

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("подготовить файл: %v", err)
	}

	repo, err := New(path)
	if err != nil {
		t.Fatalf("открыть файл: %v", err)
	}

	id, err := repo.Save(context.Background(), &model.Game{Name: "Новая", Type: model.MainGameType})
	if err != nil {
		t.Fatalf("сохранить игру: %v", err)
	}

	if id <= 5 {
		t.Fatalf("id = %d, ожидался больше 5", id)
	}
}

func TestCreatesMissingDirectories(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "data", "games.json")

	if _, err := New(path); err != nil {
		t.Fatalf("открыть репозиторий: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("файл не создан: %v", err)
	}
}
