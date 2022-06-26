package repository

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrIndicatorNotFound = errors.New("indicator not existing in inmemry db")
)

type Repository interface {
	AddIndicator(ctx context.Context, indicator int64) error
	GetIndicator(ctx context.Context, indicator int64) (int64, error)
	RemoveIndicator(ctx context.Context, indicator int64)
}

type inMemoryDB struct {
	hashcashIndicators map[int64]time.Time
	rw                 *sync.RWMutex
}

// NewInMemoryDB ..
func NewInMemoryDB() Repository {
	return &inMemoryDB{hashcashIndicators: map[int64]time.Time{}, rw: &sync.RWMutex{}}
}

func (r *inMemoryDB) AddIndicator(ctx context.Context, indicator int64) error {
	r.rw.Lock()
	defer r.rw.Unlock()

	r.hashcashIndicators[indicator] = time.Now()

	return nil
}

// GetStamp ..
func (r *inMemoryDB) GetIndicator(ctx context.Context, requestedIndicator int64) (int64, error) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	_, ok := r.hashcashIndicators[requestedIndicator]
	if ok {
		return requestedIndicator, nil
	}

	return 0, ErrIndicatorNotFound
}

// UpdateStamp - removes indicator from db
func (r *inMemoryDB) RemoveIndicator(ctx context.Context, newIndicator int64) {
	r.rw.Lock()
	defer r.rw.Unlock()

	delete(r.hashcashIndicators, newIndicator)
}
