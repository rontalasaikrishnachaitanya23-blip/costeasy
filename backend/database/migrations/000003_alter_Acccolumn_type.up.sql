CREATE TYPE account_type AS ENUM ('Asset', 'Liability', 'Equity', 'Revenue', 'Expense');

ALTER TABLE accounts ALTER COLUMN type TYPE account_type USING type::account_type;
