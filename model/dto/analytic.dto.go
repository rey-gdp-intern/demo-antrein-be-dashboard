package dto

import "time"

type Event struct {
	events chan *string
}

type Analytic struct {
	ProjectID         string    `json:"project_id"`
	TimeStamp         time.Time `json:"timestamp"`
	TotalUsersInQueue int       `json:"total_users_in_queue"`
	TotalUsersInRoom  int       `json:"total_users_in_room"`
	TotalUsers        int       `json:"total_users"`
}
