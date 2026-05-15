package middleware

import (
	"errors"
	"net/http"

	"feed_system/internal/errs"
	"feed_system/internal/response"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				response.Error(c, http.StatusInternalServerError, errs.CodeInternal, "internal server error")
				c.Abort()
			}
		}()

		c.Next()

		if len(c.Errors) == 0 || c.Writer.Written() {
			return
		}

		err := c.Errors.Last().Err

		var appErr *errs.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr.HTTPStatus, appErr.Code, appErr.Message)
			return
		}

		response.Error(c, http.StatusInternalServerError, errs.CodeInternal, "internal server error")
	}
}

func Wrap(handler func(*gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := handler(c); err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}
	}
}
