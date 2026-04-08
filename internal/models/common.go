package models

// PaginatedResponse wraps any list with pagination metadata.
type PaginatedResponse struct {
	Data     interface{} `json:"data"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// ErrorResponse is the standard error envelope.
type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// PaginationParams holds parsed query params for listing endpoints.
type PaginationParams struct {
	Page     int
	PageSize int
	Search   string
}

func (p *PaginationParams) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *PaginationParams) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
}
