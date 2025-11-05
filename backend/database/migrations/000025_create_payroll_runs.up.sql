-- ===============================================
-- 000025_create_payroll_runs.up.sql
-- Tracks payroll runs and their journal entries
-- ===============================================

CREATE TABLE IF NOT EXISTS payroll_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    payroll_period_id UUID REFERENCES payroll_periods(id) ON DELETE CASCADE,
    processed_by UUID REFERENCES users(id),
    gl_journal_id UUID REFERENCES journal_entries(id),
    status VARCHAR(20) DEFAULT 'DRAFT', -- DRAFT, APPROVED, POSTED
    remarks TEXT,
    processed_at TIMESTAMP,
    approved_at TIMESTAMP,
    posted_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE payroll_runs IS 'Tracks payroll calculation and posting batches linked to GL.';

CREATE INDEX IF NOT EXISTS idx_payroll_runs_period ON payroll_runs(payroll_period_id);
