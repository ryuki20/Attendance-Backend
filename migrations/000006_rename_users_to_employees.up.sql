ALTER TABLE users RENAME TO employees;

ALTER TABLE attendances RENAME COLUMN user_id TO employee_id;

ALTER INDEX IF EXISTS users_pkey RENAME TO employees_pkey;
ALTER INDEX IF EXISTS users_email_key RENAME TO employees_email_key;
