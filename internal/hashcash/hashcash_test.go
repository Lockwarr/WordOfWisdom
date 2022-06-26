package hashcash_test

import (
	"context"
	"crypto/sha1"
	"fmt"
	"sitemapGenerator/WordOfWisdom/internal/hashcash"
	"sitemapGenerator/WordOfWisdom/internal/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToString(t *testing.T) {
	// Arrange
	stamp := hashcash.Stamp{
		Version:    1,
		ZerosCount: 5,
		Date:       1546300800,
		Resource:   "http://example.com",
		Rand:       "123456789",
		Counter:    1}

	// Act
	result := stamp.ToString()

	// Assert
	assert.Equal(t, "1:5:1546300800:http://example.com::123456789:1", result)
}

func TestIsHashSolved(t *testing.T) {
	// Arrange
	stamp := hashcash.Stamp{
		Version:    1,
		ZerosCount: 5,
		Date:       int64(1656246214),
		Resource:   "test",
		Rand:       "123456789",
		Counter:    808598,
	}
	// Act
	solved := stamp.IsHashSolved()

	// Assert
	assert.Equal(t, true, solved)
}

func TestHashToStamp(t *testing.T) {
	// Arrange
	tests := []struct {
		name          string
		hashString    string
		expectedStamp hashcash.Stamp
		wantedErr     error
	}{
		{
			name:       "valid url",
			hashString: "1:5:1546300800:testMail@asd.ss::123456789:1",
			expectedStamp: hashcash.Stamp{
				Version:    1,
				ZerosCount: 5,
				Date:       int64(1546300800),
				Resource:   "testMail@asd.ss",
				Rand:       "123456789",
				Counter:    1,
			},
			wantedErr: nil,
		},
		{
			name:          "hash with invalid version",
			hashString:    ":5:1546300800:testMail@asd.ss::123456789:1:55",
			expectedStamp: hashcash.Stamp{},
			wantedErr:     fmt.Errorf("invalid version"),
		},
		{
			name:          "hash with invalid zeros count",
			hashString:    "1::1546300800:testMail@asd.ss::123456789:1",
			expectedStamp: hashcash.Stamp{},
			wantedErr:     fmt.Errorf("invalid zeros count"),
		},
		{
			name:          "hash with invalid date",
			hashString:    "1:5::testMail@asd.ss::123456789:1",
			expectedStamp: hashcash.Stamp{},
			wantedErr:     fmt.Errorf("invalid date"),
		},
		{
			name:          "invalid hash",
			hashString:    "5",
			expectedStamp: hashcash.Stamp{},
			wantedErr:     fmt.Errorf("invalid hash"),
		},
		{
			name:          "hash with invalid counter",
			hashString:    "1:5:1546300800:testMail@asd.ss::123456789:",
			expectedStamp: hashcash.Stamp{},
			wantedErr:     fmt.Errorf("invalid counter"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			stamp, err := hashcash.HashToStamp(tt.hashString)

			// Assert
			assert.Equal(t, tt.wantedErr, err)
			assert.Equal(t, tt.expectedStamp, stamp)
		})
	}
}

func TestValidStamp(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryDB()
	_ = repo.AddIndicator(context.Background(), 123456789)
	tests := []struct {
		name          string
		stamp         hashcash.Stamp
		expectedValid bool
	}{
		// As of June 26th, 2022, this test is valid, but it will be invalid after 28 days pass.
		{
			name: "valid solved stamp",
			stamp: hashcash.Stamp{
				Version:    1,
				ZerosCount: 5,
				Date:       int64(1656246214),
				Resource:   "test",
				Rand:       "123456789",
				Counter:    808598,
			},
			expectedValid: true,
		},
		{
			name: "unsolved stamp",
			stamp: hashcash.Stamp{
				Version:    1,
				ZerosCount: 5,
				Date:       time.Now().Unix(),
				Resource:   "test",
				Rand:       "123456789",
				Counter:    1,
			},
			expectedValid: false,
		},
		{
			name: "expired stamp",
			stamp: hashcash.Stamp{
				Version:    1,
				ZerosCount: 5,
				Date:       time.Now().Add(-time.Hour * 24 * 30).Unix(),
				Resource:   "test",
				Rand:       "123456789",
				Counter:    1,
			},
			expectedValid: false,
		},
		{
			name: "futuristic stamp",
			stamp: hashcash.Stamp{
				Version:    1,
				ZerosCount: 5,
				Date:       time.Now().Add(time.Hour * 24 * 5).Unix(),
				Resource:   "test",
				Rand:       "123456789",
				Counter:    1,
			},
			expectedValid: false,
		},
		{
			name: "invalid rand in stamp",
			stamp: hashcash.Stamp{
				Version:    1,
				ZerosCount: 5,
				Date:       time.Now().Unix(),
				Resource:   "test",
				Rand:       "invalid",
				Counter:    1,
			},
			expectedValid: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			valid := tt.stamp.ValidStamp(context.Background(), tt.stamp, repo)

			// Assert
			assert.Equal(t, tt.expectedValid, valid)
		})
	}
}

func TestComputeHashcash(t *testing.T) {
	// Arrange
	tests := []struct {
		name          string
		stamp         hashcash.Stamp
		maxIterations int
		wantedErr     error
	}{
		{
			name: "success scenario",
			stamp: hashcash.Stamp{
				Version:    1,
				ZerosCount: 5,
				Date:       time.Now().Unix(),
				Resource:   "test",
				Rand:       "123456789",
				Counter:    1,
			},
			// I read that the chances of getting the right hash are ~ one in a million with 20 wanted zeros
			maxIterations: 10000000,
			wantedErr:     nil,
		},
		{
			name: "max iterations reached",
			stamp: hashcash.Stamp{
				Version:    1,
				ZerosCount: 5,
				Date:       time.Now().Unix(),
				Resource:   "test",
				Rand:       "123456789",
				Counter:    1,
			},
			maxIterations: 1,
			wantedErr:     fmt.Errorf("max iterations exceeded"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			_, err := tt.stamp.ComputeHashcash(tt.maxIterations)

			// Assert
			assert.Equal(t, tt.wantedErr, err)
		})
	}
}

// sha1Hash - calculates sha1 hash from given string
func sha1Hash(data string) string {
	h := sha1.New()
	h.Write([]byte(data))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}
