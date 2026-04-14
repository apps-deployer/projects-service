package postgres

import (
	"context"

	"github.com/apps-deployer/projects-service/internal/domain/models"
)

func (r *Repo) ListEnvVars(ctx context.Context, args *models.ListEnvVarsParams) ([]*models.Var, error) {
	query := `
		SELECT id, key, created_at, updated_at
		FROM projects.env_vars
		WHERE env_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.executor.Query(ctx, query, args.EnvId, args.Limit, args.Offset)
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

func (r *Repo) SaveEnvVar(ctx context.Context, args *models.CreateEnvVarParams) (*models.SaveVarResponse, error) {
	query := `
		INSERT INTO projects.env_vars (env_id, key, value)
		VALUES ($1, $2, crypto.pgp_sym_encrypt($3::text, $4))
		RETURNING id, created_at, updated_at
	`
	row := r.executor.QueryRow(ctx, query, args.EnvId, args.Key, args.Value, r.encryptionKey)
	var res models.SaveVarResponse
	err := row.Scan(&res.Id, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		return nil, mapError(err)
	}
	return &res, nil
}

func (r *Repo) UpdateEnvVar(ctx context.Context, args *models.UpdateVarParams) error {
	query := `
		UPDATE projects.env_vars
		SET value = CASE WHEN $2::text IS NOT NULL THEN crypto.pgp_sym_encrypt($2::text, $3) ELSE value END
		WHERE id = $1
	`
	_, err := r.executor.Exec(ctx, query, args.Id, args.Value, r.encryptionKey)
	return mapError(err)
}

func (r *Repo) DeleteEnvVar(ctx context.Context, id string) error {
	query := `DELETE FROM projects.env_vars WHERE id = $1`
	_, err := r.executor.Exec(ctx, query, id)
	return mapError(err)
}

func (r *Repo) ProjectOwnerByEnvVarID(ctx context.Context, varID string) (string, error) {
	query := `
		SELECT p.owner_id
		FROM projects.env_vars ev
		JOIN projects.envs e ON ev.env_id = e.id
		JOIN projects.projects p ON e.project_id = p.id
		WHERE ev.id = $1
	`
	var ownerID string
	err := r.executor.QueryRow(ctx, query, varID).Scan(&ownerID)
	if err != nil {
		return "", mapError(err)
	}
	return ownerID, nil
}
