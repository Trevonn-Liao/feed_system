package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const SuccessCode = 0

type Body struct {
	Code    int         `json:"code" example:"0"`
	Message string      `json:"message" example:"ok"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Body{
		Code:    SuccessCode,
		Message: "ok",
		Data:    data,
	})
}

func Error(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Body{
		Code:    code,
		Message: message,
	})
}
