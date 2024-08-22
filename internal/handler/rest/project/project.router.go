package project

import (
	guard "antrein/bc-dashboard/application/middleware"
	"antrein/bc-dashboard/internal/usecase/configuration"
	"antrein/bc-dashboard/internal/usecase/project"
	validate "antrein/bc-dashboard/internal/utils/validator"
	"antrein/bc-dashboard/model/config"
	"antrein/bc-dashboard/model/dto"
	"context"
	"mime/multipart"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Router struct {
	cfg           *config.Config
	usecase       *project.Usecase
	configUsecase *configuration.Usecase
	vld           *validator.Validate
}

func New(cfg *config.Config, usecase *project.Usecase, configUsecase *configuration.Usecase, vld *validator.Validate) *Router {
	return &Router{
		cfg:           cfg,
		usecase:       usecase,
		configUsecase: configUsecase,
		vld:           vld,
	}
}

func (r *Router) RegisterRoute(app *mux.Router) {
	app.HandleFunc("/bc/dashboard/project/list", guard.AuthGuard(r.cfg, r.GetListProjects))
	app.HandleFunc("/bc/dashboard/project/health/{id}", guard.AuthGuard(r.cfg, r.CheckHealthProject))
	app.HandleFunc("/bc/dashboard/project/detail/{id}", guard.AuthGuard(r.cfg, r.GetProjectDetail))
	app.HandleFunc("/bc/dashboard/project", guard.AuthGuard(r.cfg, r.CreateProject))
	app.HandleFunc("/bc/dashboard/project/config", guard.AuthGuard(r.cfg, r.UpdateProjectConfig))
	app.HandleFunc("/bc/dashboard/project/style", guard.AuthGuard(r.cfg, r.UpdateProjectStyle))
	app.HandleFunc("/bc/dashboard/project/clear", guard.DefaultGuard(r.ClearAllProjects))
}

func (r *Router) CreateProject(g *guard.AuthGuardContext) error {
	ok := guard.IsMethod(g.Request, "POST")
	if !ok {
		return g.ReturnError(http.StatusMethodNotAllowed, "Method not allowed")
	}

	req := dto.CreateProjectRequest{}

	err := guard.BodyParser(g.Request, &req)
	if err != nil {
		return g.ReturnError(http.StatusBadRequest, "Request tidak sesuai format")
	}

	ctx := context.Background()

	err = r.vld.StructCtx(ctx, &req)
	if err != nil {
		return g.ReturnError(http.StatusBadRequest, "Request tidak sesuai format")
	}

	err = validate.ValidateCreateProject(req)
	if err != nil {
		return g.ReturnError(http.StatusBadRequest, err.Error())
	}

	userID := g.Claims.UserID
	resp, errRes := r.usecase.RegisterNewProject(ctx, req, userID)
	if errRes != nil {
		return g.ReturnError(errRes.Status, errRes.Error)
	}

	return g.ReturnCreated(resp)
}

func (r *Router) UpdateProjectConfig(g *guard.AuthGuardContext) error {
	ok := guard.IsMethod(g.Request, "PUT")
	if !ok {
		return g.ReturnError(http.StatusMethodNotAllowed, "Method not allowed")
	}

	req := dto.UpdateProjectConfig{}

	err := guard.BodyParser(g.Request, &req)
	if err != nil {
		return g.ReturnError(http.StatusBadRequest, "Request tidak sesuai format")
	}

	ctx := context.Background()

	err = r.vld.StructCtx(ctx, &req)
	if err != nil {
		return g.ReturnError(http.StatusBadRequest, "Request tidak sesuai format")
	}

	errRes := r.configUsecase.UpdateProjectConfig(ctx, req)
	if errRes != nil {
		return g.ReturnError(errRes.Status, errRes.Error)
	}

	return g.ReturnSuccess("Berhasil mengupdate konfigurasi project")
}

func (r *Router) UpdateProjectStyle(g *guard.AuthGuardContext) error {
	ok := guard.IsMethod(g.Request, "PUT")
	if !ok {
		return g.ReturnError(http.StatusMethodNotAllowed, "Method not allowed")
	}

	req := dto.UpdateProjectStyle{}

	err := g.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		return g.ReturnError(http.StatusBadRequest, "Request tidak sesuai format")
	}

	form := g.Request.MultipartForm

	if val, ok := form.Value["project_id"]; ok && len(val) > 0 {
		req.ProjectID = val[0]
	}

	if val, ok := form.Value["queue_page_style"]; ok && len(val) > 0 {
		req.QueuePageStyle = val[0]
	}

	if val, ok := form.Value["queue_page_base_color"]; ok && len(val) > 0 {
		req.QueuePageBaseColor = val[0]
	}

	if val, ok := form.Value["queue_page_title"]; ok && len(val) > 0 {
		req.QueuePageTitle = val[0]
	}

	ctx := context.Background()

	err = r.vld.StructCtx(ctx, &req)
	if err != nil {
		return g.ReturnError(http.StatusBadRequest, "Request tidak sesuai format")
	}

	var imageFile *multipart.FileHeader
	_, imageFile, err = g.Request.FormFile("image")
	if err != nil {
		if err != http.ErrMissingFile {
			return g.ReturnError(http.StatusBadRequest, "Gagal mendapatkan file logo")
		}
		imageFile = nil
	}

	var htmlFile *multipart.FileHeader
	_, htmlFile, err = g.Request.FormFile("file")
	if err != nil {
		if err != http.ErrMissingFile {
			return g.ReturnError(http.StatusBadRequest, "Gagal mendapatkan file html")
		}
		htmlFile = nil
	}

	errRes := r.configUsecase.UpdateProjectStyle(ctx, req, imageFile, htmlFile)
	if errRes != nil {
		return g.ReturnError(errRes.Status, errRes.Error)
	}

	return g.ReturnSuccess("Berhasil mengupdate tampilan project")
}

