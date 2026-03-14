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
