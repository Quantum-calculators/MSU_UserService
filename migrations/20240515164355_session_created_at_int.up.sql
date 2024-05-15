DROP TABLE sessions;

CREATE TABLE sessions (
    id SERIAL NOT NULL PRIMARY KEY,
    user_id bigint REFERENCES users(id) ON DELETE CASCADE,
    refresh_token varchar NOT NULL,
    fingerprint varchar NOT NULL,
    expires_in bigint NOT NULL,
    created_at bigint NOT NULL
);