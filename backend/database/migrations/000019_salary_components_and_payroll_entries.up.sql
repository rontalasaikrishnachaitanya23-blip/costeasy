-- ============================================
-- 030 - Salary Components, Employee Salary Details,
--       Payroll Entries, Payroll Entry Lines, and Journal Links
-- ============================================

-- Ensure pgcrypto available for gen_random_uuid
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- ============================================
-- Salary Components (master)
-- ============================================
CREATE TABLE IF NOT EXISTS salary_components (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(150) NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('earning','deduction')),
    calculation_type VARCHAR(20) NOT NULL CHECK (calculation_type IN ('fixed','percentage','formula')) DEFAULT 'fixed',
    value NUMERIC(18,2),        -- used for fixed amounts or percentage base
    percentage NUMERIC(7,4),    -- for percentage based components (e.g., 0.12 = 12%)
    taxable BOOLEAN DEFAULT true,
    applies_to_basic BOOLEAN DEFAULT false,
    gl_account_id UUID,         -- optional link to GL account
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (organization_id, code)
);

CREATE INDEX IF NOT EXISTS idx_salary_components_org ON salary_components(organization_id);
COMMENT ON TABLE salary_components IS 'Master list of earning and deduction components (basic, hra, pf, tax, etc.)';

-- Seed some common components (idempotent)
INSERT INTO salary_components (organization_id, code, name, type, calculation_type, value, percentage, taxable, applies_to_basic, is_active)
VALUES
  -- Note: replace NULL organization_id with actual org UUID when seeding per-organization.
  (NULL, 'BASIC', 'Basic Salary', 'earning', 'fixed', NULL, NULL, true, false, true),
  (NULL, 'HRA', 'House Rent Allowance', 'earning', 'percentage', NULL, 0.40, true, true, true),
  (NULL, 'TRANSPORT', 'Transport Allowance', 'earning', 'fixed', 0.00, NULL, true, false, true),
  (NULL, 'PF', 'Provident Fund / Social Security', 'deduction', 'percentage', NULL, 0.12, false, false, true),
  (NULL, 'TAX', 'Income Tax', 'deduction', 'percentage', NULL, NULL, true, false, true)
ON CONFLICT (organization_id, code) DO NOTHING;

-- ============================================
-- Employee Salary Details (employee level components)
-- ============================================
CREATE TABLE IF NOT EXISTS employee_salary_details (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    component_id UUID NOT NULL REFERENCES salary_components(id) ON DELETE CASCADE,
    amount NUMERIC(18,2) DEFAULT 0.00,   -- fixed override value (if calculation_type = fixed)
    percentage NUMERIC(7,4),             -- override percentage (e.g., 0.50 for 50% of basic)
    is_active BOOLEAN DEFAULT true,
    start_date DATE DEFAULT CURRENT_DATE,
    end_date DATE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(employee_id, component_id)
);

CREATE INDEX IF NOT EXISTS idx_employee_salary_employee ON employee_salary_details(employee_id);
COMMENT ON TABLE employee_salary_details IS 'Mapping of employees to salary components (amount or percent, with effective dates)';

-- ============================================
-- Payroll Entries (header) - one per employee per payroll_period
-- ============================================
CREATE TABLE IF NOT EXISTS payroll_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payroll_period_id UUID NOT NULL REFERENCES payroll_periods(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL,
    gross_earnings NUMERIC(18,2) DEFAULT 0.00,
    total_deductions NUMERIC(18,2) DEFAULT 0.00,
    net_pay NUMERIC(18,2) DEFAULT 0.00,
    status VARCHAR(20) NOT NULL DEFAULT 'draft', -- draft, processed, posted, cancelled
    processed_at TIMESTAMPTZ,
    posted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (payroll_period_id, employee_id)
);

CREATE INDEX IF NOT EXISTS idx_payroll_entries_period ON payroll_entries(payroll_period_id);
CREATE INDEX IF NOT EXISTS idx_payroll_entries_emp ON payroll_entries(employee_id);
COMMENT ON TABLE payroll_entries IS 'Payroll header per employee per payroll period';

-- ============================================
-- Payroll Entry Lines (detail per salary component)
-- ============================================
CREATE TABLE IF NOT EXISTS payroll_entry_lines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payroll_entry_id UUID NOT NULL REFERENCES payroll_entries(id) ON DELETE CASCADE,
    component_id UUID NOT NULL REFERENCES salary_components(id) ON DELETE RESTRICT,
    component_code VARCHAR(50),
    description TEXT,
    amount NUMERIC(18,2) NOT NULL DEFAULT 0.00,
    is_earning BOOLEAN NOT NULL DEFAULT true,
    gl_account_id UUID,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payroll_entry_lines_entry ON payroll_entry_lines(payroll_entry_id);
COMMENT ON TABLE payroll_entry_lines IS 'Detailed breakdown of payroll per component (earnings/deductions)';

-- ============================================
-- Payroll -> GL Journal Link
-- ============================================
CREATE TABLE IF NOT EXISTS payroll_journal_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payroll_entry_id UUID NOT NULL REFERENCES payroll_entries(id) ON DELETE CASCADE,
    journal_entry_id UUID, -- reference to GL journal (store UUID of GL journal entry if available)
    posted_by UUID,
    posted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payroll_journal_entry ON payroll_journal_links(journal_entry_id);
COMMENT ON TABLE payroll_journal_links IS 'Links payroll entries to GL journal entries once posted';

-- ============================================
-- Optional: trigger to update updated_at timestamps
-- (simple generic trigger function for tables that use updated_at)
-- ============================================
CREATE OR REPLACE FUNCTION set_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Attach trigger to tables
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_trigger WHERE tgname = 'salary_components_updated_at_trg'
  ) THEN
    CREATE TRIGGER salary_components_updated_at_trg
    BEFORE UPDATE ON salary_components
    FOR EACH ROW EXECUTE FUNCTION set_updated_at_column();
  END IF;
END;
$$;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_trigger WHERE tgname = 'employee_salary_details_updated_at_trg'
  ) THEN
    CREATE TRIGGER employee_salary_details_updated_at_trg
    BEFORE UPDATE ON employee_salary_details
    FOR EACH ROW EXECUTE FUNCTION set_updated_at_column();
  END IF;
END;
$$;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_trigger WHERE tgname = 'payroll_entries_updated_at_trg'
  ) THEN
    CREATE TRIGGER payroll_entries_updated_at_trg
    BEFORE UPDATE ON payroll_entries
    FOR EACH ROW EXECUTE FUNCTION set_updated_at_column();
  END IF;
END;
$$;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_trigger WHERE tgname = 'payroll_entry_lines_updated_at_trg'
  ) THEN
    CREATE TRIGGER payroll_entry_lines_updated_at_trg
    BEFORE UPDATE ON payroll_entry_lines
    FOR EACH ROW EXECUTE FUNCTION set_updated_at_column();
  END IF;
END;
$$;
