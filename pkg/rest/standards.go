package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"strconv"
)

const (
	OK   ResponseStatus = "ok"
	Fail ResponseStatus = "fail"

	Page    = "page"
	PerPage = "per_page"
)

var (
	InternalServerError = GetFailedResponseFromMessage("Internal Server Error")
	NotFound            = GetFailedResponseFromMessage("Not Found!")
)

type ResponseStatus string

// Filter is a map used to filter entities in queries.
type Filter map[string]string

// Query represents a query for listing entities with filters and pagination metadata.
type Query struct {
	Filter         Filter `json:"filter"` // Dynamic filtering by fields like "status" or "assignee_id"
	PaginationMeta        // Embeds pagination info (page, per_page, total)
}

type PaginationMeta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type StandardResponse struct {
	Data    interface{}    `json:"data"`
	Meta    PaginationMeta `json:"meta,omitempty"`
	Message string         `json:"message"`
	Errors  []string       `json:"errors,omitempty"`
	Status  ResponseStatus `json:"status"`
}

func GetSuccessResponse(data interface{}) StandardResponse {
	return StandardResponse{
		Status:  OK,
		Message: "success",
		Data:    data,
	}
}

func GetFailedResponseFromMessage(message string) StandardResponse {
	return StandardResponse{
		Status:  Fail,
		Message: message,
		Data:    nil,
	}
}

func GetFailedResponseFromMessageAndErrors(message string, errors []error) StandardResponse {
	var errMsg []string
	for _, e := range errors {
		errMsg = append(errMsg, e.Error())
	}
	return StandardResponse{
		Status:  Fail,
		Message: message,
		Data:    nil,
		Errors:  errMsg,
	}
}

func GetFailedValidationResponse(err error) StandardResponse {
	var errorsList []error
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, e := range ve {
			errorsList = append(errorsList, e)
		}
	} else {
		errorsList = append(errorsList, err)
	}

	return GetFailedResponseFromMessageAndErrors("validation error", errorsList)
}

func GetSuccessResponseWithMeta(data interface{}, meta PaginationMeta) StandardResponse {
	return StandardResponse{
		Status:  OK,
		Message: "success",
		Data:    data,
		Meta:    meta,
	}
}

func ParseQuery(c *gin.Context) (query Query) {
	query = Query{
		Filter: make(Filter),
	}

	// iterate all query parameters
	for key, values := range c.Request.URL.Query() {
		if len(values) == 0 {
			continue
		}
		value := values[0]

		switch key {
		case Page:
			query.Page, _ = strconv.Atoi(value)
		case PerPage:
			query.PerPage, _ = strconv.Atoi(value)
		default:
			query.Filter[key] = value
		}
	}

	if query.Page == 0 {
		query.Page = 1
	}

	if query.PerPage == 0 {
		query.PerPage = 20
	}

	return
}
