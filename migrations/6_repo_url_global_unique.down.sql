ALTER TABLE projects.projects DROP CONSTRAINT IF EXISTS projects_repo_url_key;
ALTER TABLE projects.projects ADD CONSTRAINT projects_owner_repo_url_key UNIQUE (owner_id, repo_url);
