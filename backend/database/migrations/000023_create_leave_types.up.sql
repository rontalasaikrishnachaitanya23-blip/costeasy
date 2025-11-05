-- ===============================================
-- 000023_create_leave_types.up.sql
-- Defines standard leave categories (Paid, Unpaid, Sick, etc.)
-- ===============================================

CREATE TABLE IF NOT EXISTS leave_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_paid BOOLEAN DEFAULT true,
    max_days_per_year INT DEFAULT 30,
    carry_forward BOOLEAN DEFAULT true,
    requires_approval BOOLEAN DEFAULT true,
    affects_payroll BOOLEAN DEFAULT true,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (organization_id, code)
);

COMMENT ON TABLE leave_types IS 'Master table for leave categories affecting payroll and attendance.';

-- ✅ Fixed seeding
INSERT INTO leave_types (organization_id, code, name, description, is_paid, affects_payroll)
VALUES
    (NULL, 'AL', 'Annual Leave', 'Paid yearly leave', true, true),
    (NULL, 'SL', 'Sick Leave', 'Paid sick leave', true, true),
    (NULL, 'UL', 'Unpaid Leave', 'Leave without pay', false, true),
    (NULL, 'ML', 'Maternity Leave', 'Maternity leave as per labor law', true, true),
    (NULL, 'PL', 'Paternity Leave', 'Short paid leave for fathers', true, true)
ON CONFLICT DO NOTHING;
