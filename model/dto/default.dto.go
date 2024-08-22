package dto

type File struct {
	Filename string
	Content  []byte
}

type ErrorResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

type DefaultResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginationResponse struct {
	PageSize    int         `json:"pageSize"`
	Page        int         `json:"page"`
	TotalRecord int         `json:"totalRecord"`
	TotalPage   int         `json:"totalPage"`
	Data        interface{} `json:"data"`
}
