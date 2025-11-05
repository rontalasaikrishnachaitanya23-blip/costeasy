-- Migration 000015: Update Organizations Table Structure

-- Add new columns (one at a time to handle existing columns gracefully)
DO $$ 
BEGIN
    -- Add code column
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='organizations' AND column_name='code') THEN
        ALTER TABLE organizations ADD COLUMN code VARCHAR(50);
    END IF;

    -- Add type column
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='organizations' AND column_name='type') THEN
        ALTER TABLE organizations ADD COLUMN type VARCHAR(50);
    END IF;

    -- Add emirate column
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='organizations' AND column_name='emirate') THEN
        ALTER TABLE organizations ADD COLUMN emirate VARCHAR(50);
    END IF;

    -- Add area column
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='organizations' AND column_name='area') THEN
        ALTER TABLE organizations ADD COLUMN area VARCHAR(100);
    END IF;

    -- Add website column
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='organizations' AND column_name='website') THEN
        ALTER TABLE organizations ADD COLUMN website VARCHAR(255);
    END IF;

    -- Add currency column
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='organizations' AND column_name='currency') THEN
        ALTER TABLE organizations ADD COLUMN currency VARCHAR(3);
    END IF;

    -- Add license_number column
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='organizations' AND column_name='license_number') THEN
        ALTER TABLE organizations ADD COLUMN license_number VARCHAR(100);
    END IF;

    -- Add license_expiry column
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='organizations' AND column_name='license_expiry') THEN
        ALTER TABLE organizations ADD COLUMN license_expiry TIMESTAMP;
    END IF;

    -- Add establishment_id column
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='organizations' AND column_name='establishment_id') THEN
        ALTER TABLE organizations ADD COLUMN establishment_id VARCHAR(100);
    END IF;

    -- Add description column
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='organizations' AND column_name='description') THEN
        ALTER TABLE organizations ADD COLUMN description TEXT;
    END IF;

    -- Add mfa_enabled column
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='organizations' AND column_name='mfa_enabled') THEN
        ALTER TABLE organizations ADD COLUMN mfa_enabled BOOLEAN DEFAULT FALSE;
    END IF;

    -- Add mfa_enforced column
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='organizations' AND column_name='mfa_enforced') THEN
        ALTER TABLE organizations ADD COLUMN mfa_enforced BOOLEAN DEFAULT FALSE;
    END IF;

    -- Add mfa_method column
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='organizations' AND column_name='mfa_method') THEN
        ALTER TABLE organizations ADD COLUMN mfa_method VARCHAR(20) DEFAULT 'none';
    END IF;

    -- Add allowed_mfa_methods column
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='organizations' AND column_name='allowed_mfa_methods') THEN
        ALTER TABLE organizations ADD COLUMN allowed_mfa_methods TEXT[];
    END IF;
END $$;

-- Update existing rows with default values
UPDATE organizations 
SET 
    type = COALESCE(type, 'OTHER'),
    currency = COALESCE(currency, 'AED'),
    license_number = COALESCE(license_number, 'PENDING'),
    establishment_id = COALESCE(establishment_id, 'PENDING'),
    mfa_method = COALESCE(mfa_method, 'none')
WHERE type IS NULL OR currency IS NULL OR license_number IS NULL OR establishment_id IS NULL;

-- Make required fields NOT NULL
ALTER TABLE organizations 
    ALTER COLUMN type SET NOT NULL,
    ALTER COLUMN currency SET NOT NULL,
    ALTER COLUMN license_number SET NOT NULL,
    ALTER COLUMN establishment_id SET NOT NULL;

-- Add new indexes
CREATE INDEX IF NOT EXISTS idx_organizations_type ON organizations(type);
CREATE INDEX IF NOT EXISTS idx_organizations_emirate ON organizations(emirate) WHERE emirate IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_organizations_code ON organizations(code) WHERE code IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_organizations_currency ON organizations(currency);

-- Add constraints (drop first if exists to avoid errors)
DO $$ 
BEGIN
    -- Add type constraint
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_organizations_type') THEN
        ALTER TABLE organizations 
            ADD CONSTRAINT chk_organizations_type 
            CHECK (type IN ('HEALTHCARE', 'RETAIL', 'MANUFACTURING', 'FINANCE', 
                            'EDUCATION', 'HOSPITALITY', 'LOGISTICS', 'REAL_ESTATE', 
                            'SERVICE', 'OTHER'));
    END IF;

    -- Add emirate constraint
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_organizations_emirate') THEN
        ALTER TABLE organizations 
            ADD CONSTRAINT chk_organizations_emirate 
            CHECK (emirate IS NULL OR emirate IN ('ABU_DHABI', 'DUBAI', 'SHARJAH', 
                                                   'RAS_AL_KHAIMAH', 'UMM_AL_QUWAIN', 
                                                   'FUJAIRAH', 'AJMAN'));
    END IF;

    -- Add mfa_method constraint
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_organizations_mfa_method') THEN
        ALTER TABLE organizations 
            ADD CONSTRAINT chk_organizations_mfa_method 
            CHECK (mfa_method IN ('none', 'sms', 'email', 'totp'));
    END IF;
END $$;

-- Add table and column comments
COMMENT ON TABLE organizations IS 'Business organizations with support for multiple industries';
COMMENT ON COLUMN organizations.type IS 'Organization type: HEALTHCARE, RETAIL, MANUFACTURING, etc.';
COMMENT ON COLUMN organizations.establishment_id IS 'Shafafiya provider code (required for healthcare organizations)';
COMMENT ON COLUMN organizations.emirate IS 'UAE emirate location (required for healthcare organizations)';
COMMENT ON COLUMN organizations.mfa_method IS 'Primary MFA method: none, sms, email, totp';
COMMENT ON COLUMN organizations.allowed_mfa_methods IS 'Array of allowed MFA methods for the organization';
COMMENT ON COLUMN organizations.ip_whitelist_enabled IS 'Enable IP-based access restrictions';