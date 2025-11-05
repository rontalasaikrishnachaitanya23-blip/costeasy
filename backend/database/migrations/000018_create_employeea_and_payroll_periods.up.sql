-- ============================================
-- Payroll Core Tables
-- Employees + Payroll Periods
-- ============================================

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- ============================================
-- Employees
-- ============================================
CREATE TABLE IF NOT EXISTS employees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL,
    country_id UUID REFERENCES countries(id) ON DELETE SET NULL,
    employee_code VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100),
    email VARCHAR(150) UNIQUE,
    phone VARCHAR(30),
    date_of_birth DATE,
    gender VARCHAR(10),
    nationality VARCHAR(50),
    date_of_joining DATE NOT NULL,
    date_of_exit DATE,
    department VARCHAR(100),
    designation VARCHAR(100),
    work_location VARCHAR(100),
    contract_type VARCHAR(50) DEFAULT 'permanent', -- permanent, temporary, intern
    salary_currency VARCHAR(3) DEFAULT 'AED',
    base_salary DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (organization_id, employee_code)
);

CREATE INDEX IF NOT EXISTS idx_employees_org ON employees(organization_id);
CREATE INDEX IF NOT EXISTS idx_employees_country ON employees(country_id);
CREATE INDEX IF NOT EXISTS idx_employees_active ON employees(is_active);

-- ============================================
-- Payroll Periods
-- ============================================
CREATE TABLE IF NOT EXISTS payroll_periods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    month INT NOT NULL,
    year INT NOT NULL,
    is_locked BOOLEAN DEFAULT false,
    processed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (organization_id, month, year)
);

CREATE INDEX IF NOT EXISTS idx_payroll_periods_org ON payroll_periods(organization_id);
CREATE INDEX IF NOT EXISTS idx_payroll_periods_range ON payroll_periods(start_date, end_date);

-- ============================================
-- Payroll Period Locks (for audit)
-- ============================================
CREATE TABLE IF NOT EXISTS payroll_period_locks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payroll_period_id UUID REFERENCES payroll_periods(id) ON DELETE CASCADE,
    locked_by UUID,
    locked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    reason TEXT
);

CREATE INDEX IF NOT EXISTS idx_payroll_period_lock_period ON payroll_period_locks(payroll_period_id);
