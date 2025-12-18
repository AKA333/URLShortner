package models 

import (
	"time"
)

// URL represents a shortened URL entity.
type URL struct {
	ID int64 `json:"id"`
	ShortCode string `json:"short_code"`
	LongURL string `json:"long_url"`
	UserID *int64 `json:"user_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	Clicks int64 `json:"clicks"`
}

// Req to create a short URL
type ShortenRequest struct {
	LongURL string `json:"long_url" validate:"required,url"`
	CustomAlias string `json:"custom_alias,omitempty" validate:"omitempty,alphanum, max=10"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// Response after creating a short URL
type ShortenResponse struct {
	ShortCode string `json:"short_code"`
	ShortURL string `json:"short_url"`
	LongURL string `json:"long_url"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type Config struct {
	ServerPort string `env:"SERVER_PORT" envDefault:"8080"`
	RedisURL string `env:"REDIS_URL" envDefault:"redis://localhost:6379"`
	DBHost string `env:"DB_HOST" envDefault:"localhost"`
	DBPort string `env:"DB_PORT" envDefault:"5432"`
	DBName string `env:"DB_NAME" envDefault:"urlshortener"`
	DBUser string `env:"DB_USER" envDefault:"admin"`
	DBPassword string `env:"DB_PASSWORD" envDefault:"password"`
	BaseURL string `env:"BASE_URL" envDefault:"http://localhost:8080"`
}
