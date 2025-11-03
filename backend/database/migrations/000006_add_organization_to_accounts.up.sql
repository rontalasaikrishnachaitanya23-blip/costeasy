-- Change gl_accounts.id from INTEGER to UUID
-- WARNING: This will regenerate all account IDs

-- Step 1: Add temporary UUID column
ALTER TABLE gl_accounts ADD COLUMN temp_id UUID DEFAULT gen_random_uuid();

-- Step 2: Fill all rows with UUID values
UPDATE gl_accounts SET temp_id = gen_random_uuid() WHERE temp_id IS NULL;

-- Step 3: Drop old INTEGER id column and its primary key constraint
ALTER TABLE gl_accounts DROP CONSTRAINT IF EXISTS gl_accounts_pkey CASCADE;

-- Step 4: Drop the old id column
ALTER TABLE gl_accounts DROP COLUMN id;

-- Step 5: Rename temp_id to id
ALTER TABLE gl_accounts RENAME COLUMN temp_id TO id;

-- Step 6: Set id as NOT NULL
ALTER TABLE gl_accounts ALTER COLUMN id SET NOT NULL;

-- Step 7: Make id the primary key
ALTER TABLE gl_accounts ADD PRIMARY KEY (id);

-- Step 8: Recreate only the indexes that should exist
-- Check if code column exists before creating index
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'gl_accounts' AND column_name = 'code'
    ) THEN
        CREATE INDEX IF NOT EXISTS idx_gl_accounts_code ON gl_accounts(code);
    END IF;
END $$;
