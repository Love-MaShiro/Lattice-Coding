package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"lattice-coding/internal/common/errors"
)

type Result[T any] struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

func Ok[T any](c *gin.Context, data T) {
	c.JSON(http.StatusOK, Result[T]{
		Code:    string(errors.Success),
		Message: "success",
		Data:    data,
	})
}

func OkWithMessage[T any](c *gin.Context, message string, data T) {
	c.JSON(http.StatusOK, Result[T]{
		Code:    string(errors.Success),
		Message: message,
		Data:    data,
	})
}

func Fail(c *gin.Context, err *errors.BizError) {
	c.JSON(http.StatusOK, Result[any]{
		Code:    string(err.Code),
		Message: err.Message,
	})
}

func FailWithMessage(c *gin.Context, code string, message string) {
	c.JSON(http.StatusOK, Result[any]{
		Code:    code,
		Message: message,
	})
}

func FailWithStatus(c *gin.Context, status int, err *errors.BizError) {
	c.JSON(status, Result[any]{
		Code:    string(err.Code),
		Message: err.Message,
	})
}
