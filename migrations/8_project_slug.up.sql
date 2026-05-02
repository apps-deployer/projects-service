ALTER TABLE projects.projects
    ADD COLUMN IF NOT EXISTS slug VARCHAR(64);

UPDATE projects.projects
SET slug = lower(regexp_replace(regexp_replace(trim(name), '[^a-zA-Z0-9-]+', '-', 'g'), '-+', '-', 'g'))
WHERE slug IS NULL OR slug = '';

ALTER TABLE projects.projects
    ALTER COLUMN slug SET NOT NULL,
    ADD CONSTRAINT projects_slug_not_empty CHECK (char_length(slug) > 0);

ALTER TABLE projects.projects
    ADD CONSTRAINT projects_owner_slug_key UNIQUE (owner_id, slug);
