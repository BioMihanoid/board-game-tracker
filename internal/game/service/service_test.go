package service

import (
	"testing"

	"board-game-tracker/internal/game/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewGameService(t *testing.T) {
	t.Parallel()

	t.Run(
		"nil repository",
		func(t *testing.T) {
			t.Parallel()

			expected := &GameService{repo: nil}

			actual := NewGameService(nil)

			assert.Equal(t, expected, actual)
		},
	)

	t.Run(
		"real repository",
		func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewMockGameRepository(t)
			expected := &GameService{repo: repo}

			actual := NewGameService(repo)

			assert.Equal(t, expected, actual)
		},
	)
}
