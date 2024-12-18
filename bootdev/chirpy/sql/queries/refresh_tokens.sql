-- name: CreateRefeshToken :one
INSERT INTO refresh_tokens (token, expires_at, created_at, updated_at, user_id, revoked_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING token;


-- name: GetRefreshToken :one
SELECT token, expires_at, created_at, updated_at, user_id, revoked_at
FROM refresh_tokens
WHERE token = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = $2,
    updated_at = $3
WHERE token = $1;
