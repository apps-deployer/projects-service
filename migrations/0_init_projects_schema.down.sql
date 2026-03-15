DELETE FROM projects.frameworks WHERE name = 'Custom' AND base_image = 'scratch';

DROP TRIGGER IF EXISTS trg_update_projects ON projects.projects;
DROP TRIGGER IF EXISTS trg_update_frameworks ON projects.frameworks;
DROP TRIGGER IF EXISTS trg_update_deploy_configs ON projects.deploy_configs;
DROP TRIGGER IF EXISTS trg_update_envs ON projects.envs;
DROP TRIGGER IF EXISTS trg_update_project_vars ON projects.project_vars;
DROP TRIGGER IF EXISTS trg_update_env_vars ON projects.env_vars;

DROP FUNCTION IF EXISTS utils.update_updated_at();

DROP TABLE IF EXISTS projects.env_vars CASCADE;
DROP TABLE IF EXISTS projects.project_vars CASCADE;
DROP TABLE IF EXISTS projects.envs CASCADE;
DROP TABLE IF EXISTS projects.deploy_configs CASCADE;
DROP TABLE IF EXISTS projects.frameworks CASCADE;
DROP TABLE IF EXISTS projects.projects CASCADE;

DROP SCHEMA IF EXISTS utils CASCADE;
DROP SCHEMA IF EXISTS projects CASCADE;
DROP SCHEMA IF EXISTS crypto CASCADE;

