package balance

import "github.com/shopspring/decimal"

type Store interface {
	Get(id int64) decimal.Decimal
	Update(id int64, fn func(v decimal.Decimal) (decimal.Decimal, error)) error
}

type Repository struct {
	db Store
}

func NewRepository(db Store) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Get(userID int64) decimal.Decimal {
	return r.db.Get(userID)
}

func (r *Repository) Add(userID int64, value decimal.Decimal) error {
	err := r.db.Update(userID, func(v decimal.Decimal) (decimal.Decimal, error) {
		return v.Add(value), nil
	})

	return err
}

func (r *Repository) Sub(userID int64, value decimal.Decimal) error {
	err := r.db.Update(userID, func(v decimal.Decimal) (decimal.Decimal, error) {
		next := v.Sub(value)
		if next.LessThan(decimal.Zero) {
			return v, ErrNotEnoughFunds
		}

		return next, nil
	})

	return err
}
