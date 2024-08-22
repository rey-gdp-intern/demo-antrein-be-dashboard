package infra

import (
	"antrein/bc-dashboard/model/config"
	"antrein/bc-dashboard/model/dto"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type Repository struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Repository {
	return &Repository{
		cfg: cfg,
	}
}

type InfraBody struct {
	ProjectID     string `json:"project_id"`
	ProjectDomain string `json:"project_domain"`
	URLPath       string `json:"url_path"`
}

func (r *Repository) GetInfraProjects(client *http.Client) ([]string, error) {
	req, err := http.NewRequest("GET", r.cfg.Infra.ManagerURL+"/kube/project", nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data []string `json:"data"`
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (r *Repository) CreateInfraProject(client *http.Client, infraBody InfraBody) error {
	jsonData, err := json.Marshal(infraBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", r.cfg.Infra.ManagerURL+"/kube/project", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create project, status code: %d", resp.StatusCode)
	}
	return nil
}

func (r *Repository) UploadLogoFile(client *http.Client, file dto.File) (string, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormFile("file", file.Filename)
	if err != nil {
		return "", err
	}
	if _, err = fw.Write(file.Content); err != nil {
		return "", err
	}

	w.Close()

	req, err := http.NewRequest("POST", r.cfg.Infra.ManagerURL+"/storage/assets", &b)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    struct {
			URL string `json:"url"`
		} `json:"data"`
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to upload file, status code: %d", resp.StatusCode)
	}

	return result.Data.URL, nil
}

func (r *Repository) UploadHTMLFile(client *http.Client, file dto.File) error {
	encodedContent := base64.StdEncoding.EncodeToString(file.Content)

	payload := map[string]string{
		"file_name":   file.Filename,
		"html_base64": encodedContent,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", r.cfg.Infra.ManagerURL+"/storage/html", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to upload file, status code: %d", resp.StatusCode)
	}

	return nil
}

func (r *Repository) CheckHealthProject(client *http.Client, projectId string) (bool, error) {
	req, err := http.NewRequest("GET", r.cfg.Infra.ManagerURL+"/kube/health/"+projectId, nil)
	if err != nil {
		return false, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var result struct {
		Status string `json:"status"`
		Data   struct {
			Healthiness bool `json:"healthiness"`
		} `json:"data"`
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return false, err
	}

	if result.Status == "failed" {
		return false, errors.New("Project not found")
	}
	return result.Data.Healthiness, nil
}

func (r *Repository) DeleteInfraProject(client *http.Client, projectId string) error {
	req, err := http.NewRequest("DELETE", r.cfg.Infra.ManagerURL+"/kube/project/"+projectId, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete project, status code: %d", resp.StatusCode)
	}
	return nil
}

func (r *Repository) ClearInfraProject(client *http.Client) error {
	req, err := http.NewRequest("DELETE", r.cfg.Infra.ManagerURL+"/kube/restart/project", nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to clear project, status code: %d", resp.StatusCode)
	}
	return nil
}
