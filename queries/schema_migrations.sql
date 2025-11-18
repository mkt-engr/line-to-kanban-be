-- name: GetCurrentMigrationVersion :one
SELECT COALESCE(MAX(version), -1)::int AS version
FROM schema_migrations
WHERE NOT dirty;

-- name: InsertMigrationVersion :exec
INSERT INTO schema_migrations (version, dirty)
VALUES ($1::int, FALSE);
