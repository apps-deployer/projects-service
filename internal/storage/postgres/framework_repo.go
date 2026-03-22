package postgres

import (
	"context"

	"github.com/apps-deployer/projects-service/internal/domain/models"
)

func (r *Repo) Framework(ctx context.Context, id string) (*models.Framework, error) {
	query := `
		SELECT id, name, base_image, root_dir, output_dir, install_cmd, build_cmd, run_cmd, created_at, updated_at
		FROM projects.frameworks
		WHERE id = $1
	`
	row := r.executor.QueryRow(ctx, query, id)
	var f models.Framework
	err := row.Scan(&f.Id, &f.Name, &f.BaseImage, &f.RootDir, &f.OutputDir, &f.InstallCmd, &f.BuildCmd, &f.RunCmd, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, mapError(err)
	}
	return &f, nil
}

func (r *Repo) ListFrameworks(ctx context.Context, args *models.ListFrameworksParams) ([]*models.Framework, error) {
	query := `
		SELECT id, name, base_image, root_dir, output_dir, install_cmd, build_cmd, run_cmd, created_at, updated_at
		FROM projects.frameworks
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.executor.Query(ctx, query, args.Limit, args.Offset)
	if err != nil {
		return nil, mapError(err)
	}
	defer rows.Close()
	var frameworks []*models.Framework
	for rows.Next() {
		var f models.Framework
		err := rows.Scan(&f.Id, &f.Name, &f.BaseImage, &f.RootDir, &f.OutputDir, &f.InstallCmd, &f.BuildCmd, &f.RunCmd, &f.CreatedAt, &f.UpdatedAt)
		if err != nil {
			return nil, err
		}
		frameworks = append(frameworks, &f)
	}
	return frameworks, rows.Err()
}

func (r *Repo) SaveFramework(ctx context.Context, args *models.CreateFrameworkParams) (*models.SaveFrameworkResponse, error) {
	query := `
		INSERT INTO projects.frameworks (name, base_image, root_dir, output_dir, install_cmd, build_cmd, run_cmd)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	row := r.executor.QueryRow(ctx, query, args.Name, args.BaseImage, args.RootDir, args.OutputDir, args.InstallCmd, args.BuildCmd, args.RunCmd)
	var res models.SaveFrameworkResponse
	err := row.Scan(&res.Id, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		return nil, mapError(err)
	}
	return &res, nil
}

func (r *Repo) UpdateFramework(ctx context.Context, args *models.UpdateFrameworkParams) error {
	query := `
		UPDATE projects.frameworks
		SET name = COALESCE($2, name),
		    base_image = COALESCE($3, base_image),
		    root_dir = COALESCE($4, root_dir),
		    output_dir = COALESCE($5, output_dir),
		    install_cmd = COALESCE($6, install_cmd),
		    build_cmd = COALESCE($7, build_cmd),
		    run_cmd = COALESCE($8, run_cmd)
		WHERE id = $1
	`
	_, err := r.executor.Exec(ctx, query, args.Id, args.Name, args.BaseImage, args.RootDir, args.OutputDir, args.InstallCmd, args.BuildCmd, args.RunCmd)
	return mapError(err)
}

func (r *Repo) DeleteFramework(ctx context.Context, id string) error {
	query := `DELETE FROM projects.frameworks WHERE id = $1`
	_, err := r.executor.Exec(ctx, query, id)
	return mapError(err)
}
