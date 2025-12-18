// package tests

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"github.com/AKA333/URLShortner/internal/models"
// 	"github.com/AKA333/URLShortner/internal/repository"
// 	"github.com/AKA333/URLShortner/internal/service"
// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/alicebob/miniredis/v2"
// 	"github.com/go-redis/redis/v8"
// 	"github.com/stretchr/testify/assert"
// )

// func TestURLService_ShortenURL(t *testing.T) {
// 	// Setup mock PostgreSQL
// 	db, mock, err := sqlmock.New()
// 	assert.NoError(t, err)
// 	defer db.Close()

// 	postgresRepo := &repository.PostgresRepository{db: db}

// 	// Setup mock Redis
// 	mr, err := miniredis.Run()
// 	assert.NoError(t, err)
// 	defer mr.Close()

// 	redisClient := redis.NewClient(&redis.Options{
// 		Addr: mr.Addr(),
// 	})
// 	redisRepo := &repository.RedisRepository{
// 		client: redisClient,
// 		ctx:    context.Background(),
// 	}

// 	// Create service
// 	urlService := service.NewURLService(postgresRepo, redisRepo, "http://localhost:8080")

// 	tests := []struct {
// 		name          string
// 		req           *models.ShortenRequest
// 		setupMocks    func()
// 		expectError   bool
// 		errorContains string
// 	}{
// 		{
// 			name: "Valid URL with custom alias",
// 			req: &models.ShortenRequest{
// 				LongURL:     "https://example.com",
// 				CustomAlias: "test123",
// 			},
// 			setupMocks: func() {
// 				// Mock short code check
// 				rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
// 				mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM urls WHERE short_code = \$1\)`).
// 					WithArgs("test123").
// 					WillReturnRows(rows)

// 				// Mock URL creation
// 				rows = sqlmock.NewRows([]string{"id", "created_at"}).
// 					AddRow(1, time.Now())
// 				mock.ExpectQuery(`INSERT INTO urls`).
// 					WillReturnRows(rows)
// 			},
// 			expectError: false,
// 		},
// 		{
// 			name: "Duplicate custom alias",
// 			req: &models.ShortenRequest{
// 				LongURL:     "https://example.com",
// 				CustomAlias: "duplicate",
// 			},
// 			setupMocks: func() {
// 				rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
// 				mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM urls WHERE short_code = \$1\)`).
// 					WithArgs("duplicate").
// 					WillReturnRows(rows)
// 			},
// 			expectError:   true,
// 			errorContains: "already in use",
// 		},
// 		{
// 			name: "Invalid URL",
// 			req: &models.ShortenRequest{
// 				LongURL: "not-a-valid-url",
// 			},
// 			setupMocks:    func() {},
// 			expectError:   true,
// 			errorContains: "invalid URL",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.setupMocks()

// 			response, err := urlService.ShortenURL(context.Background(), tt.req)

// 			if tt.expectError {
// 				assert.Error(t, err)
// 				if tt.errorContains != "" {
// 					assert.Contains(t, err.Error(), tt.errorContains)
// 				}
// 			} else {
// 				assert.NoError(t, err)
// 				assert.NotNil(t, response)
// 				assert.Equal(t, tt.req.CustomAlias, response.ShortCode)
// 				assert.Equal(t, tt.req.LongURL, response.LongURL)
// 			}

// 			// Ensure all expectations were met
// 			assert.NoError(t, mock.ExpectationsWereMet())
// 		})
// 	}
// }

// func TestURLService_RedirectURL(t *testing.T) {
// 	// Setup similar to above test
// 	db, mock, err := sqlmock.New()
// 	assert.NoError(t, err)
// 	defer db.Close()

// 	postgresRepo := &repository.PostgresRepository{db: db}

// 	mr, err := miniredis.Run()
// 	assert.NoError(t, err)
// 	defer mr.Close()

// 	redisClient := redis.NewClient(&redis.Options{
// 		Addr: mr.Addr(),
// 	})
// 	redisRepo := &repository.RedisRepository{
// 		client: redisClient,
// 		ctx:    context.Background(),
// 	}

// 	urlService := service.NewURLService(postgresRepo, redisRepo, "http://localhost:8080")

// 	t.Run("Redirect to existing URL", func(t *testing.T) {
// 		// Mock database response
// 		createdAt := time.Now()
// 		rows := sqlmock.NewRows([]string{"id", "short_code", "long_url", "user_id", "created_at", "expires_at", "clicks"}).
// 			AddRow(1, "abc123", "https://example.com", nil, createdAt, nil, 0)

// 		mock.ExpectQuery(`SELECT id, short_code, long_url, user_id, created_at, expires_at, clicks FROM urls`).
// 			WithArgs("abc123").
// 			WillReturnRows(rows)

// 		// Mock increment clicks
// 		mock.ExpectExec(`UPDATE urls SET clicks = clicks \+ 1 WHERE short_code = \$1`).
// 			WithArgs("abc123").
// 			WillReturnResult(sqlmock.NewResult(0, 1))

// 		longURL, err := urlService.RedirectURL(context.Background(), "abc123")

// 		assert.NoError(t, err)
// 		assert.Equal(t, "https://example.com", longURL)
// 		assert.NoError(t, mock.ExpectationsWereMet())
// 	})

// 	t.Run("Non-existent URL", func(t *testing.T) {
// 		mock.ExpectQuery(`SELECT id, short_code, long_url, user_id, created_at, expires_at, clicks FROM urls`).
// 			WithArgs("nonexistent").
// 			WillReturnRows(sqlmock.NewRows([]string{}))

// 		longURL, err := urlService.RedirectURL(context.Background(), "nonexistent")

// 		assert.Error(t, err)
// 		assert.Equal(t, "", longURL)
// 		assert.Contains(t, err.Error(), "not found")
// 		assert.NoError(t, mock.ExpectationsWereMet())
// 	})
// }
