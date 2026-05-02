ALTER TABLE projects.deploy_configs DROP COLUMN IF EXISTS app_port_override;
ALTER TABLE projects.frameworks DROP COLUMN IF EXISTS app_port;
