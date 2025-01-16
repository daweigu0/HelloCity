package service

import (
	"HelloCity/internal/repository"
	"context"
)

type TokenService interface {
	Set(ctx context.Context, prefix, key, value string) error
	Verify(ctx context.Context, prefix, key, value string) (bool, error)
	Del(ctx context.Context, prefix, key string) error
	Get(ctx context.Context, prefix, key string) (string, error)
}

type tokenService struct {
	repo repository.TokenRepository
}

func NewTokenService(repo repository.TokenRepository) TokenService {
	return &tokenService{
		repo: repo,
	}
}
func (t *tokenService) Get(ctx context.Context, prefix, key string) (string, error) {
	return t.repo.Get(ctx, prefix, key)
}

func (t *tokenService) Set(ctx context.Context, prefix, key, value string) error {
	return t.repo.Set(ctx, prefix, key, value)
}

func (t *tokenService) Verify(ctx context.Context, prefix, key, value string) (bool, error) {
	return t.repo.Verify(ctx, prefix, key, value)
}

func (t *tokenService) Del(ctx context.Context, prefix, key string) error {
	return t.repo.Del(ctx, prefix, key)
}
