package vars

import (
	"log/slog"

	"github.com/apps-deployer/projects-service/internal/services"
)

type Vars struct {
	log      *slog.Logger
	pv       services.ProjectVarRepository
	ev       services.EnvVarRepository
	rv       services.ResolvedVarsRepository
	projects services.ProjectRepository
	envs     services.EnvRepository
}

func New(log *slog.Logger, storage services.Storage) *Vars {
	repos := storage.Repos()
	return &Vars{
		log:      log,
		pv:       repos.ProjectVars(),
		ev:       repos.EnvVars(),
		rv:       repos.ResolvedVars(),
		projects: repos.Projects(),
		envs:     repos.Envs(),
	}
}
