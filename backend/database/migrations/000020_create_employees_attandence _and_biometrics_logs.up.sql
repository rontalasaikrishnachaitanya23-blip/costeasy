-- ===============================
-- 000020_create_attendance_tables.up.sql
-- Attendance + Biometric + Upload Tracking
-- ===============================

-- 1️⃣ Upload batch tracking (must come first)
CREATE TABLE IF NOT EXISTS attendance_upload_batches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id),
    uploaded_by UUID REFERENCES users(id),
    file_name VARCHAR(255),
    total_rows INT DEFAULT 0,
    processed_rows INT DEFAULT 0,
    failed_rows INT DEFAULT 0,
    upload_status VARCHAR(20) DEFAULT 'PENDING', -- PENDING, COMPLETED, FAILED
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE attendance_upload_batches IS 'Tracks attendance file uploads (Excel/CSV) for audit and rollback.';


-- 2️⃣ Attendance records
CREATE TABLE IF NOT EXISTS attendance_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    attendance_date DATE NOT NULL,
    check_in TIMESTAMP,
    check_out TIMESTAMP,
    total_hours DECIMAL(5,2),
    status VARCHAR(20) DEFAULT 'PRESENT',  -- PRESENT, ABSENT, LEAVE, HOLIDAY, WEEKOFF
    source VARCHAR(20) DEFAULT 'MANUAL',   -- MANUAL, BIOMETRIC, UPLOAD, API
    device_id UUID,
    remarks TEXT,
    approved_by UUID REFERENCES users(id),
    upload_batch_id UUID REFERENCES attendance_upload_batches(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(employee_id, attendance_date)
);

COMMENT ON TABLE attendance_records IS 'Stores attendance for each employee (manual or biometric).';


-- 3️⃣ Biometric logs (optional raw data)
CREATE TABLE IF NOT EXISTS biometric_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID,
    employee_code VARCHAR(50),
    log_time TIMESTAMP NOT NULL,
    log_type VARCHAR(10), -- IN / OUT
    sync_status VARCHAR(20) DEFAULT 'PENDING', -- PENDING / PROCESSED
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE biometric_logs IS 'Raw biometric logs synced from attendance devices.';


-- 4️⃣ Attendance exceptions (late in, missing check-out)
CREATE TABLE IF NOT EXISTS attendance_exceptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attendance_id UUID REFERENCES attendance_records(id) ON DELETE CASCADE,
    exception_type VARCHAR(50), -- LATE_IN, EARLY_OUT, MISSING_OUT, NO_CHECKIN
    remarks TEXT,
    resolved BOOLEAN DEFAULT false,
    resolved_by UUID REFERENCES users(id),
    resolved_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE attendance_exceptions IS 'Tracks anomalies or manual exceptions in attendance.';
