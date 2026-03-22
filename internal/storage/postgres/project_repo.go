package postgres

import (
	"context"

	"github.com/apps-deployer/projects-service/internal/domain/models"
)

func (r *Repo) Project(ctx context.Context, id string) (*models.Project, error) {
	query := `
		SELECT id, name, repo_url, owner_id, created_at, updated_at
		FROM projects.projects
		WHERE id = $1
	`
	row := r.executor.QueryRow(ctx, query, id)
	var p models.Project
	err := row.Scan(&p.Id, &p.Name, &p.RepoUrl, &p.OwnerId, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, mapError(err)
	}
	return &p, nil
}

func (r *Repo) ListProjects(ctx context.Context, args *models.ListProjectsParams) ([]*models.Project, error) {
	query := `
		SELECT id, name, repo_url, owner_id, created_at, updated_at
		FROM projects.projects
		WHERE owner_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.executor.Query(ctx, query, args.OwnerId, args.Limit, args.Offset)
	if err != nil {
		return nil, mapError(err)
	}
	defer rows.Close()
	var projects []*models.Project
	for rows.Next() {
		var p models.Project
		err := rows.Scan(&p.Id, &p.Name, &p.RepoUrl, &p.OwnerId, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		projects = append(projects, &p)
	}
	return projects, rows.Err()
}

func (r *Repo) SaveProject(ctx context.Context, args *models.SaveProjectParams) (*models.SaveProjectResponse, error) {
	query := `
		INSERT INTO projects.projects (name, repo_url, owner_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	row := r.executor.QueryRow(ctx, query, args.Name, args.RepoUrl, args.OwnerId)
	var res models.SaveProjectResponse
	err := row.Scan(&res.Id, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		return nil, mapError(err)
	}
	return &res, nil
}

func (r *Repo) UpdateProject(ctx context.Context, args *models.UpdateProjectParams) error {
	query := `
		UPDATE projects.projects
		SET name = COALESCE($2, name),
		    repo_url = COALESCE($3, repo_url),
		    owner_id = COALESCE($4, owner_id)
		WHERE id = $1
	`
	_, err := r.executor.Exec(ctx, query, args.Id, args.Name, args.RepoUrl, args.OwnerId)
	return mapError(err)
}

func (r *Repo) DeleteProject(ctx context.Context, id string) error {
	query := `DELETE FROM projects.projects WHERE id = $1`
	_, err := r.executor.Exec(ctx, query, id)
	return mapError(err)
}
