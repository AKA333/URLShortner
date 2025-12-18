package service

import (
	"context"
	"fmt"
	"time"

	"github.com/AKA333/URLShortner/internal/models"
	"github.com/AKA333/URLShortner/internal/repository"
	"github.com/AKA333/URLShortner/pkg/utils"
)

type URLService struct {
	postgresRepo *repository.PostgresRepository
	redisRepo    *repository.RedisRepository
	baseURL      string
}

func NewURLService(postgresRepo *repository.PostgresRepository, redisRepo *repository.RedisRepository, baseURL string) *URLService {
	return &URLService{
		postgresRepo: postgresRepo,
		redisRepo:    redisRepo,
		baseURL:      baseURL,
	}
}

func (s *URLService) ShortenURL(ctx context.Context, req *models.ShortenRequest) (*models.ShortenResponse, error) {
	// Validate long URL
	if !utils.ValidateURL(req.LongURL) {
		return nil, fmt.Errorf("invalid long URL")
	}

	var shortCode string

	// Check if custom alias is provided
	if req.CustomAlias != "" {
		// Validate custom alias
		if !utils.ValidateShortCode(req.CustomAlias) {
			return nil, fmt.Errorf("Custom alias must be 3-10 alphanumeric characters")
		}

		// Check if custom alias already exists
		exists, err := s.postgresRepo.ShortCodeExists(req.CustomAlias)
		if err != nil {
			return nil, fmt.Errorf("failed to check custom alias: %w", err)
		}

		if exists {
			return nil, fmt.Errorf("custom alias already in use")
		}
		shortCode = req.CustomAlias
	} else {
		// Generate unqiue short code
		for i := 0; i < 5; i++ { // Trying 5 times to generate unique code
			code, err := utils.GenerateShortCode(6)
			if err != nil {
				return nil, fmt.Errorf("failed to generate short code %w", err)
			}
			exists, err := s.postgresRepo.ShortCodeExists(code)
			if err != nil {
				return nil, fmt.Errorf("failed to check short code %w", err)
			}

			if !exists {
				shortCode = code
				break
			}
		}
		if shortCode == "" {
			return nil, fmt.Errorf("failed to generate unique short code after multiple attempts")
		}
	}

	// Create URL record
	url := &models.URL{
		ShortCode: shortCode,
		LongURL:   req.LongURL,
		UserID:    nil, // will update while adding authentication
		ExpiresAt: req.ExpiresAt,
		Clicks:    0,
	}

	// Saving to postgres
	if err := s.postgresRepo.CreateURL(url); err != nil {
		return nil, fmt.Errorf("failed to save URL %w", &err)
	}

	// Cache in Redis with ttl - 24hrs
	if s.redisRepo != nil {
		expiration := 24 * time.Hour
		if req.ExpiresAt != nil && req.ExpiresAt.Before(time.Now().Add(expiration)) {
			expiration = time.Until(*req.ExpiresAt)
		}

		if err := s.redisRepo.SetURL(url, expiration); err != nil {
			fmt.Printf("Warning: failed to cache URL in redis %v\n", err)
		}
	}

	response := &models.ShortenResponse{
		ShortCode: shortCode,
		ShortURL:  fmt.Sprintf("%s%s", s.baseURL, shortCode),
		LongURL:   req.LongURL,
		CreatedAt: url.CreatedAt,
		ExpiresAt: req.ExpiresAt,
	}
	return response, nil
}

func (s *URLService) RedirectURL(ctx context.Context, shortCode string) (string, error) {
	var url *models.URL
	var err error

	if s.redisRepo != nil {
		url, err = s.redisRepo.GetURL(shortCode)
		if err != nil {
			return "", fmt.Errorf("failed to get URL from cache %w", err)
		}
	}
	if url == nil {
		url, err = s.postgresRepo.GetURLByShortCode(shortCode)
		if err != nil {
			return "", fmt.Errorf("failed to get URL from database %w", err)
		}

		if url == nil {
			return "", fmt.Errorf("short URL not found or expired")
		}

		if s.redisRepo != nil {
			expiration := 24 * time.Hour
			if url.ExpiresAt != nil {
				expiration = time.Until(*url.ExpiresAt)
				if expiration < 0 {
					expiration = 0
				}
			}
			if expiration > 0 {
				if err := s.redisRepo.SetURL(url, expiration); err != nil {
					fmt.Printf("Warning: failed to cache URL in redis %v\n", err)
				}
			}
		}
	}

	go func() {
		if err := s.postgresRepo.IncrementClicks(shortCode); err != nil {
			fmt.Printf("Warning: failed to increment clicks %v\n", err)
		}
	}()

	return url.LongURL, nil
}

func (s *URLService) GetURLStats(ctx context.Context, shortCode string) (*models.URL, error) {
	url, err := s.postgresRepo.GetURLByShortCode(shortCode)

	if err != nil {
		return nil, fmt.Errorf("failed to get URL stats %w", err)
	}

	if url == nil {
		return nil, fmt.Errorf("URL not found")
	}

	return url, nil
}
