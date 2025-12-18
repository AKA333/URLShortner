-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    api_key_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create URLs table
CREATE TABLE IF NOT EXISTS urls (
    id SERIAL PRIMARY KEY,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    long_url TEXT NOT NULL,
    user_id INT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP,
    clicks INT DEFAULT 0,
    CONSTRAINT short_code_format CHECK (short_code ~ '^[a-zA-Z0-9]+$')
);

-- Create analytics table 
CREATE TABLE IF NOT EXISTS analytics (
    id SERIAL PRIMARY KEY,
    short_code VARCHAR(10) NOT NULL REFERENCES urls(short_code) ON DELETE CASCADE,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ip_address INET,
    user_agent TEXT,
    referrer TEXT,
    country_code VARCHAR(2)
);

-- Create indexes for performance
CREATE INDEX idx_urls_short_code ON urls(short_code);
CREATE INDEX idx_analytics_short_code ON analytics(short_code);
CREATE INDEX idx_analytics_timestamp ON analytics(timestamp);
CREATE INDEX idx_urls_user_id ON urls(user_id);

-- Insert a test user (optional, for development)
INSERT INTO users (email, api_key_hash) 
VALUES ('test@example.com', '$2a$10$YourTestHashHere')
ON CONFLICT (email) DO NOTHING;