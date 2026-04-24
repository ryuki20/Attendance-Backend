ALTER TABLE employees RENAME TO users;

ALTER TABLE attendances RENAME COLUMN employee_id TO user_id;

ALTER INDEX IF EXISTS employees_pkey RENAME TO users_pkey;
ALTER INDEX IF EXISTS employees_email_key RENAME TO users_email_key;
