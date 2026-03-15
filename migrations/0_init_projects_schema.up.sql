CREATE SCHEMA IF NOT EXISTS crypto;
CREATE SCHEMA IF NOT EXISTS projects;
CREATE SCHEMA IF NOT EXISTS utils;

CREATE EXTENSION IF NOT EXISTS pgcrypto SCHEMA crypto;

CREATE TABLE IF NOT EXISTS projects.projects (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    name VARCHAR(128) NOT NULL CHECK char_length(name) > 0,
    owner_id UUID NOT NULL,
    repo_url VARCHAR(512) NOT NULL UNIQUE CHECK char_length(repo_url) > 0,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE(owner_id, name)
);

CREATE TABLE IF NOT EXISTS projects.frameworks (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    name VARCHAR(128) UNIQUE NOT NULL CHECK char_length(name) > 0,
    base_image VARCHAR(128) NOT NULL CHECK char_length(base_image) > 0,
    root_dir VARCHAR(128),
    output_dir VARCHAR(128),
    install_cmd VARCHAR(128),
    build_cmd VARCHAR(128),
    run_cmd VARCHAR(128),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS projects.deploy_configs (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    project_id UUID UNIQUE NOT NULL REFERENCES projects.projects (id) ON DELETE CASCADE,
    framework_id UUID NOT NULL REFERENCES projects.frameworks (id) ON DELETE RESTRICT,
    base_image_override VARCHAR(128),
    root_dir_override VARCHAR(128),
    output_dir_override VARCHAR(128),
    install_cmd_override VARCHAR(128),
    build_cmd_override VARCHAR(128),
    run_cmd_override VARCHAR(128),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS projects.envs (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    name VARCHAR(128) NOT NULL CHECK char_length(name) > 0,
    project_id UUID NOT NULL REFERENCES projects.projects (id) ON DELETE CASCADE,
    target_branch VARCHAR(128) NOT NULL CHECK char_length(target_branch) > 0,
    domain_name VARCHAR(128) UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE(project_id, name),
    UNIQUE(project_id, target_branch)
);

CREATE TABLE IF NOT EXISTS projects.project_vars (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    project_id UUID NOT NULL REFERENCES projects.projects (id) ON DELETE CASCADE,
    key VARCHAR(128) NOT NULL CHECK char_length(key) > 0,
    value BYTEA NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE(project_id, key)

);

CREATE TABLE IF NOT EXISTS projects.env_vars (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    env_id UUID NOT NULL REFERENCES projects.envs (id) ON DELETE CASCADE,
    key VARCHAR(128) NOT NULL CHECK char_length(key) > 0,
    value BYTEA NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE(env_id, key)
);

CREATE INDEX IF NOT EXISTS idx_envs_project_id ON projects.envs (project_id);
CREATE INDEX IF NOT EXISTS idx_project_vars_project_id ON projects.project_vars (project_id);
CREATE INDEX IF NOT EXISTS idx_env_vars_env_id ON projects.env_vars (env_id);

CREATE OR REPLACE FUNCTION utils.update_updated_at()
RETURN TRIGGER AS $$
BEGIN 
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_projects
BEFORE UPDATE ON projects.projects
FOR EACH ROW
EXECUTE FUNCTION utils.update_updated_at();

CREATE TRIGGER trg_update_frameworks
BEFORE UPDATE ON projects.frameworks
FOR EACH ROW
EXECUTE FUNCTION utils.update_updated_at();

CREATE TRIGGER trg_update_deploy_configs
BEFORE UPDATE ON projects.deploy_configs
FOR EACH ROW
EXECUTE FUNCTION utils.update_updated_at();

CREATE TRIGGER trg_update_envs
BEFORE UPDATE ON projects.envs
FOR EACH ROW
EXECUTE FUNCTION utils.update_updated_at();

CREATE TRIGGER trg_update_project_vars
BEFORE UPDATE ON projects.project_vars
FOR EACH ROW
EXECUTE FUNCTION utils.update_updated_at();

CREATE TRIGGER trg_update_env_vars
BEFORE UPDATE ON projects.env_vars
FOR EACH ROW
EXECUTE FUNCTION utils.update_updated_at();

INSERT INTO projects.frameworks (name, base_image)
VALUES ('Custom', 'scratch')
ON CONFLICT (name) DO NOTHING;
