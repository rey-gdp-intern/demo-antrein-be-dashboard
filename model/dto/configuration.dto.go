package dto

import "time"

type ProjectConfig struct {
	ProjectID          string    `json:"project_id"`
	Threshold          int       `json:"threshold"`
	SessionTime        int       `json:"session_time"`
	Host               string    `json:"host"`
	BaseURL            string    `json:"base_url"`
	MaxUsersInQueue    int       `json:"max_users_in_queue"`
	QueueStart         time.Time `json:"queue_start"`
	QueueEnd           time.Time `json:"queue_end"`
	QueuePageStyle     string    `json:"queue_page_style"`
	QueueHTMLPage      string    `json:"queue_html_page"`
	QueuePageBaseColor string    `json:"queue_page_base_color"`
	QueuePageTitle     string    `json:"queue_page_title"`
	QueuePageLogo      string    `json:"queue_page_logo"`
	IsConfigure        bool      `json:"is_configure"`
}

type UpdateProjectConfig struct {
	ProjectID       string `json:"project_id"`
	Threshold       int    `json:"threshold"`
	SessionTime     int    `json:"session_time"`
	Host            string `json:"host"`
	BaseURL         string `json:"base_url"`
	MaxUsersInQueue int    `json:"max_users_in_queue"`
	QueueStart      string `json:"queue_start"`
	QueueEnd        string `json:"queue_end"`
}

type UpdateProjectStyle struct {
	ProjectID          string `json:"project_id"`
	QueuePageStyle     string `json:"queue_page_style"`
	QueuePageBaseColor string `json:"queue_page_base_color,omitempty"`
	QueuePageTitle     string `json:"queue_page_title,omitempty"`
}
