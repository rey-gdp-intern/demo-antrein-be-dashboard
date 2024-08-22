package dto

type Project struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`
	Name     string `json:"name"`
}

type CreateProjectRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProjectDetailResponse struct {
	ID            string        `json:"id"`
	TenantID      string        `json:"tenant_id"`
	Name          string        `json:"name"`
	Configuration ProjectConfig `json:"configuration"`
}

type ListProjectResponse struct {
	TenantID string    `json:"tenant_id"`
	Projects []Project `json:"projects"`
}

type CheckHealthProjectResponse struct {
	ID          string `json:"id"`
	Healthiness bool   `json:"healthiness"`
}

type CreateProjectResponse struct {
	Project
}
