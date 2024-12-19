-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, password_hash)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING id, created_at, updated_at, email, is_chirpy_red;

-- name: ResetUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, password_hash, is_chirpy_red
FROM users
WHERE email = $1;

-- name: GetUserById :one
SELECT id, created_at, updated_at, email, password_hash, is_chirpy_red
FROM users
WHERE id = $1;

-- name: UpdateUser :one
UPDATE users
SET email = $2,
    password_hash = $3,
    updated_at = $4
WHERE id = $1
RETURNING id, created_at, updated_at, email, is_chirpy_red;

-- name: UpgradeUser :exec
UPDATE users
SET is_chirpy_red = true
WHERE id = $1;
