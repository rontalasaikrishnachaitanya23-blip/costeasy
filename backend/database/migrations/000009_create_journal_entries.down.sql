-- Drop triggers
DROP TRIGGER IF EXISTS journal_entries_updated_at ON journal_entries;
DROP FUNCTION IF EXISTS update_journal_entry_updated_at();

-- Drop indexes
DROP INDEX IF EXISTS idx_journal_lines_reference;          -- ✅ ADDED
DROP INDEX IF EXISTS idx_journal_lines_account_id;
DROP INDEX IF EXISTS idx_journal_lines_journal_entry_id;
DROP INDEX IF EXISTS idx_journal_entries_reference;        -- ✅ ADDED
DROP INDEX IF EXISTS idx_journal_entries_entry_number;
DROP INDEX IF EXISTS idx_journal_entries_transaction_date;
DROP INDEX IF EXISTS idx_journal_entries_status;
DROP INDEX IF EXISTS idx_journal_entries_org_id;

-- Drop tables
DROP TABLE IF EXISTS journal_lines;
DROP TABLE IF EXISTS journal_entries;
