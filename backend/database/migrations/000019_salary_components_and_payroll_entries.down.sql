-- Rollback: drop payroll components and entries
DROP TRIGGER IF EXISTS payroll_entry_lines_updated_at_trg ON payroll_entry_lines;
DROP TRIGGER IF EXISTS payroll_entries_updated_at_trg ON payroll_entries;
DROP TRIGGER IF EXISTS employee_salary_details_updated_at_trg ON employee_salary_details;
DROP TRIGGER IF EXISTS salary_components_updated_at_trg ON salary_components;

DROP FUNCTION IF EXISTS set_updated_at_column();

DROP TABLE IF EXISTS payroll_journal_links;
DROP TABLE IF EXISTS payroll_entry_lines;
DROP TABLE IF EXISTS payroll_entries;
DROP TABLE IF EXISTS employee_salary_details;
DROP TABLE IF EXISTS salary_components;