func (r *Router) GetListProjects(g *guard.AuthGuardContext) error {
	ok := guard.IsMethod(g.Request, "GET")
	if !ok {
		return g.ReturnError(http.StatusMethodNotAllowed, "Method not allowed")
	}

	ctx := context.Background()
	tenantID := g.Claims.UserID
	resp, errRes := r.usecase.GetListProject(ctx, tenantID)
	if errRes != nil {
		return g.ReturnError(errRes.Status, errRes.Error)
	}

	return g.ReturnSuccess(resp)
}

func (r *Router) GetProjectDetail(g *guard.AuthGuardContext) error {
	ok := guard.IsMethod(g.Request, "GET")
	if !ok {
		return g.ReturnError(http.StatusMethodNotAllowed, "Method not allowed")
	}

	projectID := guard.GetParam(g.Request, "id")
	ctx := context.Background()
	tenantID := g.Claims.UserID
	resp, errRes := r.usecase.GetProjectDetail(ctx, projectID, tenantID)
	if errRes != nil {
		return g.ReturnError(errRes.Status, errRes.Error)
	}

	return g.ReturnSuccess(resp)
}

func (r *Router) CheckHealthProject(g *guard.AuthGuardContext) error {
	ok := guard.IsMethod(g.Request, "GET")
	if !ok {
		return g.ReturnError(http.StatusMethodNotAllowed, "Method not allowed")
	}

	projectID := guard.GetParam(g.Request, "id")
	ctx := context.Background()
	resp, errRes := r.usecase.CheckHealthProject(ctx, projectID)
	if errRes != nil {
		return g.ReturnError(errRes.Status, errRes.Error)
	}

	return g.ReturnSuccess(resp)
}

func (r *Router) ClearAllProjects(g *guard.GuardContext) error {
	ok := guard.IsMethod(g.Request, "DELETE")
	if !ok {
		return g.ReturnError(http.StatusMethodNotAllowed, "Method not allowed")
	}

	ctx := context.Background()
	errRes := r.usecase.ClearProject(ctx)
	if errRes != nil {
		return g.ReturnError(errRes.Status, errRes.Error)
	}

	return g.ReturnSuccess("Berhasil clear semua project")
}
