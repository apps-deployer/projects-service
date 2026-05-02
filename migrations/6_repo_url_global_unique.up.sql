ALTER TABLE projects.projects DROP CONSTRAINT IF EXISTS projects_owner_repo_url_key;
ALTER TABLE projects.projects DROP CONSTRAINT IF EXISTS projects_repo_url_key;
ALTER TABLE projects.projects ADD CONSTRAINT projects_repo_url_key UNIQUE (repo_url);
