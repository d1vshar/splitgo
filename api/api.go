package api

import (
	"github.com/d1vshar/splitgo/core"
	"github.com/labstack/echo/v4"
)

type Api struct {
	App *core.BaseApp
}

type Response struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func SuccessResponse(c echo.Context, code int, data interface{}) error {
	return c.JSON(code, Response{
		Code:    code,
		Success: true,
		Data:    data,
	})
}

func ErrorResponse(c echo.Context, code int, err error) error {
	return c.JSON(code, Response{
		Code:    code,
		Success: false,
		Error: func() string {
			if err == nil {
				return ""
			}

			return err.Error()
		}()})
}
