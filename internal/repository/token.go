package repository

import (
	"HelloCity/internal/repository/cache"
	"context"
)

type TokenRepository interface {
	Set(ctx context.Context, prefix, key, value string) error
	Verify(ctx context.Context, prefix, key, value string) (bool, error)
	Del(ctx context.Context, prefix, key string) error
	Get(ctx context.Context, prefix, key string) (string, error)
}

type TokenCachedRepository struct {
	cache cache.TokenCache
}

func (t *TokenCachedRepository) Get(ctx context.Context, prefix, key string) (string, error) {
	return t.cache.Get(ctx, prefix, key)
}

func (t *TokenCachedRepository) Set(ctx context.Context, prefix, key, value string) error {
	return t.cache.Set(ctx, prefix, key, value)
}

func (t *TokenCachedRepository) Verify(ctx context.Context, prefix, key, value string) (bool, error) {
	return t.cache.Verify(ctx, prefix, key, value)
}

func (t *TokenCachedRepository) Del(ctx context.Context, prefix, key string) error {
	return t.cache.Del(ctx, prefix, key)
}

func NewTokenCachedRepository(cache cache.TokenCache) TokenRepository {
	return &TokenCachedRepository{
		cache: cache,
	}
}
