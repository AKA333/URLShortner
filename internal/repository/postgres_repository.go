package repository

import (
	"context"
	"fmt"
	"time"
	"database/sql"

	"github.com/AKA333/URLShortner/internal/models"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(connStr string) (*PostgresRepository, error){
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database %w", err)
	}
	 ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err:= db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database %w", err)
	}
	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(time.Minute * 5)
	
	return &PostgresRepository{db: db}, nil
}

func (r * PostgresRepository) CreateURL (url *models.URL) error {
	query := `
			INSERT INTO urls (short_code, long_url, user_id, expires_at, clicks)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, created_at
		`
	err := r.db.QueryRow(
		query,
		url.ShortCode,
		url.LongURL,
		url.UserID,
		url.ExpiresAt,
		url.Clicks,
	).Scan(&url.ID, &url.CreatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create url %w", err)
	}
	return nil
}


func (r *PostgresRepository) GetURLByShortCode(shortCode string) (*models.URL, error) {
	query := `
			SELECT id, short_code, long_url, user_id, created_at, expires_at, clicks
			FROM urls
			WHERE short_code = $1
			AND (expires_at IS NULL OR expires_at > NOW())
		`
	url := &models.URL{}
	err := r.db.QueryRow(query, shortCode).Scan(
		&url.ID,
		&url.ShortCode,
		&url.LongURL,
		&url.UserID,
		&url.CreatedAt,
		&url.ExpiresAt,
		&url.Clicks,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get url by short code %w", err)
	}
	return url, nil
}

func (r * PostgresRepository)  IncrementClicks(shortCode string) error {
	query := `
			UPDATE urls
			SET clicks = clicks + 1
			WHERE short_code = $1
		`
	_, err := r.db.Exec(query, shortCode)
	if err != nil {
		return fmt.Errorf("failed to increment clicks %w", err)
	}
	return nil
}

func (r *PostgresRepository) ShortCodeExists(shortCode string) (bool, error) {
	query := `
			SELECT EXISTS (
				SELECT 1
				FROM urls
				WHERE short_code = $1
			)
		`
	var exists bool
	err := r.db.QueryRow(query, shortCode).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if short code exists %w", err)
	}
	return exists, nil
}

func (r *PostgresRepository) Close() error {
	return r.db.Close()
}
