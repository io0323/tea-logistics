-- +migrate Down
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP INDEX IF EXISTS idx_password_resets_token;
DROP INDEX IF EXISTS idx_access_tokens_user_id;
DROP INDEX IF EXISTS idx_access_tokens_token;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_username;
DROP TABLE IF EXISTS password_resets;
DROP TABLE IF EXISTS access_tokens;
DROP TABLE IF EXISTS users; 