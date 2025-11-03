DELETE FROM schema_migrations WHERE version = '000008' OR version = '000009';

DROP TABLE IF EXISTS journal_lines CASCADE;
DROP TABLE IF EXISTS journal_entries CASCADE;
