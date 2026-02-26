-- This migration adds a role_id column to the users table, sets a default role for existing users,
-- makes the column required, and enforces a foreign key constraint to roles.
ALTER TABLE users
ADD COLUMN IF NOT EXISTS role_id INT;

-- Set default role for existing users
UPDATE users
set role_id = (
SELECT id FROM roles WHERE name = 'user' LIMIT 1
);

-- Make role_id NOT NULL and add foreign key constraint
ALTER TABLE users
ALTER COLUMN role_id SET NOT NULL;

-- Add foreign key constraint to ensure role_id references a valid role
ALTER TABLE users
ADD CONSTRAINT fk_role
FOREIGN KEY (role_id) REFERENCES roles(id)
ON DELETE RESTRICT;

