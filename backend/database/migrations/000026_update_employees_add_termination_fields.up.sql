ALTER TABLE employees
ADD COLUMN employment_status VARCHAR(50) NOT NULL DEFAULT 'active',
ADD COLUMN joined_at DATE NOT NULL,
ADD COLUMN relieved_at DATE,
ADD COLUMN termination_reason TEXT,
ADD COLUMN is_salary_stopped BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN final_settlement_generated BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN final_settlement_date TIMESTAMP;
