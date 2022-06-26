package repository_test

import (
	"context"
	"sitemapGenerator/WordOfWisdom/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddIndicator(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryDB()

	// Act
	err := repo.AddIndicator(context.Background(), 123456789)

	// Assert
	assert.Nil(t, err)
}

func TestGetIndicator(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryDB()

	// Act
	err := repo.AddIndicator(context.Background(), 123456789)
	assert.NoError(t, err)

	indicator, err := repo.GetIndicator(context.Background(), 123456789)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, int64(123456789), indicator)
}

func TestRemoveIndicator(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryDB()

	// Act
	err := repo.AddIndicator(context.Background(), 123456789)
	assert.NoError(t, err)
	repo.RemoveIndicator(context.Background(), 123456789)
	_, err = repo.GetIndicator(context.Background(), 123456789)

	// Assert
	assert.NotNil(t, err)
}
