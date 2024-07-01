ALTER TABLE sessions 
    ADD COLUMN user_id bigint; 

ALTER TABLE sessions 
    DROP CONSTRAINT sessions_email_fkey;

ALTER TABLE sessions 
    ADD CONSTRAINT sessions_user_id_fkey
    FOREIGN KEY (user_id) references users (id);

UPDATE sessions
SET user_id = u.id 
FROM sessions AS s
JOIN users AS u
ON s.email = u.email;

ALTER TABLE sessions
    DROP COLUMN email;

ALTER TABLE sessions
    ADD COLUMN email VARCHAR(30);
