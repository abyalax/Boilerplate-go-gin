
-- Create User and return data
-- name: CreateUser :one
INSERT INTO users (
    name,
    password,
    email
)
VALUES (
    $1,
    $2,
    $3
)
RETURNING id, name, password, email;

-- Get base data User by ID
-- name: GetUserByID :one
SELECT id, name, password, email
FROM users
WHERE id = $1;

-- Get base data User by Email
-- name: GetUserByEmail :one
SELECT id, name, password, email
FROM users
WHERE email = $1;

-- Get Complete data users with role and permission (exclude password)
-- name: ListUsers :many
SELECT
    u.id AS id,
    u.name AS name,
    u.email AS email,
    COALESCE(
        (
            SELECT jsonb_agg(role_obj)
            FROM (
                SELECT DISTINCT
                    r.id AS role_id,
                    r.name AS role_name,
                    COALESCE(
                        (
                            SELECT jsonb_agg(
                                jsonb_build_object(
                                    'permission_id', p.id,
                                    'permission_key', p.key,
                                    'permission_name', p.name
                                )
                            )
                            FROM role_permissions rp
                            JOIN permissions p ON rp.permission_id = p.id
                            WHERE rp.role_id = r.id
                        ), '[]'::jsonb
                    ) AS permissions
                FROM roles r
                JOIN user_roles ur ON ur.role_id = r.id
                WHERE ur.user_id = u.id
            ) role_obj
        ), '[]'::jsonb
    ) AS roles
FROM users u;

-- Update base data users and return
-- name: UpdateUser :one
UPDATE users
SET name = COALESCE($2, name),
    email = COALESCE($3, email),
    password = COALESCE($4, password)
WHERE id = $1
RETURNING id, name, password, email;

-- Delete data user by ID
-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- List all users with passwords (for repository internal use)
-- name: ListAllUsers :many
SELECT id, name, password, email
FROM users
ORDER BY id DESC;

