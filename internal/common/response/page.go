package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"lattice-coding/internal/common/errors"
)

type PageResult[T any] struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
	Total   int64  `json:"total"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
}

func OkPage[T any](c *gin.Context, data T, total int64, page, size int) {
	c.JSON(http.StatusOK, PageResult[T]{
		Code:    string(errors.Success),
		Message: "success",
		Data:    data,
		Total:   total,
		Page:    page,
		Size:    size,
	})
}

func FailPage(c *gin.Context, err *errors.BizError) {
	c.JSON(http.StatusOK, PageResult[any]{
		Code:    string(err.Code),
		Message: err.Message,
	})
}

func FailPageWithMessage(c *gin.Context, code string, message string) {
	c.JSON(http.StatusOK, PageResult[any]{
		Code:    code,
		Message: message,
	})
}
