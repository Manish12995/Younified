package repository

import (
	"context"
	"encoding/json"
	"time"
	"younified-backend/contracts/union/model"
	"younified-backend/providers/database"

	"github.com/go-redis/redis/v8"
)

type RedisUnionRepository struct {
	client *database.RedisClient
}

func NewRedisUnionRepository(client *database.RedisClient) *RedisUnionRepository {
	return &RedisUnionRepository{client: client}
}

func (r *RedisUnionRepository) CacheUnion(ctx context.Context, key string, union *model.Union) error {

	// Serialize union to JSON
	data, err := json.Marshal(union)
	if err != nil {
		return err
	}

	// Cache for 1 hour
	return r.client.Set(ctx, key, data, 24*time.Hour)
}

func (r *RedisUnionRepository) CacheUnions(ctx context.Context, key string, union []*model.Union) error {

	// Serialize union to JSON
	data, err := json.Marshal(union)
	if err != nil {
		return err
	}

	// Cache for 1 hour
	return r.client.Set(ctx, key, data, 24*time.Hour)
}

func (r *RedisUnionRepository) GetUnionFromCache(ctx context.Context, key string) (*model.Union, error) {
	var union model.Union
	// Retrieve from cache
	err := r.client.GetJSON(ctx, key, &union)
	if err == redis.Nil {
		return nil, nil // Cache miss
	} else if err != nil {
		return nil, err
	}

	return &union, nil
}

func (r *RedisUnionRepository) GetUnionsFromCache(ctx context.Context, key string) ([]*model.Union, error) {
	var unions []*model.Union
	// Retrieve from cache
	err := r.client.GetJSON(ctx, key, &unions)
	if err == redis.Nil {
		return nil, nil // Cache miss
	} else if err != nil {
		return nil, err
	}

	return unions, nil
}

func (r *RedisUnionRepository) InvalidateCache(ctx context.Context, key string) error {
	err := r.client.Delete(ctx, key)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisUnionRepository) CacheExists(ctx context.Context, key string) (bool, error) {
	return r.client.Exists(ctx, key)
}
