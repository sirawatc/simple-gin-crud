package dto

import (
	"net/http"
	"strconv"
)

type Code string

// Standard response codes
const (
	Success             Code = "20000"
	Updated             Code = "20010"
	Deleted             Code = "20020"
	Created             Code = "20100"
	BadRequest          Code = "40000"
	NotFound            Code = "40400"
	Conflict            Code = "40900"
	UnprocessableEntity Code = "42200"
	InternalError       Code = "50000"
)

// Custom response codes
const (
	BindingError      Code = "40010"
	UUIDFormatInvalid Code = "40011"
	ValidationError   Code = "40020"

	BookNotFound   Code = "40401"
	AuthorNotFound Code = "40402"

	BookAlreadyExists   Code = "40901"
	AuthorAlreadyExists Code = "40902"
)

var CodeMessage = map[Code]string{
	Success:             "Success",
	Updated:             "Updated successfully",
	Deleted:             "Deleted successfully",
	Created:             "Created successfully",
	BadRequest:          "Bad Request",
	NotFound:            "Not Found",
	UnprocessableEntity: "Unprocessable Entity",
	InternalError:       "Internal Server Error",

	// Custom response codes
	BindingError:        "JSON parse error",
	UUIDFormatInvalid:   "Invalid UUID format",
	BookNotFound:        "Book not found",
	AuthorNotFound:      "Author not found",
	ValidationError:     "Validation error",
	BookAlreadyExists:   "Book already exists",
	AuthorAlreadyExists: "Author already exists",
}

func (c Code) GetHTTPCode() int {
	if len(c) < 3 {
		return http.StatusInternalServerError
	}
	code, err := strconv.Atoi(string(c)[:3])
	if err != nil {
		return http.StatusInternalServerError
	}
	return code
}
