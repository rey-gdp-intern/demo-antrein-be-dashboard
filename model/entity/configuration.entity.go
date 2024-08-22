package entity

import (
	"database/sql"
)

type Configuration struct {
	ID                 string         `db:"id"`
	ProjectID          string         `db:"project_id"`
	Threshold          int            `db:"threshold"`
	SessionTime        int            `db:"session_time"`
	Host               sql.NullString `db:"host"`
	BaseURL            sql.NullString `db:"base_url"`
	MaxUsersInQueue    int            `db:"max_users_in_queue"`
	QueueStart         sql.NullTime   `db:"queue_start"`
	QueueEnd           sql.NullTime   `db:"queue_end"`
	QueuePageStyle     string         `db:"queue_page_style"`
	QueueHTMLPage      sql.NullString `db:"queue_html_page"`
	QueuePageBaseColor sql.NullString `db:"queue_page_base_color"`
	QueuePageTitle     sql.NullString `db:"queue_page_title"`
	QueuePageLogo      sql.NullString `db:"queue_page_logo"`
	IsConfigure        bool           `db:"is_configure"`
	UpdatedAt          sql.NullTime   `db:"updated_at,omitempty"`
}
