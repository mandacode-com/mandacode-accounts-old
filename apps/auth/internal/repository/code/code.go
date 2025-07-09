package coderepo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"github.com/redis/go-redis/v9"
	codedomain "mandacode.com/accounts/auth/internal/domain/repository/code"
	"mandacode.com/accounts/auth/internal/util"
)

type loginCodeManager struct {
	codeGen   *util.RandomGenerator
	codeTTL   time.Duration
	codeStore *redis.Client
	prefix    string
}

func (l *loginCodeManager) IssueCode(ctx context.Context, userID uuid.UUID) (string, error) {
	code, err := l.codeGen.GenerateSecureRandomCode()
	if err != nil {
		return "", err
	}

	key := l.prefix + code

	err = l.codeStore.Set(ctx, key, userID, l.codeTTL).Err()
	if err != nil {
		return "", err
	}

	return code, nil
}

func (l *loginCodeManager) ValidateCode(ctx context.Context, userID uuid.UUID, code string) (bool, error) {
	key := l.prefix + code
	storedUserID, err := l.codeStore.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil // Code does not exist
		}
		return false, errors.New(err.Error(), "Failed to get login code from store", errcode.ErrInternalFailure)
	}

	if storedUserID != userID.String() {
		return false, nil // Code exists but does not match user ID
	}

	// Delete the code after successful validation
	err = l.codeStore.Del(ctx, code).Err()
	if err != nil {
		return false, errors.New(err.Error(), "Failed to delete login code from store", errcode.ErrInternalFailure)
	}

	return true, nil // Code is valid and deleted
}

func NewCodeManager(codeGen *util.RandomGenerator, codeTTL time.Duration, codeStore *redis.Client, prefix string) codedomain.CodeManager {
	return &loginCodeManager{
		codeGen:   codeGen,
		codeTTL:   codeTTL,
		codeStore: codeStore,
		prefix:    prefix,
	}
}
