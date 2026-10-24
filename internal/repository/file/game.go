// Package file хранит коллекцию в JSON-файле.
//
// Данные читаются в память при старте и целиком перезаписываются при каждом
// изменении: запись идёт во временный файл с последующим rename, поэтому
// прерванный процесс не оставит наполовину записанный JSON.
package file

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"board-game-tracker/internal/model"
	"board-game-tracker/internal/seed"
)

type snapshot struct {
	NextID int64         `json:"next_id"`
	Games  []*model.Game `json:"games"`
}

type GameRepository struct {
	path   string
	m      map[int64]*model.Game
	nextID int64
	rw     sync.RWMutex
}

// New открывает файл с данными. Если файла нет — создаёт его
// и заполняет стартовой коллекцией. Битый файл не перезаписывается:
// возвращается ошибка, чтобы данные можно было починить руками.
func New(path string) (*GameRepository, error) {
	repo := &GameRepository{
		path: path,
		m:    make(map[int64]*model.Game),
	}

	if err := repo.load(); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}

		repo.fillSeed()

		if err := repo.flush(); err != nil {
			return nil, err
		}

		return repo, nil
	}

	return repo, nil
}

func (repo *GameRepository) load() error {
	data, err := os.ReadFile(repo.path)
	if err != nil {
		return fmt.Errorf("открыть файл %s: %w", repo.path, err)
	}

	var stored snapshot
	if err := json.Unmarshal(data, &stored); err != nil {
		return fmt.Errorf("разобрать файл %s: %w", repo.path, err)
	}

	for _, game := range stored.Games {
		if game == nil {
			continue
		}

		if game.ID <= 0 {
			return fmt.Errorf("файл %s: у игры %q некорректный id %d", repo.path, game.Name, game.ID)
		}

		if _, ok := repo.m[game.ID]; ok {
			return fmt.Errorf("файл %s: повторяющийся id %d", repo.path, game.ID)
		}

		repo.m[game.ID] = game
	}

	// Считаем заново, а не доверяем файлу на слово: его правят руками.
	repo.nextID = nextIDAfter(repo.m)
	if stored.NextID > repo.nextID {
		repo.nextID = stored.NextID
	}

	return nil
}

func (repo *GameRepository) fillSeed() {
	for _, game := range seed.Games() {
		repo.m[game.ID] = game
	}

	repo.nextID = nextIDAfter(repo.m)
}

// flush перезаписывает файл текущим состоянием. Вызывается под lock.
func (repo *GameRepository) flush() error {
	games := make([]*model.Game, 0, len(repo.m))
	for _, game := range repo.m {
		copied := *game
		games = append(games, &copied)
	}

	sort.Slice(games, func(i, j int) bool {
		return games[i].ID < games[j].ID
	})

	data, err := json.MarshalIndent(snapshot{NextID: repo.nextID, Games: games}, "", "  ")
	if err != nil {
		return fmt.Errorf("сериализовать данные: %w", err)
	}

	if err := writeAtomic(repo.path, append(data, '\n')); err != nil {
		return fmt.Errorf("записать файл %s: %w", repo.path, err)
	}

	return nil
}

func (repo *GameRepository) GetByID(_ context.Context, id int64) (*model.Game, error) {
	repo.rw.RLock()
	defer repo.rw.RUnlock()

	game, ok := repo.m[id]
	if !ok {
		return nil, fmt.Errorf("get game by id %d: %w", id, model.ErrNotFound)
	}

	res := *game

	return &res, nil
}

func (repo *GameRepository) ListBaseGames(_ context.Context) ([]*model.Game, error) {
	repo.rw.RLock()
	defer repo.rw.RUnlock()

	return repo.filter(func(game *model.Game) bool {
		return game.Type == model.MainGameType
	}), nil
}

func (repo *GameRepository) ListExpansions(_ context.Context, baseID int64) ([]*model.Game, error) {
	repo.rw.RLock()
	defer repo.rw.RUnlock()

	return repo.filter(func(game *model.Game) bool {
		return game.BaseGameID == baseID && game.Type != model.MainGameType
	}), nil
}

// filter вызывается под RLock, поэтому копирует сам.
func (repo *GameRepository) filter(match func(*model.Game) bool) []*model.Game {
	res := make([]*model.Game, 0, len(repo.m))
	for _, game := range repo.m {
		if !match(game) {
			continue
		}

		copied := *game
		res = append(res, &copied)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].ID < res[j].ID
	})

	return res
}

func (repo *GameRepository) Save(_ context.Context, game *model.Game) (int64, error) {
	if game == nil {
		return 0, fmt.Errorf("%w: пустая игра", model.ErrInvalidGame)
	}

	repo.rw.Lock()
	defer repo.rw.Unlock()

	id := repo.nextID

	stored := *game
	stored.ID = id
	repo.m[id] = &stored
	repo.nextID = id + 1

	// Если запись на диск не удалась — откатываем состояние в памяти,
	// чтобы оно не разошлось с файлом.
	if err := repo.flush(); err != nil {
		delete(repo.m, id)
		repo.nextID = id

		return 0, err
	}

	game.ID = id

	return id, nil
}

func nextIDAfter(games map[int64]*model.Game) int64 {
	var maxID int64

	for id := range games {
		if id > maxID {
			maxID = id
		}
	}

	return maxID + 1
}

// writeAtomic пишет файл через временный файл и rename:
// читатель никогда не увидит недописанный JSON.
func writeAtomic(path string, data []byte) (err error) {
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("создать каталог %s: %w", dir, err)
	}

	tmp, err := os.CreateTemp(dir, ".games-*.json")
	if err != nil {
		return fmt.Errorf("создать временный файл: %w", err)
	}

	tmpName := tmp.Name()

	defer func() {
		if err != nil {
			os.Remove(tmpName)
		}
	}()

	if _, err = tmp.Write(data); err != nil {
		tmp.Close()
		return fmt.Errorf("записать временный файл: %w", err)
	}

	if err = tmp.Sync(); err != nil {
		tmp.Close()
		return fmt.Errorf("сбросить на диск: %w", err)
	}

	if err = tmp.Close(); err != nil {
		return fmt.Errorf("закрыть временный файл: %w", err)
	}

	if err = os.Chmod(tmpName, 0o644); err != nil {
		return fmt.Errorf("выставить права: %w", err)
	}

	if err = os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("переименовать временный файл: %w", err)
	}

	return nil
}
