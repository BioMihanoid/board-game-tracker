package memory

import (
	"board-game-tracker/internal/model"
	"board-game-tracker/internal/seed"
)

// NewGameRepository создаёт репозиторий в памяти со стартовой коллекцией.
// Данные живут только до выхода из программы — для хранения на диске
// используйте repository/file.
func NewGameRepository() *GameRepository {
	m := make(map[int64]*model.Game)

	for _, game := range seed.Games() {
		m[game.ID] = game
	}

	nextID := int64(1)
	for id := range m {
		if id >= nextID {
			nextID = id + 1
		}
	}

	return &GameRepository{
		m:      m,
		nextID: nextID,
	}
}
