package tenant

import (
	"antrein/bc-dashboard/model/config"
	"antrein/bc-dashboard/model/entity"
	"context"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	cfg *config.Config
	db  *sqlx.DB
}

func New(cfg *config.Config, db *sqlx.DB) *Repository {
	return &Repository{
		cfg: cfg,
		db:  db,
	}
}

func (r *Repository) CreateNewTenant(ctx context.Context, req entity.Tenant) (*entity.Tenant, error) {
	tenant := req
	q := `INSERT INTO tenants (email, password, name, created_at) VALUES ($1, $2, $3, $4) returning id`
	var id string
	err := r.db.GetContext(ctx, &id, q, req.Email, req.Password, req.Name, req.CreatedAt)
	tenant.ID = id
	return &tenant, err
}

func (r *Repository) GetTenantByID(ctx context.Context, id string) (*entity.Tenant, error) {
	tenant := entity.Tenant{}
	q := `SELECT * FROM tenants WHERE id = $1 LIMIT 1`
	err := r.db.GetContext(ctx, &tenant, q, id)
	if err != nil {
		return nil, err
	}
	return &tenant, err
}

func (r *Repository) GetTenantByEmail(ctx context.Context, email string) (*entity.Tenant, error) {
	tenant := entity.Tenant{}
	q := `SELECT * FROM tenants WHERE email = $1 LIMIT 1`
	err := r.db.GetContext(ctx, &tenant, q, email)
	if err != nil {
		return nil, err
	}
	return &tenant, err
}

func (r *Repository) GetTenants(ctx context.Context, page int, pageSize int) ([]entity.Tenant, error) {
	tenants := []entity.Tenant{}
	q := `SELECT * FROM tenants ORDER BY name LIMIT $1 OFFSET $2`
	offset := (page - 1) * pageSize
	err := r.db.SelectContext(ctx, &tenants, q, pageSize, offset)
	return tenants, err
}
