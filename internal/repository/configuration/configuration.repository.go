package configuration

import (
	"antrein/bc-dashboard/internal/repository/infra"
	"antrein/bc-dashboard/model/config"
	"antrein/bc-dashboard/model/entity"
	"context"
	"database/sql"
	"errors"
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

func (r *Repository) GetConfigByProjectID(ctx context.Context, projectID string) (*entity.Configuration, error) {
	config := entity.Configuration{}
	q := `SELECT * FROM configurations WHERE project_id = $1 LIMIT 1`
	err := r.db.GetContext(ctx, &config, q, projectID)
	if err != nil {
		return nil, err
	}
	return &config, err
}

func (r *Repository) GetConfigByHost(ctx context.Context, host string) (*entity.Configuration, error) {
	config := entity.Configuration{}
	q := `SELECT * FROM configurations WHERE host = $1 LIMIT 1`
	err := r.db.GetContext(ctx, &config, q, host)
	if err != nil {
		return nil, err
	}
	return &config, err
}

func (r *Repository) UpdateProjectConfig(ctx context.Context, req entity.Configuration) error {
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: 1,
		ReadOnly:  false,
	})
	q := `UPDATE configurations 
		  SET threshold = $1, 
		  session_time = $2, 
		  host = $3, 
		  base_url = $4,
		  max_users_in_queue = $5,
		  queue_start = $6,
		  queue_end = $7,
		  is_configure = $8,
		  updated_at = now()
		  WHERE project_id = $9`

	resp, err := tx.ExecContext(ctx, q, req.Threshold, req.SessionTime, req.Host, req.BaseURL, req.MaxUsersInQueue, req.QueueStart, req.QueueEnd, true, req.ProjectID)

	if err != nil {
		tx.Rollback()
		return err
	}

	affected, err := resp.RowsAffected()

	if err != nil {
		tx.Rollback()
		return err
	}

	if affected == 0 {
		tx.Rollback()
		return errors.New("Project tidak terdaftar")
	}

	client := &http.Client{}

	err = r.infraRepo.CreateInfraProject(client, infra.InfraBody{
		ProjectID:     req.ProjectID,
		ProjectDomain: req.Host.String,
		URLPath:       req.BaseURL.String,
	})

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *Repository) UpdateProjectStyle(ctx context.Context, req entity.Configuration) error {
	q := `UPDATE configurations 
		  SET queue_page_style = $1,
		  queue_html_page = $2,
		  queue_page_base_color = $3,
		  queue_page_title = $4,
		  queue_page_logo = $5,
		  updated_at = now()
		  WHERE project_id = $6`
	_, err := r.db.ExecContext(ctx, q, req.QueuePageStyle, req.QueueHTMLPage, req.QueuePageBaseColor, req.QueuePageTitle, req.QueuePageLogo, req.ProjectID)
	return err
}
