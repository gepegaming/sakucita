-- roles
CREATE TABLE roles (
  id    SERIAL PRIMARY KEY,
  name  VARCHAR(50) UNIQUE NOT NULL
);

-- seeder role
INSERT INTO roles (name) VALUES
  ('super_admin'),
  ('supporter'),
  ('creator')
ON CONFLICT DO NOTHING;