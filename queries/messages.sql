-- -- name: GetMessage :one
-- SELECT * FROM messages
-- WHERE id = $1 LIMIT 1;

-- -- name: ListMessages :many
-- SELECT * FROM messages
-- ORDER BY created_at DESC;

-- name: CreateMessage :one
INSERT INTO messages (
  content, status
) VALUES (
  $1, $2
)
RETURNING *;

-- -- name: UpdateMessageStatus :exec
-- UPDATE messages
-- SET status = $2, updated_at = now()
-- WHERE id = $1;

-- -- name: DeleteMessage :exec
-- DELETE FROM messages
-- WHERE id = $1;
