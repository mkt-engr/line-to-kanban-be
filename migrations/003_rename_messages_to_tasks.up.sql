-- Rename ENUM type from message_status to task_status (if exists)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'message_status') THEN
        ALTER TYPE message_status RENAME TO task_status;
    END IF;
END $$;

-- Rename table from messages to tasks (if exists)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'messages') THEN
        ALTER TABLE messages RENAME TO tasks;
    END IF;
END $$;

-- Rename indexes to match new table name (if they exist)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'messages_created_at_idx') THEN
        ALTER INDEX messages_created_at_idx RENAME TO tasks_created_at_idx;
    END IF;

    IF EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'messages_status_idx') THEN
        ALTER INDEX messages_status_idx RENAME TO tasks_status_idx;
    END IF;

    IF EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'messages_user_id_idx') THEN
        ALTER INDEX messages_user_id_idx RENAME TO tasks_user_id_idx;
    END IF;
END $$;
