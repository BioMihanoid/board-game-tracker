package memory

import "board-game-tracker/internal/model"

// seedGames возвращает новую копию стартовых данных,
// чтобы каждый репозиторий владел своим состоянием.
func seedGames() map[int64]*model.Game {
	return map[int64]*model.Game{
		1: {
			ID:          1,
			Name:        "Древний ужас",
			CostKopecks: 5999_00,
			Type:        model.MainGameType,
			Publisher:   "Hobby Games",
		},
		2: {
			ID:          2,
			BaseGameID:  1,
			Name:        "Древний ужас. Забытые тайны",
			CostKopecks: 2999_00,
			Type:        model.AdditionalContentGameType,
			Publisher:   "Hobby Games",
		},
		3: {
			ID:          3,
			BaseGameID:  1,
			Name:        "Древний ужас. Знамения каркозы",
			CostKopecks: 2999_00,
			Type:        model.AdditionalContentGameType,
			Publisher:   "Hobby Games",
		},
		4: {
			ID:          4,
			BaseGameID:  1,
			Name:        "Древний ужас. Таинственные руины",
			CostKopecks: 2999_00,
			Type:        model.AdditionalContentGameType,
			Publisher:   "Hobby Games",
		},
		5: {
			ID:          5,
			BaseGameID:  1,
			Name:        "Древний ужас. Разрушенные города",
			CostKopecks: 2999_00,
			Type:        model.AdditionalContentGameType,
			Publisher:   "Hobby Games",
		},
		6: {
			ID:          6,
			BaseGameID:  1,
			Name:        "Древний ужас. Мир грез",
			CostKopecks: 4999_00,
			Type:        model.AdditionalContentGameType,
			Publisher:   "Hobby Games",
		},
	}
}
