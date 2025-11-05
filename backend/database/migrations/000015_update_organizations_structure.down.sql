-- Migration 000015 Rollback: Revert Organizations Table Structure

-- Drop constraints
ALTER TABLE organizations DROP CONSTRAINT IF EXISTS chk_organizations_type;
ALTER TABLE organizations DROP CONSTRAINT IF EXISTS chk_organizations_emirate;
ALTER TABLE organizations DROP CONSTRAINT IF EXISTS chk_organizations_mfa_method;

-- Drop indexes
DROP INDEX IF EXISTS idx_organizations_type;
DROP INDEX IF EXISTS idx_organizations_emirate;
DROP INDEX IF EXISTS idx_organizations_code;
DROP INDEX IF EXISTS idx_organizations_currency;

-- Drop columns (in reverse order)
ALTER TABLE organizations 
    DROP COLUMN IF EXISTS allowed_mfa_methods,
    DROP COLUMN IF EXISTS mfa_method,
    DROP COLUMN IF EXISTS mfa_enforced,
    DROP COLUMN IF EXISTS mfa_enabled,
    DROP COLUMN IF EXISTS description,
    DROP COLUMN IF EXISTS establishment_id,
    DROP COLUMN IF EXISTS license_expiry,
    DROP COLUMN IF EXISTS license_number,
    DROP COLUMN IF EXISTS currency,
    DROP COLUMN IF EXISTS website,
    DROP COLUMN IF EXISTS area,
    DROP COLUMN IF EXISTS emirate,
    DROP COLUMN IF EXISTS type,
    DROP COLUMN IF EXISTS code;