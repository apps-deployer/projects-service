package frameworks

import (
	"context"
	"log/slog"

	"github.com/apps-deployer/projects-service/internal/domain/models"
)

type FrameworkStorage interface {
	Framework(ctx context.Context, id string) (*models.Framework, error)
	ListFrameworks(ctx context.Context, args *models.ListFrameworksParams) ([]*models.Framework, error)
	SaveFramework(ctx context.Context, args *models.SaveFrameworkParams) (*models.SaveFrameworkResponse, error)
	UpdateFramework(ctx context.Context, args *models.UpdateFrameworkParams) error
	DeleteFramework(ctx context.Context, id string) error
}

type Frameworks struct {
	log        *slog.Logger
	frameworks FrameworkStorage
}

func New(log *slog.Logger, frameworks FrameworkStorage) *Frameworks {
	return &Frameworks{log: log, frameworks: frameworks}
}

func (f *Frameworks) Get(ctx context.Context, id string) (*models.Framework, error) {
	framework, err := f.frameworks.Framework(ctx, id)
	return framework, err
}

func (f *Frameworks) List(ctx context.Context, args *models.ListFrameworksParams) ([]*models.Framework, error) {
	frameworks, err := f.frameworks.ListFrameworks(ctx, args)
	return frameworks, err
}

func (f *Frameworks) Create(ctx context.Context, args *models.CreateFrameworkParams) (*models.Framework, error) {
	// TODO: check admin
	res, err := f.frameworks.SaveFramework(ctx, args)
	if err != nil {
		return nil, err
	}
	return &models.Framework{
		Id:         res.Id,
		Name:       args.Name,
		RootDir:    args.RootDir,
		OutputDir:  args.OutputDir,
		BaseImage:  args.BaseImage,
		InstallCmd: args.InstallCmd,
		BuildCmd:   args.BuildCmd,
		RunCmd:     args.RunCmd,
	}, nil
}
func (f *Frameworks) Update(ctx context.Context, args *models.UpdateFrameworkParams) error {
	// TODO: check admin
	err := f.frameworks.UpdateFramework(ctx, args)
	return err
}
func (f *Frameworks) Delete(ctx context.Context, id string) error {
	// TODO: check admin
	err := f.frameworks.DeleteFramework(ctx, id)
	return err
}
