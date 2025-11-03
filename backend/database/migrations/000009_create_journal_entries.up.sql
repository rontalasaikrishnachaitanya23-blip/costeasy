-- Drop tables if they exist (clean slate)
DROP TABLE IF EXISTS journal_lines CASCADE;
DROP TABLE IF EXISTS journal_entries CASCADE;
DROP FUNCTION IF EXISTS update_journal_entry_updated_at() CASCADE;

-- Create journal_entries table
CREATE TABLE journal_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL,           -- UUID to match shafafiya_org_settings.id
    entry_number VARCHAR(50) NOT NULL,
    transaction_date DATE NOT NULL,
    posting_date TIMESTAMP,
    reference VARCHAR(100),
    description TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    total_debit DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    total_credit DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    created_by UUID NOT NULL,
    posted_by UUID,
    reversed_by UUID,
    reversal_of UUID,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create journal_lines table
CREATE TABLE journal_lines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    journal_entry_id UUID NOT NULL,
    account_id UUID NOT NULL,                -- NOW UUID to match gl_accounts.id
    line_number INT NOT NULL,
    reference VARCHAR(100),
    description VARCHAR(255) NOT NULL,
    debit DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    credit DECIMAL(15, 2) NOT NULL DEFAULT 0.00
);

-- Add foreign key constraints
ALTER TABLE journal_entries 
    ADD CONSTRAINT fk_journal_entries_organization 
    FOREIGN KEY (organization_id) REFERENCES shafafiya_org_settings(id) ON DELETE CASCADE;

ALTER TABLE journal_entries 
    ADD CONSTRAINT fk_journal_entries_reversal 
    FOREIGN KEY (reversal_of) REFERENCES journal_entries(id);

ALTER TABLE journal_lines 
    ADD CONSTRAINT fk_journal_lines_entry 
    FOREIGN KEY (journal_entry_id) REFERENCES journal_entries(id) ON DELETE CASCADE;

ALTER TABLE journal_lines 
    ADD CONSTRAINT fk_journal_lines_account 
    FOREIGN KEY (account_id) REFERENCES gl_accounts(id);

-- Add check constraints
ALTER TABLE journal_entries 
    ADD CONSTRAINT check_journal_status 
    CHECK (status IN ('DRAFT', 'POSTED', 'VOID', 'REVERSED'));

ALTER TABLE journal_entries 
    ADD CONSTRAINT check_journal_balance 
    CHECK (total_debit = total_credit OR status = 'DRAFT');

ALTER TABLE journal_lines 
    ADD CONSTRAINT check_debit_credit 
    CHECK ((debit > 0 AND credit = 0) OR (credit > 0 AND debit = 0));

ALTER TABLE journal_lines 
    ADD CONSTRAINT check_amounts_positive 
    CHECK (debit >= 0 AND credit >= 0);

-- Add unique constraints
ALTER TABLE journal_entries 
    ADD CONSTRAINT unique_entry_number 
    UNIQUE (organization_id, entry_number);

ALTER TABLE journal_lines 
    ADD CONSTRAINT unique_line_number 
    UNIQUE (journal_entry_id, line_number);

-- Create indexes
CREATE INDEX idx_journal_entries_org_id ON journal_entries(organization_id);
CREATE INDEX idx_journal_entries_status ON journal_entries(status);
CREATE INDEX idx_journal_entries_transaction_date ON journal_entries(transaction_date);
CREATE INDEX idx_journal_entries_entry_number ON journal_entries(entry_number);
CREATE INDEX idx_journal_entries_reference ON journal_entries(reference);
CREATE INDEX idx_journal_lines_journal_entry_id ON journal_lines(journal_entry_id);
CREATE INDEX idx_journal_lines_account_id ON journal_lines(account_id);
CREATE INDEX idx_journal_lines_reference ON journal_lines(reference);

-- Create trigger function
CREATE OR REPLACE FUNCTION update_journal_entry_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger
CREATE TRIGGER journal_entries_updated_at
    BEFORE UPDATE ON journal_entries
    FOR EACH ROW
    EXECUTE FUNCTION update_journal_entry_updated_at();
