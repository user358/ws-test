package balance

import "github.com/shopspring/decimal"

type Repo interface {
	Get(userID int64) decimal.Decimal
	Add(userID int64, value decimal.Decimal) error
	Sub(userID int64, value decimal.Decimal) error
}

type Service struct {
	repo Repo
}

func NewService(repo Repo) *Service {
	return &Service{repo: repo}
}

func (s Service) GetByUserID(userID int64) decimal.Decimal {
	return s.repo.Get(userID)
}

func (s Service) Deposit(userID int64, value decimal.Decimal) error {
	if value.LessThan(decimal.Zero) {
		return ErrInvalidValue
	}

	return s.repo.Add(userID, value)
}

func (s Service) Withdraw(userID int64, value decimal.Decimal) error {
	if value.LessThan(decimal.Zero) {
		return ErrInvalidValue
	}

	return s.repo.Sub(userID, value)
}
