package repository

import (
	"antrein/bc-dashboard/application/common/resource"
	"antrein/bc-dashboard/internal/repository/configuration"
	"antrein/bc-dashboard/internal/repository/infra"
	"antrein/bc-dashboard/internal/repository/project"
	"antrein/bc-dashboard/internal/repository/tenant"
	"antrein/bc-dashboard/model/config"
)

type CommonRepository struct {
	TenantRepo  *tenant.Repository
	ProjectRepo *project.Repository
	ConfigRepo  *configuration.Repository
	InfraRepo   *infra.Repository
}

func NewCommonRepository(cfg *config.Config, rsc *resource.CommonResource) (*CommonRepository, error) {
	tenantRepo := tenant.New(cfg, rsc.Db)
	infraRepo := infra.New(cfg)
	projectRepo := project.New(cfg, rsc.Db, infraRepo)
	configRepo := configuration.New(cfg, rsc.Db, infraRepo)

	commonRepo := CommonRepository{
		TenantRepo:  tenantRepo,
		ProjectRepo: projectRepo,
		ConfigRepo:  configRepo,
		InfraRepo:   infraRepo,
	}
	return &commonRepo, nil
}
