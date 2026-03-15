package frameworks

import (
	"context"
	"log/slog"

	"github.com/apps-deployer/projects-service/internal/domain/models"
	"github.com/apps-deployer/projects-service/internal/lib/logger/sl"
)

type FrameworkStorage interface {
	Framework(ctx context.Context, id string) (*models.Framework, error)
	ListFrameworks(ctx context.Context, args *models.ListFrameworksParams) ([]*models.Framework, error)
	SaveFramework(ctx context.Context, args *models.CreateFrameworkParams) (*models.SaveFrameworkResponse, error)
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
	op := "Frameworks.Get"
	log := f.log.With(
		slog.String("op", op),
		slog.String("id", id),
	)
	log.Info("getting framework")
	framework, err := f.frameworks.Framework(ctx, id)
	if err != nil {
		log.Error("failed to get framework", sl.Err(err))
		return nil, err
	}
	return framework, nil
}

func (f *Frameworks) List(ctx context.Context, args *models.ListFrameworksParams) ([]*models.Framework, error) {
	op := "Frameworks.List"
	log := f.log.With(
		slog.String("op", op),
	)
	log.Info("listing frameworks")
	frameworks, err := f.frameworks.ListFrameworks(ctx, args)
	if err != nil {
		log.Error("failed to list frameworks", sl.Err(err))
		return nil, err
	}
	return frameworks, nil
}

func (f *Frameworks) Create(ctx context.Context, args *models.CreateFrameworkParams) (*models.Framework, error) {
	// TODO: check admin
	op := "Frameworks.Create"
	log := f.log.With(
		slog.String("op", op),
		slog.String("name", args.Name),
	)
	log.Info("creating framework")
	res, err := f.frameworks.SaveFramework(ctx, args)
	if err != nil {
		log.Error("failed to create framework", sl.Err(err))
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
		CreatedAt:  res.CreatedAt,
		UpdatedAt:  res.UpdatedAt,
	}, nil
}

func (f *Frameworks) Update(ctx context.Context, args *models.UpdateFrameworkParams) error {
	// TODO: check admin
	op := "Frameworks.Update"
	log := f.log.With(
		slog.String("op", op),
		slog.String("id", args.Id),
	)
	log.Info("updating framework")
	err := f.frameworks.UpdateFramework(ctx, args)
	if err != nil {
		log.Error("failed to update framework", sl.Err(err))
		return err
	}
	return nil
}

func (f *Frameworks) Delete(ctx context.Context, id string) error {
	// TODO: check admin
	op := "Frameworks.Delete"
	log := f.log.With(
		slog.String("op", op),
		slog.String("id", id),
	)
	log.Info("deleting framework")
	err := f.frameworks.DeleteFramework(ctx, id)
	if err != nil {
		log.Error("failed to delete framework", sl.Err(err))
		return err
	}
	return nil
}
