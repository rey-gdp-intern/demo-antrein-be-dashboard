package configuration

import (
	"antrein/bc-dashboard/internal/repository/configuration"
	"antrein/bc-dashboard/internal/repository/infra"
	"antrein/bc-dashboard/model/config"
	"antrein/bc-dashboard/model/dto"
	"antrein/bc-dashboard/model/entity"
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

type Usecase struct {
	cfg       *config.Config
	repo      *configuration.Repository
	infraRepo *infra.Repository
}

func New(cfg *config.Config, repo *configuration.Repository, infraRepo *infra.Repository) *Usecase {
	return &Usecase{
		cfg:       cfg,
		repo:      repo,
		infraRepo: infraRepo,
	}
}

func loadDefaultHTML() (string, error) {
	filePath := "./files/templates/queue.html"
	htmlFile, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer htmlFile.Close()

	content, err := io.ReadAll(htmlFile)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func readFileContent(file *multipart.FileHeader) ([]byte, error) {
	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

func changeValueInHTML(html, name, value string) string {
	return strings.Replace(html, name, value, 1)
}

func handleError(status int, message string) *dto.ErrorResponse {
	return &dto.ErrorResponse{
		Status: status,
		Error:  message,
	}
}

func (u *Usecase) GetProjectConfigByID(ctx context.Context, projectID string) (*dto.ProjectConfig, *dto.ErrorResponse) {
	var errRes dto.ErrorResponse

	config, err := u.repo.GetConfigByProjectID(ctx, projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			errRes = dto.ErrorResponse{
				Status: 404,
				Error:  "Project dengan id tersebut tidak ditemukan",
			}
			return nil, &errRes
		}
		log.Println(err)
		errRes = dto.ErrorResponse{
			Status: 500,
			Error:  "Gagal mendapatkan konfigurasi project",
		}
		return nil, &errRes
	}
	return &dto.ProjectConfig{
		ProjectID:          config.ProjectID,
		Threshold:          config.Threshold,
		SessionTime:        config.SessionTime,
		Host:               config.Host.String,
		BaseURL:            config.BaseURL.String,
		MaxUsersInQueue:    config.MaxUsersInQueue,
		QueueStart:         config.QueueStart.Time,
		QueueEnd:           config.QueueEnd.Time,
		QueuePageStyle:     config.QueuePageStyle,
		QueueHTMLPage:      config.QueueHTMLPage.String,
		QueuePageBaseColor: config.QueuePageBaseColor.String,
		QueuePageTitle:     config.QueuePageTitle.String,
		QueuePageLogo:      config.QueuePageLogo.String,
	}, nil
}

func (u *Usecase) GetProjectConfigByHost(ctx context.Context, host string) (*dto.ProjectConfig, *dto.ErrorResponse) {
	var errRes dto.ErrorResponse

	config, err := u.repo.GetConfigByHost(ctx, host)
	if err != nil {
		if err == sql.ErrNoRows {
			errRes = dto.ErrorResponse{
				Status: 404,
				Error:  "Project dengan host tersebut tidak ditemukan",
			}
			return nil, &errRes
		}
		log.Println(err)
		errRes = dto.ErrorResponse{
			Status: 500,
			Error:  "Gagal mendapatkan konfigurasi project",
		}
		return nil, &errRes
	}
	return &dto.ProjectConfig{
		ProjectID:          config.ProjectID,
		Threshold:          config.Threshold,
		SessionTime:        config.SessionTime,
		Host:               config.Host.String,
		BaseURL:            config.BaseURL.String,
		MaxUsersInQueue:    config.MaxUsersInQueue,
		QueueStart:         config.QueueStart.Time,
		QueueEnd:           config.QueueEnd.Time,
		QueuePageStyle:     config.QueuePageStyle,
		QueueHTMLPage:      config.QueueHTMLPage.String,
		QueuePageBaseColor: config.QueuePageBaseColor.String,
		QueuePageTitle:     config.QueuePageTitle.String,
		QueuePageLogo:      config.QueuePageLogo.String,
	}, nil
}

func (u *Usecase) UpdateProjectConfig(ctx context.Context, req dto.UpdateProjectConfig) *dto.ErrorResponse {
	var errRes dto.ErrorResponse

	const layout = "2006-01-02T15:04:05"
	queueStart, err := time.Parse(layout, req.QueueStart)
	if err != nil {
		errRes = dto.ErrorResponse{
			Status: 400,
			Error:  "Format waktu queue mulai salah",
		}
		return &errRes
	}

	queueEnd, err := time.Parse(layout, req.QueueEnd)
	if err != nil {
		errRes = dto.ErrorResponse{
			Status: 400,
			Error:  "Format waktu queue berakhir salah",
		}
		return &errRes
	}

	config := entity.Configuration{
		ProjectID:   req.ProjectID,
		Threshold:   req.Threshold,
		SessionTime: req.SessionTime,
		Host: sql.NullString{
			Valid:  true,
			String: req.Host,
		},
		BaseURL: sql.NullString{
			Valid:  true,
			String: req.BaseURL,
		},
		MaxUsersInQueue: req.MaxUsersInQueue,
		QueueStart: sql.NullTime{
			Valid: true,
			Time:  queueStart,
		},
		QueueEnd: sql.NullTime{
			Valid: true,
			Time:  queueEnd,
		},
	}

	err = u.repo.UpdateProjectConfig(ctx, config)
	if err != nil {
		log.Println("Error gagal mengupdate konfigurasi project", err)
		if err == sql.ErrNoRows {
			errRes = dto.ErrorResponse{
				Status: 404,
				Error:  "Project dengan id tersebut tidak ditemukan",
			}
			return &errRes
		}
		errRes = dto.ErrorResponse{
			Status: 500,
			Error:  "Gagal mengupdate konfigurasi project",
		}
		return &errRes
	}

	return nil
}

func (u *Usecase) UpdateProjectStyle(ctx context.Context, req dto.UpdateProjectStyle, imageFile *multipart.FileHeader, htmlFile *multipart.FileHeader) *dto.ErrorResponse {
	logoURL := "https://lh3.googleusercontent.com/proxy/ADW02XxlWJtFJ9MfhL0gRPFhUb9pDx08u6hlXUceO35UBGZncB9B9KdKoeiZW0K6rK1cJfYlRULTZaB-8zOJBFkEuhe8jC_9xivMaDIqA9TpJHQTV_5zmCsNkFzvH0uxICaV-v_F367S8xK5fe2bXINYVkz2CpNToA"
	if req.QueuePageStyle == "base" {
		if imageFile != nil {
			imageContent, err := readFileContent(imageFile)
			if err != nil {
				return handleError(http.StatusBadRequest, "Gagal membaca file image")
			}
			logoURL, err = u.infraRepo.UploadLogoFile(&http.Client{}, dto.File{
				Filename: imageFile.Filename,
				Content:  imageContent,
			})
			if err != nil {
				return handleError(http.StatusInternalServerError, "Gagal upload file image")
			}
		}

		htmlTemplate, err := loadDefaultHTML()
		if err != nil {
			log.Println(err)
			return handleError(http.StatusInternalServerError, "Gagal membuka file template")
		}

		if req.QueuePageBaseColor != "" {
			htmlTemplate = changeValueInHTML(htmlTemplate, "var(--base-color, #f1f1f1)", req.QueuePageBaseColor)
		}

		htmlTemplate = changeValueInHTML(htmlTemplate, "{queue_logo}", logoURL)
		htmlTemplate = changeValueInHTML(htmlTemplate, "{queue_title}", req.QueuePageTitle)

		err = u.infraRepo.UploadHTMLFile(&http.Client{}, dto.File{
			Filename: req.ProjectID,
			Content:  []byte(htmlTemplate),
		})
		if err != nil {
			return handleError(http.StatusInternalServerError, "Gagal upload HTML file")
		}
	} else if req.QueuePageStyle == "custom" {
		if htmlFile == nil {
			return handleError(http.StatusBadRequest, "Mohon sertakan file HTML")
		}
		htmlContent, err := readFileContent(htmlFile)
		if err != nil {
			return handleError(http.StatusBadRequest, "Gagal membaca file HTML")
		}
		err = u.infraRepo.UploadHTMLFile(&http.Client{}, dto.File{
			Filename: req.ProjectID,
			Content:  htmlContent,
		})
		if err != nil {
			return handleError(http.StatusInternalServerError, "Gagal upload HTML file")
		}
	} else {
		return handleError(http.StatusBadRequest, "Tipe style tidak valid")
	}

	config := entity.Configuration{
		ProjectID:      req.ProjectID,
		QueuePageStyle: req.QueuePageStyle,
		QueueHTMLPage: sql.NullString{
			Valid:  true,
			String: fmt.Sprintf("https://storage.googleapis.com/antrein-ta/html_templates/%s.html", req.ProjectID),
		},
		QueuePageBaseColor: sql.NullString{
			Valid:  true,
			String: req.QueuePageBaseColor,
		},
		QueuePageTitle: sql.NullString{
			Valid:  true,
			String: req.QueuePageTitle,
		},
		QueuePageLogo: sql.NullString{
			Valid:  true,
			String: logoURL,
		},
	}

	err := u.repo.UpdateProjectStyle(ctx, config)
	if err != nil {
		log.Println("Error updating project style", err)
		if err == sql.ErrNoRows {
			return handleError(http.StatusNotFound, "Project tidak ditemukan")
		}
		return handleError(http.StatusInternalServerError, "Gagal mengupdate project style")
	}

	return nil
}
