package repository

import (
	"context"
	"encoding/json"
	"time"
	"younified-backend/contracts/user/model"
	"younified-backend/providers/database"

	"github.com/go-redis/redis/v8"
)

type RedisUserRepository struct {
	client *database.RedisClient
}

func NewRedisUserRepository(client *database.RedisClient) *RedisUserRepository {
	return &RedisUserRepository{client: client}
}

func (r *RedisUserRepository) CacheUser(ctx context.Context, key string, user *model.User) error {

	// Serialize union to JSON
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	// Cache for 1 hour
	return r.client.Set(ctx, key, data, 24*time.Hour)
}

func (r *RedisUserRepository) CacheUsers(ctx context.Context, key string, user []*model.User) error {

	// Serialize union to JSON
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	// Cache for 1 hour
	return r.client.Set(ctx, key, data, 24*time.Hour)
}

func (r *RedisUserRepository) GetUserFromCache(ctx context.Context, key string) (*model.User, error) {
	var user model.User
	// Retrieve from cache
	err := r.client.GetJSON(ctx, key, &user)
	if err == redis.Nil {
		return nil, nil // Cache miss
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *RedisUserRepository) GetUsersFromCache(ctx context.Context, key string) ([]*model.User, error) {
	var users []*model.User
	// Retrieve from cache
	err := r.client.GetJSON(ctx, key, &users)
	if err == redis.Nil {
		return nil, nil // Cache miss
	} else if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *RedisUserRepository) InvalidateCache(ctx context.Context, key string) error {
	err := r.client.Delete(ctx, key)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisUserRepository) CacheExists(ctx context.Context, key string) (bool, error) {
	return r.client.Exists(ctx, key)
}
