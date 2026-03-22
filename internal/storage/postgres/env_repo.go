package postgres

import (
	"context"

	"github.com/apps-deployer/projects-service/internal/domain/models"
)

func (r *Repo) Env(ctx context.Context, id string) (*models.Env, error) {
	query := `
		SELECT id, name, project_id, target_branch, domain_name, created_at, updated_at
		FROM projects.envs
		WHERE id = $1
	`
	row := r.executor.QueryRow(ctx, query, id)
	var e models.Env
	err := row.Scan(&e.Id, &e.Name, &e.ProjectId, &e.TargetBranch, &e.DomainName, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, mapError(err)
	}
	return &e, nil
}

func (r *Repo) EnvByGit(ctx context.Context, repoUrl string, branch string) (*models.Env, error) {
	query := `
		SELECT e.id, e.name, e.project_id, e.target_branch, e.domain_name, e.created_at, e.updated_at
		FROM projects.envs e
		JOIN projects.projects p ON e.project_id = p.id
		WHERE p.repo_url = $1 AND e.target_branch = $2
	`
	row := r.executor.QueryRow(ctx, query, repoUrl, branch)
	var e models.Env
	err := row.Scan(&e.Id, &e.Name, &e.ProjectId, &e.TargetBranch, &e.DomainName, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, mapError(err)
	}
	return &e, nil
}

func (r *Repo) ListEnvs(ctx context.Context, args *models.ListEnvsParams) ([]*models.Env, error) {
	query := `
		SELECT id, name, project_id, target_branch, domain_name, created_at, updated_at
		FROM projects.envs
		WHERE project_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.executor.Query(ctx, query, args.ProjectId, args.Limit, args.Offset)
	if err != nil {
		return nil, mapError(err)
	}
	defer rows.Close()
	var envs []*models.Env
	for rows.Next() {
		var e models.Env
		err := rows.Scan(&e.Id, &e.Name, &e.ProjectId, &e.TargetBranch, &e.DomainName, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			return nil, err
		}
		envs = append(envs, &e)
	}
	return envs, rows.Err()
}

func (r *Repo) SaveEnv(ctx context.Context, args *models.CreateEnvParams) (*models.SaveEnvResponse, error) {
	query := `
		INSERT INTO projects.envs (name, project_id, target_branch, domain_name)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	row := r.executor.QueryRow(ctx, query, args.Name, args.ProjectId, args.TargetBranch, args.DomainName)
	var res models.SaveEnvResponse
	err := row.Scan(&res.Id, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		return nil, mapError(err)
	}
	return &res, nil
}

func (r *Repo) UpdateEnv(ctx context.Context, args *models.UpdateEnvParams) error {
	query := `
		UPDATE projects.envs
		SET name = COALESCE($2, name),
		    target_branch = COALESCE($3, target_branch),
		    domain_name = COALESCE($4, domain_name)
		WHERE id = $1
	`
	_, err := r.executor.Exec(ctx, query, args.Id, args.Name, args.TargetBranch, args.DomainName)
	return mapError(err)
}

func (r *Repo) DeleteEnv(ctx context.Context, id string) error {
	query := `DELETE FROM projects.envs WHERE id = $1`
	_, err := r.executor.Exec(ctx, query, id)
	return mapError(err)
}
