package balance

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository used for test service business logic
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Get(userID int64) decimal.Decimal {
	args := m.Called(userID)
	return args.Get(0).(decimal.Decimal)
}

func (m *MockRepository) Add(userID int64, value decimal.Decimal) error {
	args := m.Called(userID, value)
	return args.Error(0)
}

func (m *MockRepository) Sub(userID int64, value decimal.Decimal) error {
	args := m.Called(userID, value)
	return args.Error(0)
}

func TestService_GetByUserID(t *testing.T) {
	mockRepo := &MockRepository{}
	mockRepo.On("Get", int64(1)).Return(decimal.Zero)
	s := NewService(mockRepo)

	assert.Equal(t, s.GetByUserID(1), decimal.Zero)
}
