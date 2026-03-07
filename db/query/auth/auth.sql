-- Register user
-- name: CreateUser :one
INSERT INTO users (name, email, password)
VALUES ($1, $2, $3)
RETURNING id, name, email, password;

-- Login user by email
-- name: GetUserByEmail :one
SELECT id, name, email, password
FROM users
WHERE email = $1;

-- Login user by name
-- name: GetUserByName :one
SELECT id, name, email, password
FROM users
WHERE name = $1;

-- Get user with roles and permissions
-- name: GetUserWithPermissions :many
SELECT 
    u.id AS user_id,
    u.name AS user_name,
    u.email AS user_email,
    r.id AS role_id,
    r.name AS role_name,
    p.id AS permission_id,
    p.key AS permission_key,
    p.name AS permission_name
FROM users u
LEFT JOIN user_roles ur ON u.id = ur.user_id
LEFT JOIN roles r ON ur.role_id = r.id
LEFT JOIN role_permissions rp ON r.id = rp.role_id
LEFT JOIN permissions p ON rp.permission_id = p.id
WHERE u.id = $1;