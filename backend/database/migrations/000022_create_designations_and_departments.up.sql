-- ===============================================
-- 000022_create_designations_and_departments.up.sql
-- Departments and Designations for Payroll & HR
-- ===============================================

-- 1️⃣ Departments
CREATE TABLE IF NOT EXISTS departments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    code VARCHAR(50),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    head_id UUID REFERENCES employees(id), -- optional: department head
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (organization_id, name)
);

COMMENT ON TABLE departments IS 'Defines organizational departments (e.g., HR, Finance, Operations).';
COMMENT ON COLUMN departments.head_id IS 'Employee who leads the department.';


-- 2️⃣ Designations
CREATE TABLE IF NOT EXISTS designations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    department_id UUID REFERENCES departments(id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50),
    level INT DEFAULT 1, -- Hierarchy level
    grade VARCHAR(20),   -- Optional grade (A1, B2, etc.)
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (organization_id, name)
);

COMMENT ON TABLE designations IS 'Defines roles or titles within departments (e.g., Accountant, Nurse, Developer).';
COMMENT ON COLUMN designations.level IS 'Hierarchy level; lower = higher seniority.';
COMMENT ON COLUMN designations.grade IS 'Optional salary or band grade for HR structure.';


-- 3️⃣ Extend employees to include department & designation
ALTER TABLE employees
ADD COLUMN IF NOT EXISTS department_id UUID REFERENCES departments(id) ON DELETE SET NULL,
ADD COLUMN IF NOT EXISTS designation_id UUID REFERENCES designations(id) ON DELETE SET NULL;
