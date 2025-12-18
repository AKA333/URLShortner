package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/AKA333/URLShortner/internal/models"
	"github.com/go-redis/redis/v8"

)

type RedisRepository struct {
	client *redis.Client
	ctx context.Context	
}

func NewRedisRepository(addr string) (*RedisRepository, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		Password: "",
		DB: 0, // using default DB
	})
	
	ctx := context.Background()
	
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	
	return &RedisRepository{
		client: client,
		ctx: ctx,
	}, nil
}

func (r *RedisRepository) SetURL(url *models.URL, expiration time.Duration) error {
	data, err := json.Marshal(url)
	if err != nil {
		return fmt.Errorf("failed to marshal URL %w", err)
	}
	
	key := fmt.Sprintf("url:%s", url.ShortCode)
	err = r.client.Set(r.ctx, key, data, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set URL in redis %w", err)
	}
	
	return nil
}

func (r *RedisRepository) GetURL(shortCode string) (*models.URL, error) {
	key := fmt.Sprintf("url: %s", shortCode)
	
	data, err := r.client.Get(r.ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to get URL from Redis %w", err)
	}
	
	var url models.URL
	if err := json.Unmarshal(data, &url); err != nil {
		return nil, fmt.Errorf("failed to unmarshal URL %w", err)
	}
	
	return &url, nil
}

func (r *RedisRepository) DeleteURL(shortCode string) error {
	key := fmt.Sprintf("url: %s", shortCode)
	
	err := r.client.Del(r.ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete URl from Redis %w", err)
	}
	return nil
}

func (r *RedisRepository) Close() error {
	return r.client.Close()
}