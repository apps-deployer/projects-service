package postgres

import (
	"context"

	"github.com/apps-deployer/projects-service/internal/domain/models"
)

func (r *Repo) DeployConfig(ctx context.Context, projectId string) (*models.DeployConfig, error) {
	query := `
		SELECT
			c.id,
			c.project_id,
			c.framework_id,
			COALESCE(c.base_image_override, ''),
			COALESCE(c.root_dir_override, ''),
			COALESCE(c.output_dir_override, ''),
			COALESCE(c.install_cmd_override, ''),
			COALESCE(c.build_cmd_override, ''),
			COALESCE(c.run_cmd_override, ''),
			c.created_at,
			c.updated_at
		FROM projects.deploy_configs AS c
		WHERE c.project_id = $1
	`
	row := r.executor.QueryRow(ctx, query, projectId)
	var dc models.DeployConfig
	err := row.Scan(&dc.Id, &dc.ProjectId, &dc.FrameworkId, &dc.BaseImageOverride, &dc.RootDirOverride, &dc.OutputDirOverride, &dc.InstallCmdOverride, &dc.BuildCmdOverride, &dc.RunCmdOverride, &dc.CreatedAt, &dc.UpdatedAt)
	if err != nil {
		return nil, mapError(err)
	}
	return &dc, nil
}

func (r *Repo) UpdateDeployConfig(ctx context.Context, args *models.UpdateDeployConfigParams) error {
	query := `
		UPDATE projects.deploy_configs
		SET framework_id = COALESCE($2, framework_id),
		    base_image_override = COALESCE($3, base_image_override),
		    root_dir_override = COALESCE($4, root_dir_override),
		    output_dir_override = COALESCE($5, output_dir_override),
		    install_cmd_override = COALESCE($6, install_cmd_override),
		    build_cmd_override = COALESCE($7, build_cmd_override),
		    run_cmd_override = COALESCE($8, run_cmd_override)
		WHERE id = $1
	`
	_, err := r.executor.Exec(ctx, query, args.Id, args.FrameworkId, args.BaseImageOverride, args.RootDirOverride, args.OutputDirOverride, args.InstallCmdOverride, args.BuildCmdOverride, args.RunCmdOverride)
	return mapError(err)
}

func (r *Repo) SaveDeployConfig(ctx context.Context, args *models.SaveDeployConfigParams) (*models.SaveDeployConfigResponse, error) {
	query := `
		INSERT INTO projects.deploy_configs (project_id, framework_id)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`
	row := r.executor.QueryRow(ctx, query, args.ProjectId, args.FrameworkId)
	var res models.SaveDeployConfigResponse
	err := row.Scan(&res.Id, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		return nil, mapError(err)
	}
	return &res, nil
}

func (r *Repo) DeleteDeployConfig(ctx context.Context, projectId string) error {
	query := `DELETE FROM projects.deploy_configs WHERE project_id = $1`
	_, err := r.executor.Exec(ctx, query, projectId)
	return mapError(err)
}

func (r *Repo) ProjectOwnerByDeployConfigID(ctx context.Context, configID string) (string, error) {
	query := `
		SELECT p.owner_id
		FROM projects.deploy_configs dc
		JOIN projects.projects p ON dc.project_id = p.id
		WHERE dc.id = $1
	`
	var ownerID string
	err := r.executor.QueryRow(ctx, query, configID).Scan(&ownerID)
	if err != nil {
		return "", mapError(err)
	}
	return ownerID, nil
}
