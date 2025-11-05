-- Drop the OTPs table and related index
DROP INDEX IF EXISTS idx_otps_expires_at;
DROP TABLE IF EXISTS otps;
