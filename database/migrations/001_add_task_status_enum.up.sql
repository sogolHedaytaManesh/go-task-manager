DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'task_status') THEN
CREATE TYPE task_status AS ENUM (
            'pending',
            'in_progress',
            'done',
            'canceled'
        );
END IF;
END$$;
