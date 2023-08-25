package balance

import "errors"

var (
	ErrInvalidValue   = errors.New("balance: invalid value")
	ErrNotEnoughFunds = errors.New("balance: not enough funds")
)
