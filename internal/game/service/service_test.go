package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"board-game-tracker/internal/game/service/mocks"
	"board-game-tracker/internal/models"
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

func TestAddGame(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	game := &models.Game{
		Name:      "test",
		Publisher: "test publisher",
	}
	expectedID := int64(1)
	repo := mocks.NewMockGameRepository(t)
	repo.EXPECT().SaveGame(ctx, game).Return(expectedID, nil)
	service := NewGameService(repo)

	actual, err := service.AddGame(ctx, game)

	assert.Nil(t, err)
	assert.Equal(t, expectedID, actual)
}

func TestGetGameById(t *testing.T) {
	t.Parallel()

	t.Run(
		"no error",
		func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			game := &models.Game{
				Name:      "test",
				Publisher: "test publisher",
			}
			id := int64(1)
			repo := mocks.NewMockGameRepository(t)
			repo.EXPECT().GetGameById(ctx, id).Return(game, nil)
			service := NewGameService(repo)

			actual, err := service.GetGameById(ctx, id)

			assert.Nil(t, err)
			assert.Equal(t, game, actual)
		},
	)

	t.Run(
		"with error",
		func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			id := int64(1)
			repo := mocks.NewMockGameRepository(t)
			repo.EXPECT().GetGameById(ctx, id).Return(nil, fmt.Errorf("error"))
			service := NewGameService(repo)

			actual, err := service.GetGameById(ctx, id)

			assert.Nil(t, actual)
			assert.Error(t, err)
		},
	)
}

func TestGetAllGame(t *testing.T) {
	t.Parallel()

	t.Run(
		"no error",
		func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			games := map[int64]*models.Game{
				1: {
					Name:      "test",
					Publisher: "test publisher",
				},
			}
			repo := mocks.NewMockGameRepository(t)
			repo.EXPECT().GetAllGame(ctx).Return(games, nil)
			service := NewGameService(repo)

			actual, err := service.GetAllGame(ctx)

			assert.Nil(t, err)
			assert.Equal(t, games, actual)
		},
	)

	t.Run(
		"with error",
		func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			repo := mocks.NewMockGameRepository(t)
			repo.EXPECT().GetAllGame(ctx).Return(nil, fmt.Errorf("error"))
			service := NewGameService(repo)

			actual, err := service.GetAllGame(ctx)

			assert.Nil(t, actual)
			assert.Error(t, err)
		},
	)
}

func TestGetAdditionalContentByGameBaseID(t *testing.T) {
	t.Parallel()

	t.Run(
		"no error",
		func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			games := map[int64]*models.Game{
				1: {
					Name:      "test",
					Publisher: "test publisher",
				},
			}
			id := int64(1)
			repo := mocks.NewMockGameRepository(t)
			repo.EXPECT().GetAdditionalContentByGameBaseID(ctx, id).Return(games, nil)
			service := NewGameService(repo)

			actual, err := service.GetAdditionalContentByGameBaseID(ctx, id)

			assert.Nil(t, err)
			assert.Equal(t, games, actual)
		},
	)

	t.Run(
		"with error",
		func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			id := int64(1)
			repo := mocks.NewMockGameRepository(t)
			repo.EXPECT().GetAdditionalContentByGameBaseID(ctx, id).Return(nil, fmt.Errorf("error"))
			service := NewGameService(repo)

			actual, err := service.GetAdditionalContentByGameBaseID(ctx, id)

			assert.Nil(t, actual)
			assert.Error(t, err)
		},
	)
}
