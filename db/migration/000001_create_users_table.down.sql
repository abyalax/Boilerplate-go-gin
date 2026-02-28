-- Drop junction tables first (they have foreign keys)
DROP TABLE IF EXISTS "user_roles" CASCADE;
DROP TABLE IF EXISTS "role_permissions" CASCADE;

-- Drop main tables
DROP TABLE IF EXISTS "users" CASCADE;
DROP TABLE IF EXISTS "permissions" CASCADE;
DROP TABLE IF EXISTS "roles" CASCADE;

-- Drop indexes
DROP INDEX IF EXISTS "users_email_key";
DROP INDEX IF EXISTS "users_name_idx";
DROP INDEX IF EXISTS "permissions_key_unique";
DROP INDEX IF EXISTS "roles_name_unique";
