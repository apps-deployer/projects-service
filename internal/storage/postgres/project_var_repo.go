package postgres

import (
	"context"

	"github.com/apps-deployer/projects-service/internal/domain/models"
)

func (r *Repo) ListProjectVars(ctx context.Context, args *models.ListProjectVarsParams) ([]*models.Var, error) {
	query := `
		SELECT id, key, created_at, updated_at
		FROM projects.project_vars
		WHERE project_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.executor.Query(ctx, query, args.ProjectId, args.Limit, args.Offset)
	if err != nil {
		return nil, mapError(err)
	}
	defer rows.Close()
	var vars []*models.Var
	for rows.Next() {
		var v models.Var
		err := rows.Scan(&v.Id, &v.Key, &v.CreatedAt, &v.UpdatedAt)
		if err != nil {
			return nil, err
		}
		vars = append(vars, &v)
	}
	return vars, rows.Err()
}

func (r *Repo) SaveProjectVar(ctx context.Context, args *models.CreateProjectVarParams) (*models.SaveVarResponse, error) {
	query := `
		INSERT INTO projects.project_vars (project_id, key, value)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	row := r.executor.QueryRow(ctx, query, args.ProjectId, args.Key, []byte(args.Value))
	var res models.SaveVarResponse
	err := row.Scan(&res.Id, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		return nil, mapError(err)
	}
	return &res, nil
}

func (r *Repo) UpdateProjectVar(ctx context.Context, args *models.UpdateVarParams) error {
	query := `
		UPDATE projects.project_vars
		SET value = COALESCE($2, value)
		WHERE id = $1
	`
	var value []byte
	if args.Value != nil {
		value = []byte(*args.Value)
	}
	_, err := r.executor.Exec(ctx, query, args.Id, value)
	return mapError(err)
}

func (r *Repo) DeleteProjectVar(ctx context.Context, id string) error {
	query := `DELETE FROM projects.project_vars WHERE id = $1`
	_, err := r.executor.Exec(ctx, query, id)
	return mapError(err)
}
