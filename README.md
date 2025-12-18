# URLShortner
Production grade URL Shortening Service (Like Bitly)

Tech Stack: Golang, Redis (for cache), PostgreSQL, REST API, GCP/AWS (free tier for deployment), Docker.

- With rate limiting per IP/API key
- Detailed analytics dashboard (count clicks)

Project Structure:
.
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── handlers/            # HTTP request handlers
│   ├── service/             # Business logic
│   ├── repository/          # Database interaction logic (Postgres/Redis)
│   └── models/              # Data structures
├── pkg/
│   └── utils/               # Shared utilities (e.g., ID generation)
└── tests/                   # Test suites

Core Shortening logic:
1. POST /api/v1/shorten
2. GET /:short_code

