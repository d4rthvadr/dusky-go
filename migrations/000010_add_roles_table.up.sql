CREATE TABLE if NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    level INT NOT NULL DEFAULT 0,
    description TEXT
);

INSERT INTO roles (name, level, description) VALUES
('user', 1, 'can create and manage their own content'),
('editor', 2, 'can update content created by users'),
('admin', 3, 'Administrator can manage all content and users');