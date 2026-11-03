package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGameService(t *testing.T) {
	t.Parallel()

	expected := &GameService{repo: nil}

	actual := NewGameService(nil)

	assert.Equal(t, expected, actual)
}
