-- Tracks each uploaded Excel file
CREATE TABLE IF NOT EXISTS employee_import_batches (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id  UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    uploaded_by      UUID NOT NULL REFERENCES users(id),
    file_name        TEXT NOT NULL,
    total_rows       INT NOT NULL DEFAULT 0,
    valid_rows       INT NOT NULL DEFAULT 0,
    invalid_rows     INT NOT NULL DEFAULT 0,
    status           VARCHAR(20) NOT NULL DEFAULT 'processing', -- processing/success/partial/failed
    created_at       TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at     TIMESTAMP
);

-- Row-level errors for visibility/traceability
CREATE TABLE IF NOT EXISTS employee_import_errors (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    batch_id    UUID NOT NULL REFERENCES employee_import_batches(id) ON DELETE CASCADE,
    row_number  INT NOT NULL,
    error_code  VARCHAR(50) NOT NULL,
    message     TEXT NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_employee_import_errors_batch ON employee_import_errors(batch_id);
