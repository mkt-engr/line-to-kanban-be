-- -- name: GetMessage :one
-- SELECT * FROM messages
-- WHERE id = $1 LIMIT 1;

-- name: ListMessagesByUser :many
SELECT * FROM messages
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: CreateMessage :one
INSERT INTO messages (
  content, status, user_id
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- -- name: UpdateMessageStatus :exec
-- UPDATE messages
-- SET status = $2, updated_at = now()
-- WHERE id = $1;

-- -- name: DeleteMessage :exec
-- DELETE FROM messages
-- WHERE id = $1;
