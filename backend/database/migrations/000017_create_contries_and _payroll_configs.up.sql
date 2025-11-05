-- ===============================
-- Payroll Base Migration
-- Countries & Country Payroll Configs
-- ===============================

-- Enable pgcrypto (for gen_random_uuid)
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- ============================================
-- Countries Master
-- ============================================
CREATE TABLE IF NOT EXISTS countries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(3) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    currency_code VARCHAR(3) NOT NULL,
    phone_code VARCHAR(10),
    date_format VARCHAR(20) DEFAULT 'DD-MM-YYYY',
    time_zone VARCHAR(50),
    working_days_per_week INT DEFAULT 5,
    weekend_days TEXT[],
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- Seed Data
-- ============================================
INSERT INTO countries (code, name, currency_code, phone_code, time_zone, weekend_days, working_days_per_week) VALUES
-- India
('IN', 'India', 'INR', '+91', 'Asia/Kolkata', ARRAY['Sunday'], 6),

-- GCC/Arab Countries
('AE', 'United Arab Emirates', 'AED', '+971', 'Asia/Dubai', ARRAY['Friday', 'Saturday'], 5),
('SA', 'Saudi Arabia', 'SAR', '+966', 'Asia/Riyadh', ARRAY['Friday', 'Saturday'], 5),
('QA', 'Qatar', 'QAR', '+974', 'Asia/Qatar', ARRAY['Friday', 'Saturday'], 5),
('KW', 'Kuwait', 'KWD', '+965', 'Asia/Kuwait', ARRAY['Friday', 'Saturday'], 5),
('OM', 'Oman', 'OMR', '+968', 'Asia/Muscat', ARRAY['Friday', 'Saturday'], 5),
('BH', 'Bahrain', 'BHD', '+973', 'Asia/Bahrain', ARRAY['Friday', 'Saturday'], 5),

-- Other Arab Countries
('EG', 'Egypt', 'EGP', '+20', 'Africa/Cairo', ARRAY['Friday', 'Saturday'], 5),
('JO', 'Jordan', 'JOD', '+962', 'Asia/Amman', ARRAY['Friday', 'Saturday'], 5),
('LB', 'Lebanon', 'LBP', '+961', 'Asia/Beirut', ARRAY['Saturday', 'Sunday'], 5),

-- Australia & Others
('AU', 'Australia', 'AUD', '+61', 'Australia/Sydney', ARRAY['Saturday', 'Sunday'], 5),
('NZ', 'New Zealand', 'NZD', '+64', 'Pacific/Auckland', ARRAY['Saturday', 'Sunday'], 5),
('US', 'United States', 'USD', '+1', 'America/New_York', ARRAY['Saturday', 'Sunday'], 5),
('GB', 'United Kingdom', 'GBP', '+44', 'Europe/London', ARRAY['Saturday', 'Sunday'], 5),
('SG', 'Singapore', 'SGD', '+65', 'Asia/Singapore', ARRAY['Saturday', 'Sunday'], 5)
ON CONFLICT (code) DO NOTHING;

-- ============================================
-- Country Payroll Configurations
-- ============================================
CREATE TABLE IF NOT EXISTS country_payroll_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country_id UUID REFERENCES countries(id) ON DELETE CASCADE,
    has_income_tax BOOLEAN DEFAULT true,
    has_social_security BOOLEAN DEFAULT true,
    has_professional_tax BOOLEAN DEFAULT false,
    has_gratuity BOOLEAN DEFAULT false,
    minimum_wage DECIMAL(15,2),
    overtime_multiplier DECIMAL(5,2) DEFAULT 1.5,
    probation_period_days INT DEFAULT 90,
    notice_period_days INT DEFAULT 30,
    annual_leave_days INT DEFAULT 21,
    sick_leave_days INT DEFAULT 12,
    maternity_leave_days INT DEFAULT 90,
    paternity_leave_days INT DEFAULT 5,
    config_json JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(country_id)
);

CREATE INDEX IF NOT EXISTS idx_country_payroll_country ON country_payroll_configs(country_id);
