CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    date BIGINT NOT NULL,
    category TEXT NOT NULL DEFAULT '',
    type TEXT NOT NULL DEFAULT '',
    amount BIGINT NOT NULL DEFAULT 0,
    description TEXT NOT NULL DEFAULT '',
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);