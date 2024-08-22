package project

import (
	"antrein/bc-dashboard/internal/repository/infra"
	"antrein/bc-dashboard/internal/repository/project"
	"antrein/bc-dashboard/model/config"
	"antrein/bc-dashboard/model/dto"
	"antrein/bc-dashboard/model/entity"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/lib/pq"
)

type Usecase struct {
	cfg       *config.Config
	repo      *project.Repository
	infraRepo *infra.Repository
}

func New(cfg *config.Config, repo *project.Repository, infraRepo *infra.Repository) *Usecase {
	return &Usecase{
		cfg:       cfg,
		repo:      repo,
		infraRepo: infraRepo,
	}
}

func (u *Usecase) RegisterNewProject(ctx context.Context, req dto.CreateProjectRequest, tenantID string) (*dto.CreateProjectResponse, *dto.ErrorResponse) {
	var errRes dto.ErrorResponse

	project := entity.Project{
		ID:        req.ID,
		Name:      req.Name,
		TenantID:  tenantID,
		CreatedAt: time.Now(),
	}

	created, err := u.repo.CreateNewProject(ctx, project)
	if err != nil {
		log.Println("Error gagal membuat project", err)
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			errRes = dto.ErrorResponse{
				Status: 400,
				Error:  "Project dengan id tersebut sudah ada",
			}
			return nil, &errRes
		}
		errRes = dto.ErrorResponse{
			Status: 500,
			Error:  "Gagal membuat project",
		}
		return nil, &errRes
	}

	return &dto.CreateProjectResponse{
		Project: dto.Project{
			ID:       created.ID,
			Name:     created.Name,
			TenantID: created.TenantID,
		},
	}, nil
}

func (u *Usecase) GetListProject(ctx context.Context, tenantID string) (*dto.ListProjectResponse, *dto.ErrorResponse) {
	var errRes dto.ErrorResponse
	projects, err := u.repo.GetTenantProjects(ctx, tenantID)
	if err != nil {
		errRes = dto.ErrorResponse{
			Status: 500,
			Error:  err.Error(),
		}
		return nil, &errRes
	}
	listProjects := make([]dto.Project, len(projects))
	for i, project := range projects {
		listProjects[i] = dto.Project{
			ID:       project.ID,
			Name:     project.Name,
			TenantID: tenantID,
		}
	}
	return &dto.ListProjectResponse{
		TenantID: tenantID,
		Projects: listProjects,
	}, nil
}

func (u *Usecase) GetProjectDetail(ctx context.Context, projectID, tenantID string) (*dto.ProjectDetailResponse, *dto.ErrorResponse) {
	var errRes dto.ErrorResponse
	project, err := u.repo.GetTenantProjectByID(ctx, projectID, tenantID)
	if err != nil {
		errRes = dto.ErrorResponse{
			Status: 500,
			Error:  err.Error(),
		}
		return nil, &errRes
	}
	return &dto.ProjectDetailResponse{
		ID:       projectID,
		Name:     project.Name,
		TenantID: project.TenantID,
		Configuration: dto.ProjectConfig{
			ProjectID:          projectID,
			Threshold:          project.Threshold,
			SessionTime:        project.SessionTime,
			Host:               project.Host.String,
			BaseURL:            project.BaseURL.String,
			MaxUsersInQueue:    project.MaxUsersInQueue,
			QueueStart:         project.QueueEnd.Time,
			QueueEnd:           project.QueueEnd.Time,
			QueuePageStyle:     project.QueuePageStyle,
			QueueHTMLPage:      project.QueueHTMLPage.String,
			QueuePageBaseColor: project.QueuePageBaseColor.String,
			QueuePageTitle:     project.QueuePageTitle.String,
			QueuePageLogo:      project.QueuePageLogo.String,
			IsConfigure:        project.IsConfigure,
		},
	}, nil
}

func (u *Usecase) CheckHealthProject(ctx context.Context, projectID string) (*dto.CheckHealthProjectResponse, *dto.ErrorResponse) {
	var errRes dto.ErrorResponse
	client := &http.Client{}
	healthiness, err := u.infraRepo.CheckHealthProject(client, projectID)
	if err != nil {
		errRes = dto.ErrorResponse{
			Status: 500,
			Error:  err.Error(),
		}
		return nil, &errRes
	}
	return &dto.CheckHealthProjectResponse{
		ID:          projectID,
		Healthiness: healthiness,
	}, nil
}

func (u *Usecase) ClearProject(ctx context.Context) *dto.ErrorResponse {
	var errRes dto.ErrorResponse
	err := u.repo.ClearAllProjects(ctx)
	if err != nil {
		errRes = dto.ErrorResponse{
			Status: 500,
			Error:  err.Error(),
		}
		return &errRes
	}
	return nil
}
