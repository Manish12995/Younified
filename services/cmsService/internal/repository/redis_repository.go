package repository

import "younified-backend/providers/database"

type RedisUserRepository struct {
	client *database.RedisClient
}

func NewRedisUserRepository(client *database.RedisClient) *RedisUserRepository {
	return &RedisUserRepository{client: client}
}