package balance

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryStore_Get(t *testing.T) {
	db := NewInMemoryStore()
	db.m = map[int64]decimal.Decimal{
		1: decimal.NewFromInt(100),
	}

	assert.True(t, db.Get(1).Equal(decimal.NewFromInt(100)))
	assert.True(t, db.Get(2).Equal(decimal.Zero))
}

func TestInMemoryStore_Update(t *testing.T) {
	db := NewInMemoryStore()
	assert.True(t, db.Get(1).Equal(decimal.Zero))

	err := db.Update(1, func(v decimal.Decimal) (decimal.Decimal, error) {
		return decimal.NewFromInt(100), nil
	})
	assert.Nil(t, err)
	assert.True(t, db.Get(1).Equal(decimal.NewFromInt(100)))
}
