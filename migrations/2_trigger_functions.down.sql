DROP TRIGGER IF EXISTS trg_update_projects ON projects.projects;
DROP TRIGGER IF EXISTS trg_update_frameworks ON projects.frameworks;
DROP TRIGGER IF EXISTS trg_update_deploy_configs ON projects.deploy_configs;
DROP TRIGGER IF EXISTS trg_update_envs ON projects.envs;
DROP TRIGGER IF EXISTS trg_update_project_vars ON projects.project_vars;
DROP TRIGGER IF EXISTS trg_update_env_vars ON projects.env_vars;

DROP FUNCTION IF EXISTS utils.update_updated_at();
