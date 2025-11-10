CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    balance INTEGER DEFAULT 0 CHECK (balance >= 0),
    referrer_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    refresh_token VARCHAR(500)
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_referrer_id ON users(referrer_id);
CREATE INDEX idx_users_refresh_token ON users(refresh_token);
CREATE INDEX idx_users_balance ON users(balance DESC);