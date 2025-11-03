-- backend/database/migrations/000006_add_organization_to_accounts.down.sql

-- Restore original unique constraint
ALTER TABLE accounts DROP CONSTRAINT IF EXISTS accounts_org_code_unique;
ALTER TABLE accounts ADD CONSTRAINT accounts_code_unique UNIQUE(code);

-- Drop index
DROP INDEX IF EXISTS idx_accounts_organization;

-- Remove organization_id column
ALTER TABLE accounts DROP COLUMN IF EXISTS organization_id;
