-- ===============================
-- 000021_create_attendance_statuses.up.sql
-- Master table for attendance statuses
-- ===============================

CREATE TABLE IF NOT EXISTS attendance_statuses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(20) UNIQUE NOT NULL,        -- PRESENT, ABSENT, LEAVE, HOLIDAY, WEEKOFF, etc.
    display_name VARCHAR(50) NOT NULL,       -- Human-friendly name
    description TEXT,
    is_paid BOOLEAN DEFAULT true,            -- Whether this status counts as paid day
    affects_payroll BOOLEAN DEFAULT true,    -- Should be considered during payroll calculation
    color_code VARCHAR(10),                  -- For UI (e.g. green, red, blue)
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE attendance_statuses IS 'Master list of attendance statuses like PRESENT, ABSENT, LEAVE, HOLIDAY, etc.';


-- Seed default statuses
INSERT INTO attendance_statuses (code, display_name, description, is_paid, affects_payroll, color_code)
VALUES
('PRESENT', 'Present', 'Employee attended work', true, true, '#4CAF50'),
('ABSENT', 'Absent', 'Employee did not attend work', false, true, '#F44336'),
('LEAVE', 'Leave', 'Employee on approved leave', true, true, '#2196F3'),
('HOLIDAY', 'Holiday', 'Public or organizational holiday', true, false, '#FFEB3B'),
('WEEKOFF', 'Weekly Off', 'Scheduled weekly off day', true, false, '#9E9E9E'),
('LATE', 'Late Arrival', 'Employee arrived late', true, true, '#FF9800'),
('EARLYOUT', 'Early Out', 'Employee left early', true, true, '#FF5722'),
('OVERTIME', 'Overtime', 'Worked beyond shift hours', true, true, '#8BC34A'),
('WORKFROMHOME', 'Work From Home', 'Approved remote work day', true, true, '#03A9F4'),
('ONFIELD', 'On Field', 'Working outside office (client/site)', true, true, '#00BCD4')
ON CONFLICT (code) DO NOTHING;
