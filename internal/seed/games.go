// Package seed содержит стартовую коллекцию игр,
// которая создаётся при первом запуске, пока файл с данными ещё не существует.
package seed

import "board-game-tracker/internal/model"

// Games возвращает стартовую коллекцию. Каждый вызов создаёт новые объекты.
func Games() []*model.Game {
	return []*model.Game{
		{
			ID:          1,
			Name:        "Древний ужас",
			CostKopecks: 5999_00,
			Type:        model.MainGameType,
			Publisher:   "Hobby Games",
		},
		{
			ID:          2,
			BaseGameID:  1,
			Name:        "Древний ужас. Забытые тайны",
			CostKopecks: 2999_00,
			Type:        model.AdditionalContentGameType,
			Publisher:   "Hobby Games",
		},
		{
			ID:          3,
			BaseGameID:  1,
			Name:        "Древний ужас. Знамения каркозы",
			CostKopecks: 2999_00,
			Type:        model.AdditionalContentGameType,
			Publisher:   "Hobby Games",
		},
		{
			ID:          4,
			BaseGameID:  1,
			Name:        "Древний ужас. Таинственные руины",
			CostKopecks: 2999_00,
			Type:        model.AdditionalContentGameType,
			Publisher:   "Hobby Games",
		},
		{
			ID:          5,
			BaseGameID:  1,
			Name:        "Древний ужас. Разрушенные города",
			CostKopecks: 2999_00,
			Type:        model.AdditionalContentGameType,
			Publisher:   "Hobby Games",
		},
		{
			ID:          6,
			BaseGameID:  1,
			Name:        "Древний ужас. Мир грез",
			CostKopecks: 4999_00,
			Type:        model.AdditionalContentGameType,
			Publisher:   "Hobby Games",
		},
	}
}
