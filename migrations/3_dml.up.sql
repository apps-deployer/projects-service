INSERT INTO projects.frameworks (name, base_image)
VALUES ('Custom', 'scratch')
ON CONFLICT (name) DO NOTHING;
