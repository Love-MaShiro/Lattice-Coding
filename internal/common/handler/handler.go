package handler

import (
	"github.com/gin-gonic/gin"

	"lattice-coding/internal/common/errors"
	"lattice-coding/internal/common/response"
)

func BindQuery[T any](c *gin.Context) (*T, error) {
	var req T
	if err := c.ShouldBindQuery(&req); err != nil {
		return nil, errors.InvalidArg("参数校验失败: " + err.Error())
	}
	return &req, nil
}

func BindJSON[T any](c *gin.Context) (*T, error) {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, errors.InvalidArg("参数校验失败: " + err.Error())
	}
	return &req, nil
}

func BindUri[T any](c *gin.Context) (*T, error) {
	var req T
	if err := c.ShouldBindUri(&req); err != nil {
		return nil, errors.InvalidArg("路径参数校验失败: " + err.Error())
	}
	return &req, nil
}

func Bind[T any](c *gin.Context) (*T, error) {
	var req T
	if err := c.ShouldBind(&req); err != nil {
		return nil, errors.InvalidArg("参数校验失败: " + err.Error())
	}
	return &req, nil
}

func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	if bizErr, ok := err.(*errors.BizError); ok {
		response.Fail(c, bizErr)
		return
	}

	response.Fail(c, errors.InternalWithErr(err, "系统内部错误"))
}

func HandleErrorWithStatus(c *gin.Context, status int, err error) {
	if err == nil {
		return
	}

	if bizErr, ok := err.(*errors.BizError); ok {
		response.FailWithStatus(c, status, bizErr)
		return
	}

	response.FailWithStatus(c, status, errors.InternalWithErr(err, "系统内部错误"))
}

func MustBindQuery[T any](c *gin.Context) *T {
	req, err := BindQuery[T](c)
	if err != nil {
		response.Fail(c, err.(*errors.BizError))
		c.Abort()
		return nil
	}
	return req
}

func MustBindJSON[T any](c *gin.Context) *T {
	req, err := BindJSON[T](c)
	if err != nil {
		response.Fail(c, err.(*errors.BizError))
		c.Abort()
		return nil
	}
	return req
}

func MustBindUri[T any](c *gin.Context) *T {
	req, err := BindUri[T](c)
	if err != nil {
		response.Fail(c, err.(*errors.BizError))
		c.Abort()
		return nil
	}
	return req
}
