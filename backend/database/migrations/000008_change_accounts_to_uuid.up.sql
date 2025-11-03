-- Change gl_accounts.id from INTEGER to UUID
-- WARNING: This will regenerate all account IDs

-- Step 1: Add temporary UUID column
ALTER TABLE gl_accounts ADD COLUMN temp_id UUID DEFAULT gen_random_uuid();

-- Step 2: Update all rows to have UUID values
UPDATE gl_accounts SET temp_id = gen_random_uuid() WHERE temp_id IS NULL;

-- Step 3: Drop old primary key constraint
ALTER TABLE gl_accounts DROP CONSTRAINT IF EXISTS gl_accounts_pkey CASCADE;

-- Step 4: Drop the old id column and its sequence
DROP SEQUENCE IF EXISTS gl_accounts_id_seq CASCADE;
ALTER TABLE gl_accounts DROP COLUMN id;

-- Step 5: Rename temp_id to id
ALTER TABLE gl_accounts RENAME COLUMN temp_id TO id;

-- Step 6: Set id as NOT NULL
ALTER TABLE gl_accounts ALTER COLUMN id SET NOT NULL;

-- Step 7: Make id the primary key
ALTER TABLE gl_accounts ADD PRIMARY KEY (id);

-- Step 8: Recreate constraints and indexes ONLY IF they don't exist
-- The code constraint already exists, so skip it
-- Just recreate the type index
CREATE INDEX IF NOT EXISTS idx_gl_accounts_type ON gl_accounts(type);
