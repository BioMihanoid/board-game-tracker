package memory

import "board-game-tracker/internal/model"

var gamesMap = map[int64]*model.Game{
	1: {
		ID:        1,
		Name:      "Древний ужас",
		Cost:      5999,
		Type:      model.MainGameType,
		Publisher: "Hobby Games",
	},
	2: {
		ID:         2,
		BaseGameID: 1,
		Name:       "Древний ужас. Забытые тайны",
		Cost:       2999,
		Type:       model.AdditionalContentGameType,
		Publisher:  "Hobby Games",
	},
	3: {
		ID:         3,
		BaseGameID: 1,
		Name:       "Древний ужас. Знамения каркозы",
		Cost:       2999,
		Type:       model.AdditionalContentGameType,
		Publisher:  "Hobby Games",
	},
	4: {
		ID:         4,
		BaseGameID: 1,
		Name:       "Древний ужас. Таинственные руины",
		Cost:       2999,
		Type:       model.AdditionalContentGameType,
		Publisher:  "Hobby Games",
	},
	5: {
		ID:         5,
		BaseGameID: 1,
		Name:       "Древний ужас. Разрушенные города",
		Cost:       2999,
		Type:       model.AdditionalContentGameType,
		Publisher:  "Hobby Games",
	},
	6: {
		ID:         6,
		BaseGameID: 1,
		Name:       "Древний ужас. Мир грез",
		Cost:       4999,
		Type:       model.AdditionalContentGameType,
		Publisher:  "Hobby Games",
	},
}
