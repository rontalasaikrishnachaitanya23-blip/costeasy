ALTER TABLE employees
DROP COLUMN IF EXISTS designation_id,
DROP COLUMN IF EXISTS department_id;

DROP TABLE IF EXISTS designations;
DROP TABLE IF EXISTS departments;
