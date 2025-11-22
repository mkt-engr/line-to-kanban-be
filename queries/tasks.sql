-- -- name: GetTask :one
-- SELECT * FROM tasks
-- WHERE id = $1 LIMIT 1;

-- name: ListTasksByUser :many
SELECT * FROM tasks
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = $1 AND user_id = $2;

-- name: CreateTask :one
INSERT INTO tasks (
  content, status, user_id
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- -- name: UpdateTaskStatus :exec
-- UPDATE tasks
-- SET status = $2, updated_at = now()
-- WHERE id = $1;

-- -- name: DeleteTask :exec
-- DELETE FROM tasks
-- WHERE id = $1;
