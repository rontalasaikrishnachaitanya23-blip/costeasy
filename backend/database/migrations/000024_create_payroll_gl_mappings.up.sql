-- ===============================================
-- 000024_create_payroll_gl_mappings.up.sql
-- Links payroll components to GL accounts for auto-posting
-- ===============================================

CREATE TABLE IF NOT EXISTS payroll_gl_mappings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    component_type VARCHAR(50) NOT NULL, -- BASIC, ALLOWANCE, DEDUCTION, etc.
    component_name VARCHAR(100) NOT NULL,
    debit_account_id UUID REFERENCES gl_accounts(id) ON DELETE SET NULL,  -- Expense
    credit_account_id UUID REFERENCES gl_accounts(id) ON DELETE SET NULL, -- Liability/Payable
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (organization_id, component_name)
);

COMMENT ON TABLE payroll_gl_mappings IS 'Maps payroll components (like Basic, HRA, PF) to GL accounts for auto journal posting.';
