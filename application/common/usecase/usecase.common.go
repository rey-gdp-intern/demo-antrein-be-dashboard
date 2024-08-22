package usecase

import (
	"antrein/bc-dashboard/application/common/repository"
	"antrein/bc-dashboard/internal/usecase/auth"
	"antrein/bc-dashboard/internal/usecase/configuration"
	"antrein/bc-dashboard/internal/usecase/project"
	"antrein/bc-dashboard/model/config"
)

type CommonUsecase struct {
	AuthUsecase    *auth.Usecase
	ProjectUsecase *project.Usecase
	ConfigUsecase  *configuration.Usecase
}

func NewCommonUsecase(cfg *config.Config, repo *repository.CommonRepository) (*CommonUsecase, error) {
	authUsecase := auth.New(cfg, repo.TenantRepo)
	configUsecase := configuration.New(cfg, repo.ConfigRepo, repo.InfraRepo)
	projectUsecase := project.New(cfg, repo.ProjectRepo, repo.InfraRepo)

	commonUC := CommonUsecase{
		AuthUsecase:    authUsecase,
		ProjectUsecase: projectUsecase,
		ConfigUsecase:  configUsecase,
	}
	return &commonUC, nil
}
