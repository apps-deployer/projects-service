ALTER TABLE projects.projects DROP CONSTRAINT IF EXISTS projects_owner_slug_key;
ALTER TABLE projects.projects DROP CONSTRAINT IF EXISTS projects_slug_not_empty;
ALTER TABLE projects.projects DROP COLUMN IF EXISTS slug;
