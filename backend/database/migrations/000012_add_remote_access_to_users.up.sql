-- Migration 000012: Add Remote Access Control to Users
-- File: database/migrations/000012_add_remote_access_to_users.up.sql

-- Add remote access toggle to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS allow_remote_access BOOLEAN DEFAULT FALSE;

-- Add reason/notes field for audit
ALTER TABLE users ADD COLUMN IF NOT EXISTS remote_access_reason TEXT;

-- Add approved by and date for tracking
ALTER TABLE users ADD COLUMN IF NOT EXISTS remote_access_approved_by UUID REFERENCES users(id);
ALTER TABLE users ADD COLUMN IF NOT EXISTS remote_access_approved_at TIMESTAMP;

-- Create index for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_remote_access ON users(allow_remote_access);

-- Example data:
-- allow_remote_access: true = Can access from anywhere (home, travel, etc.)
-- allow_remote_access: false = Office IP restriction applies (default)
-- remote_access_reason: 'CEO - Frequent travel', 'Remote employee', 'Field worker'
