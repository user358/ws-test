package balance

import (
	"sync"

	"github.com/shopspring/decimal"
)

// InMemoryStore simple thread safe store
type InMemoryStore struct {
	m  map[int64]decimal.Decimal
	rw sync.RWMutex
}

// NewInMemoryStore constructor
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		m: map[int64]decimal.Decimal{},
	}
}

// Get returns stored value or default 0.
func (r *InMemoryStore) Get(id int64) decimal.Decimal {
	r.rw.RLock()
	defer r.rw.RUnlock()

	return r.m[id]
}

// Update fn is used for safety work with mutex. Takes current value and should return next value or error.
func (r *InMemoryStore) Update(id int64, fn func(v decimal.Decimal) (decimal.Decimal, error)) error {
	r.rw.Lock()
	defer r.rw.Unlock()

	v := r.m[id]
	v, err := fn(v)
	if err != nil {
		return err
	}

	r.m[id] = v

	return nil
}
