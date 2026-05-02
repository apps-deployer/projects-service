ALTER TABLE projects.frameworks
    ADD COLUMN IF NOT EXISTS app_port INTEGER NOT NULL DEFAULT 8080
        CHECK (app_port BETWEEN 1 AND 65535);

ALTER TABLE projects.deploy_configs
    ADD COLUMN IF NOT EXISTS app_port_override INTEGER
        CHECK (app_port_override IS NULL OR app_port_override BETWEEN 0 AND 65535);
