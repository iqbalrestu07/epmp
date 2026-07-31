package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Success struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

type Error struct {
	Success bool        `json:"success"`
	Error   ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func OK(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, Success{Success: true, Data: data})
}

func Created(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusCreated, Success{Success: true, Data: data})
}

func NoContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

func BadRequest(c echo.Context, msg string) error {
	return c.JSON(http.StatusBadRequest, Error{
		Success: false,
		Error:   ErrorDetail{Code: "BAD_REQUEST", Message: msg},
	})
}

func NotFound(c echo.Context, msg string) error {
	return c.JSON(http.StatusNotFound, Error{
		Success: false,
		Error:   ErrorDetail{Code: "NOT_FOUND", Message: msg},
	})
}

func Unauthorized(c echo.Context, msg string) error {
	return c.JSON(http.StatusUnauthorized, Error{
		Success: false,
		Error:   ErrorDetail{Code: "UNAUTHORIZED", Message: msg},
	})
}

func Forbidden(c echo.Context, msg string) error {
	return c.JSON(http.StatusForbidden, Error{
		Success: false,
		Error:   ErrorDetail{Code: "FORBIDDEN", Message: msg},
	})
}

func Conflict(c echo.Context, msg string) error {
	return c.JSON(http.StatusConflict, Error{
		Success: false,
		Error:   ErrorDetail{Code: "CONFLICT", Message: msg},
	})
}

func InternalError(c echo.Context, msg string) error {
	return c.JSON(http.StatusInternalServerError, Error{
		Success: false,
		Error:   ErrorDetail{Code: "INTERNAL_ERROR", Message: msg},
	})
}
