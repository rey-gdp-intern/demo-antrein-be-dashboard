package auth

import (
	"antrein/bc-dashboard/internal/repository/tenant"
	"antrein/bc-dashboard/internal/utils/generator"
	"antrein/bc-dashboard/model/config"
	"antrein/bc-dashboard/model/dto"
	"antrein/bc-dashboard/model/entity"
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type Usecase struct {
	cfg  *config.Config
	repo *tenant.Repository
}

func New(cfg *config.Config, repo *tenant.Repository) *Usecase {
	return &Usecase{
		cfg:  cfg,
		repo: repo,
	}
}
func (u *Usecase) RegisterNewTenant(ctx context.Context, req dto.CreateTenantRequest) (*dto.CreateTenantResponse, *dto.ErrorResponse) {
	var errRes dto.ErrorResponse

	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		errRes = dto.ErrorResponse{
			Status: 500,
			Error:  "Gagal membuat akun tenant",
		}
		return nil, &errRes
	}

	tenant := entity.Tenant{
		Email:     req.Email,
		Name:      req.Name,
		Password:  string(encryptedPass),
		CreatedAt: time.Now(),
	}

	created, err := u.repo.CreateNewTenant(ctx, tenant)
	if err != nil {
		log.Println("Error gagal membuat akun", err)
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			errRes = dto.ErrorResponse{
				Status: 400,
				Error:  "Email telah terdaftar",
			}
			return nil, &errRes
		}
		errRes = dto.ErrorResponse{
			Status: 500,
			Error:  "Gagal membuat akun tenant",
		}
		return nil, &errRes
	}

	claims := entity.JWTClaim{
		UserID: created.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "rest",
			Subject:   "",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 5)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token, err := generator.GenerateJWTToken(u.cfg.Secrets.JWTSecret, claims)
	if err != nil {
		errRes = dto.ErrorResponse{
			Status: 500,
			Error:  "Gagal membuat akun tenant",
		}
		return nil, &errRes
	}

	return &dto.CreateTenantResponse{
		Tenant: dto.Tenant{
			ID:    created.ID,
			Name:  created.Name,
			Email: created.Email,
		},
		Token: token,
	}, nil
}

func (u *Usecase) LoginTenantAccount(ctx context.Context, req dto.LoginRequest) (*dto.CreateTenantResponse, *dto.ErrorResponse) {
	var errRes dto.ErrorResponse

	tenant, err := u.repo.GetTenantByEmail(ctx, req.Email)
	if tenant == nil {
		if err == sql.ErrNoRows {
			errRes = dto.ErrorResponse{
				Status: 401,
				Error:  "Email atau password salah",
			}
			return nil, &errRes
		}
		log.Println("Error login akun", err)
		errRes = dto.ErrorResponse{
			Status: 500,
			Error:  "Gagal login ke akun",
		}
		return nil, &errRes
	}
	err = bcrypt.CompareHashAndPassword([]byte(tenant.Password), []byte(req.Password))
	if err != nil {
		errRes = dto.ErrorResponse{
			Status: 401,
			Error:  "Email atau password salah",
		}
		return nil, &errRes
	}

	claims := entity.JWTClaim{
		UserID: tenant.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "rest",
			Subject:   "",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 5)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token, err := generator.GenerateJWTToken(u.cfg.Secrets.JWTSecret, claims)
	if err != nil {
		errRes = dto.ErrorResponse{
			Status: 500,
			Error:  "Gagal login ke akun",
		}
		return nil, &errRes
	}

	return &dto.CreateTenantResponse{
		Tenant: dto.Tenant{
			ID:    tenant.ID,
			Name:  tenant.Name,
			Email: tenant.Email,
		},
		Token: token,
	}, nil
}
