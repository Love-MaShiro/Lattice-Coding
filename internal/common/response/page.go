package response

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"lattice-coding/internal/common/errors"
)

type PageRequest struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}

type PageResult[T any] struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
	Total   int64  `json:"total"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
}

func ParsePageRequest(c *gin.Context) PageRequest {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}

	return PageRequest{
		Page: page,
		Size: size,
	}
}

func (pr PageRequest) Offset() int {
	return (pr.Page - 1) * pr.Size
}

func (pr PageRequest) Limit() int {
	return pr.Size
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
