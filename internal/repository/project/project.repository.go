package project

import (
	"antrein/bc-dashboard/internal/repository/infra"
	"antrein/bc-dashboard/model/config"
	"antrein/bc-dashboard/model/entity"
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	cfg       *config.Config
	db        *sqlx.DB
	infraRepo *infra.Repository
}

func New(cfg *config.Config, db *sqlx.DB, infraRepo *infra.Repository) *Repository {
	return &Repository{
		cfg:       cfg,
		db:        db,
		infraRepo: infraRepo,
	}
}

func (r *Repository) CreateNewProject(ctx context.Context, req entity.Project) (*entity.Project, error) {
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: 2,
		ReadOnly:  false,
	})
	project := req
	q1 := `INSERT INTO projects (id, name, tenant_id, created_at) VALUES ($1, $2, $3, $4)`
	_, err = tx.ExecContext(ctx, q1, req.ID, req.Name, req.TenantID, req.CreatedAt)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	q2 := `INSERT INTO configurations (project_id) VALUES ($1)`
	_, err = tx.ExecContext(ctx, q2, req.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &project, err
}

func (r *Repository) GetTenantByID(ctx context.Context, id string) (*entity.Project, error) {
	project := entity.Project{}
	q := `SELECT * FROM projects WHERE id = $1 LIMIT 1`
	err := r.db.GetContext(ctx, &project, q, id)
	if err != nil {
		return nil, err
	}
	return &project, err
}

func (r *Repository) GetTenantProjectByID(ctx context.Context, id, tenantID string) (*entity.ProjectWithConfig, error) {
	project := entity.ProjectWithConfig{}
	q := `SELECT * FROM projects INNER JOIN configurations ON projects.id = configurations.project_id WHERE projects.id = $1 AND projects.tenant_id = $2 LIMIT 1`
	err := r.db.GetContext(ctx, &project, q, id, tenantID)
	if err != nil {
		return nil, err
	}
	return &project, err
}

func (r *Repository) GetTenantProjects(ctx context.Context, tenantID string) ([]entity.Project, error) {
	projects := []entity.Project{}
	q := `SELECT * FROM projects WHERE tenant_id = $1 ORDER BY id`
	err := r.db.SelectContext(ctx, &projects, q, tenantID)
	return projects, err
}

func (r *Repository) GetProjects(ctx context.Context, page int, pageSize int) ([]entity.Project, error) {
	projects := []entity.Project{}
	q := `SELECT * FROM projects ORDER BY name LIMIT $1 OFFSET $2`
	offset := (page - 1) * pageSize
	err := r.db.SelectContext(ctx, &projects, q, pageSize, offset)
	return projects, err
}

func (r *Repository) ClearAllProjects(ctx context.Context) error {
	// Start transaction
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault, // Using default isolation level for clarity
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}

	q1 := `TRUNCATE TABLE configurations`
	if _, err = tx.ExecContext(ctx, q1); err != nil {
		fmt.Println(err)
		tx.Rollback()
		return err
	}

	q2 := `DELETE FROM projects`
	if _, err = tx.ExecContext(ctx, q2); err != nil {
		fmt.Println(err)
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	client := &http.Client{}
	if err = r.infraRepo.ClearInfraProject(client); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
