ALTER TABLE attendances
    DROP COLUMN IF EXISTS break_start,
    DROP COLUMN IF EXISTS break_end,
    DROP COLUMN IF EXISTS status,
    DROP COLUMN IF EXISTS notes;

DROP INDEX IF EXISTS idx_attendances_status;
