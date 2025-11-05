-- Create OTPs table for MFA and one-time codes
CREATE TABLE IF NOT EXISTS otps (
    key TEXT PRIMARY KEY,
    code TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for efficient cleanup and expiry queries
CREATE INDEX IF NOT EXISTS idx_otps_expires_at ON otps (expires_at);

-- Add comment for documentation
COMMENT ON TABLE otps IS 'Stores OTP/MFA codes (email, SMS, TOTP) with expiry.';
