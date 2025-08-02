package dto

import "strconv"

type PaginationRequest struct {
	Page     int `json:"page" form:"page" binding:"min=1"`
	PageSize int `json:"pageSize" form:"pageSize" binding:"min=1,max=100"`
}

func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *PaginationRequest) GetLimit() int {
	return p.PageSize
}

type PaginationResponse struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalPages int   `json:"totalPages"`
	TotalItems int64 `json:"totalItems"`
}

type PaginationDataResponse[T any] struct {
	Items      []T                `json:"items"`
	Pagination PaginationResponse `json:"pagination"`
}

func NewPaginationRequest(page, pageSize string) (*PaginationRequest, []string) {
	errors := []string{}

	pagination := &PaginationRequest{
		Page:     1,
		PageSize: 10,
	}

	if page != "" {
		if page, err := strconv.Atoi(page); err == nil && page > 0 {
			pagination.Page = page
		} else {
			errors = append(errors, "Page must be greater than 0")
		}
	}

	if pageSize != "" {
		if pageSize, err := strconv.Atoi(pageSize); err == nil && pageSize > 0 {
			pagination.PageSize = pageSize
		} else {
			errors = append(errors, "Page size must be greater than 0")
		}
	}

	return pagination, errors
}

func NewPaginationDataResponse[T any](items []T, req *PaginationRequest, totalItems int64) *PaginationDataResponse[T] {
	return &PaginationDataResponse[T]{
		Items:      items,
		Pagination: *NewPaginationResponse(req, totalItems),
	}
}

func NewPaginationResponse(req *PaginationRequest, totalItems int64) *PaginationResponse {
	totalPages := int((totalItems + int64(req.PageSize) - 1) / int64(req.PageSize))
	if totalPages < 1 {
		totalPages = 1
	}

	return &PaginationResponse{
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
		TotalItems: totalItems,
	}
}
