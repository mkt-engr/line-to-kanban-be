-- Migration 003: Rename messages to tasks
--
-- This migration renames the database objects from 'message' to 'task'
-- to align with the domain model.

-- Rename ENUM type from message_status to task_status
ALTER TYPE message_status RENAME TO task_status;

-- Rename table from messages to tasks
ALTER TABLE messages RENAME TO tasks;

-- Rename indexes to match new table name
ALTER INDEX messages_created_at_idx RENAME TO tasks_created_at_idx;
ALTER INDEX messages_status_idx RENAME TO tasks_status_idx;
ALTER INDEX messages_user_id_idx RENAME TO tasks_user_id_idx;
