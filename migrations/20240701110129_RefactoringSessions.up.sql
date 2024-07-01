ALTER TABLE sessions
    ADD COLUMN email VARCHAR(30);

ALTER TABLE sessions 
    DROP CONSTRAINT sessions_user_id_fkey;

ALTER TABLE sessions 
    ADD CONSTRAINT sessions_email_fkey
    FOREIGN KEY (email) references users (email);

UPDATE sessions
SET email = u.email 
FROM sessions AS s
    JOIN users AS u
    ON s.user_id = u.id;

ALTER TABLE sessions
    DROP COLUMN user_id;