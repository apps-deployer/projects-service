package postgres

import (
	"context"

	"github.com/apps-deployer/projects-service/internal/domain/models"
)

func (r *Repo) ResolvedVars(ctx context.Context, envId string) ([]*models.ResolvedVar, error) {
	query := `
		WITH env_vars AS (
			SELECT key, crypto.pgp_sym_decrypt(value, $2)::text AS value
			FROM projects.env_vars
			WHERE env_id = $1
		),
		project_vars AS (
			SELECT pv.key, crypto.pgp_sym_decrypt(pv.value, $2)::text AS value
			FROM projects.project_vars pv
			JOIN projects.envs e ON pv.project_id = e.project_id
			WHERE e.id = $1
		),
		combined AS (
			SELECT key, value FROM env_vars
			UNION ALL
			SELECT key, value FROM project_vars
			WHERE key NOT IN (SELECT key FROM env_vars)
		)
		SELECT key, value FROM combined
	`
	rows, err := r.executor.Query(ctx, query, envId, r.encryptionKey)
	if err != nil {
		return nil, mapError(err)
	}
	defer rows.Close()
	var vars []*models.ResolvedVar
	for rows.Next() {
		var key string
		var value string
		err := rows.Scan(&key, &value)
		if err != nil {
			return nil, err
		}
		vars = append(vars, &models.ResolvedVar{
			Key:   key,
			Value: value,
		})
	}
	return vars, rows.Err()
}
