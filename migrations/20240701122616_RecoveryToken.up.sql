CREATE TABLE recovery_tokens (
    email VARCHAR,
    token VARCHAR(128) NOT NULL UNIQUE CHECK (length(token) >= 64),
    created_at integer NOT NULL,
    FOREIGN KEY (email) REFERENCES users(email) ON DELETE CASCADE
);