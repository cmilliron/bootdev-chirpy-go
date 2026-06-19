-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1, 
    $2
)
RETURNING *;

-- name: GetAllChirps :many
SELECT * FROM chirps
ORDER BY 
    CASE WHEN $1::text = 'asc' THEN created_at END ASC,
    CASE WHEN $1::text = 'desc' THEN created_at END DESC,
    CASE WHEN $1::text NOT IN ('asc', 'desc') THEN created_at END ASC;
    
-- name: GetChirpsByAuthor :many
SELECT * FROM chirps
WHERE user_id = $1
ORDER BY
    CASE WHEN $2::text = 'asc' THEN created_at END ASC,
    CASE WHEN $2::text = 'desc' THEN created_at END DESC,
    CASE WHEN $2::text NOT IN ('asc', 'desc') THEN created_at END ASC;

-- name: GetChirpByID :one
SELECT *
FROM chirps
WHERE id = $1;

-- name: DeleteChirpByID :exec
DELETE FROM chirps
WHERE id = $1;