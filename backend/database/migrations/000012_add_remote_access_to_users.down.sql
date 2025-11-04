-- Migration 000012: Remove Remote Access Control from Users
-- File: database/migrations/000012_add_remote_access_to_users.down.sql

DROP INDEX IF EXISTS idx_users_remote_access;

ALTER TABLE users DROP COLUMN IF EXISTS remote_access_approved_at;
ALTER TABLE users DROP COLUMN IF EXISTS remote_access_approved_by;
ALTER TABLE users DROP COLUMN IF EXISTS remote_access_reason;
ALTER TABLE users DROP COLUMN IF EXISTS allow_remote_access;
