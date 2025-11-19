package http

import (
	"github.com/labstack/echo/v4"
)

type SuccessResponse struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

func Success(c echo.Context, status int, data interface{}) error {
	return c.JSON(status, SuccessResponse{
		Status: status,
		Data:   data,
	})
}
