ALTER TABLE users ADD COLUMN refresh_token varchar unique;
ALTER TABLE users ADD COLUMN exp_refresh_token integer;