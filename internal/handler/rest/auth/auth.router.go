package auth

import (
	guard "antrein/bc-dashboard/application/middleware"
	"antrein/bc-dashboard/internal/usecase/auth"
	validate "antrein/bc-dashboard/internal/utils/validator"
	"antrein/bc-dashboard/model/config"
	"antrein/bc-dashboard/model/dto"
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Router struct {
	cfg     *config.Config
	usecase *auth.Usecase
	vld     *validator.Validate
}

func New(cfg *config.Config, usecase *auth.Usecase, vld *validator.Validate) *Router {
	return &Router{
		cfg:     cfg,
		usecase: usecase,
		vld:     vld,
	}
}

func (r *Router) RegisterRoute(app *mux.Router) {
	app.HandleFunc("/bc/dashboard/auth/register", guard.DefaultGuard(r.RegisterTenant))
	app.HandleFunc("/bc/dashboard/auth/login", guard.DefaultGuard(r.LoginTenantAccount))
}

func (r *Router) RegisterTenant(g *guard.GuardContext) error {
	ok := guard.IsMethod(g.Request, "POST")
	if !ok {
		return g.ReturnError(http.StatusMethodNotAllowed, "Method not allowed")
	}
	req := dto.CreateTenantRequest{}

	err := guard.BodyParser(g.Request, &req)
	if err != nil {
		return g.ReturnError(http.StatusBadRequest, "Request tidak sesuai format")
	}

	ctx := context.Background()

	err = r.vld.StructCtx(ctx, &req)
	if err != nil {
		return g.ReturnError(http.StatusBadRequest, "Request tidak sesuai format")
	}

	err = validate.ValidateCreateAccount(req)
	if err != nil {
		return g.ReturnError(http.StatusBadRequest, err.Error())
	}

	resp, errRes := r.usecase.RegisterNewTenant(ctx, req)
	if errRes != nil {
		return g.ReturnError(errRes.Status, errRes.Error)
	}

	return g.ReturnCreated(resp)
}

func (r *Router) LoginTenantAccount(g *guard.GuardContext) error {
	ok := guard.IsMethod(g.Request, "POST")
	if !ok {
		return g.ReturnError(http.StatusMethodNotAllowed, "Method not allowed")
	}
	req := dto.LoginRequest{}

	err := guard.BodyParser(g.Request, &req)
	if err != nil {
		return g.ReturnError(http.StatusBadRequest, "Request tidak sesuai format")
	}

	ctx := context.Background()

	err = r.vld.StructCtx(ctx, &req)
	if err != nil {
		return g.ReturnError(http.StatusBadRequest, "Request tidak sesuai format")
	}

	resp, errRes := r.usecase.LoginTenantAccount(ctx, req)
	if errRes != nil {
		return g.ReturnError(errRes.Status, errRes.Error)
	}

	return g.ReturnSuccess(resp)
}
