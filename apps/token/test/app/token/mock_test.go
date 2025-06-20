package token_test

import "github.com/stretchr/testify/mock"

// MockTokenGenerator implements domain.TokenGenerator for unit testing
type MockTokenGenerator struct {
	mock.Mock
}

func (m *MockTokenGenerator) GenerateToken(claims map[string]string) (string, int64, error) {
	args := m.Called(claims)
	if args.Get(0) == nil {
		return "", 0, args.Error(1)
	}
	return args.String(0), args.Get(1).(int64), args.Error(2)
}

func (m *MockTokenGenerator) VerifyToken(token string) (map[string]string, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]string), args.Error(1)
}
