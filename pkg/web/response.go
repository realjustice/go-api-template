package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应（200）
func Success(c *Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 成功响应（自定义消息）
func SuccessWithMessage(c *Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应（自定义状态码和消息）
func Error(c *Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
	})
}

// BadRequest 请求参数错误（400）
func BadRequest(c *Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    400,
		Message: message,
	})
}

// Unauthorized 未授权（401）
func Unauthorized(c *Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    401,
		Message: message,
	})
}

// Forbidden 禁止访问（403）
func Forbidden(c *Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Code:    403,
		Message: message,
	})
}

// NotFound 资源不存在（404）
func NotFound(c *Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Code:    404,
		Message: message,
	})
}

// InternalError 服务器内部错误（500）
func InternalError(c *Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    500,
		Message: message,
	})
}

// Created 创建成功（201）
func Created(c *Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:    201,
		Message: "创建成功",
		Data:    data,
	})
}

// NoContent 无内容（204）
func NoContent(c *Context) {
	c.Status(http.StatusNoContent)
}

// ========== 兼容性方法（用于非 web.HandlerFunc 的场景）==========

// SuccessGin 成功响应（兼容 gin.Context）
func SuccessGin(c *gin.Context, data interface{}) {
	Success(&Context{Context: c}, data)
}

// SuccessWithMessageGin 成功响应（兼容 gin.Context）
func SuccessWithMessageGin(c *gin.Context, message string, data interface{}) {
	SuccessWithMessage(&Context{Context: c}, message, data)
}

// BadRequestGin 请求参数错误（兼容 gin.Context）
func BadRequestGin(c *gin.Context, message string) {
	BadRequest(&Context{Context: c}, message)
}

// UnauthorizedGin 未授权（兼容 gin.Context）
func UnauthorizedGin(c *gin.Context, message string) {
	Unauthorized(&Context{Context: c}, message)
}

// NotFoundGin 资源不存在（兼容 gin.Context）
func NotFoundGin(c *gin.Context, message string) {
	NotFound(&Context{Context: c}, message)
}
