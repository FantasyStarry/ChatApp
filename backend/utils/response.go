package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 响应码常量定义
const (
	// 成功响应码
	CODE_SUCCESS = 1000

	// 客户端错误 (4xxx)
	CODE_BAD_REQUEST       = 4000 // 请求参数错误
	CODE_UNAUTHORIZED      = 4001 // 未认证
	CODE_FORBIDDEN         = 4003 // 无权限
	CODE_NOT_FOUND         = 4004 // 资源不存在
	CODE_VALIDATION_ERROR  = 4005 // 数据验证失败

	// 服务端错误 (5xxx)
	CODE_INTERNAL_ERROR    = 5000 // 服务器内部错误
	CODE_DATABASE_ERROR    = 5001 // 数据库错误
	CODE_THIRD_PARTY_ERROR = 5002 // 第三方服务错误
)

// 响应码对应的默认消息
var codeMessages = map[int]string{
	CODE_SUCCESS:           "成功",
	CODE_BAD_REQUEST:       "请求参数错误",
	CODE_UNAUTHORIZED:      "未认证或认证失败",
	CODE_FORBIDDEN:         "无权限访问",
	CODE_NOT_FOUND:         "资源不存在",
	CODE_VALIDATION_ERROR:  "数据验证失败",
	CODE_INTERNAL_ERROR:    "服务器内部错误",
	CODE_DATABASE_ERROR:    "数据库操作失败",
	CODE_THIRD_PARTY_ERROR: "第三方服务异常",
}

// ApiResponse 统一API响应结构
type ApiResponse struct {
	Code     int         `json:"code"`
	Messages string      `json:"messages"`
	Data     interface{} `json:"data"`
}

// SuccessResponse 成功响应
func SuccessResponse(c *gin.Context, data interface{}) {
	response := ApiResponse{
		Code:     CODE_SUCCESS,
		Messages: codeMessages[CODE_SUCCESS],
		Data:     data,
	}
	c.JSON(http.StatusOK, response)
}

// SuccessResponseWithMessage 带自定义消息的成功响应
func SuccessResponseWithMessage(c *gin.Context, message string, data interface{}) {
	response := ApiResponse{
		Code:     CODE_SUCCESS,
		Messages: message,
		Data:     data,
	}
	c.JSON(http.StatusOK, response)
}

// ErrorResponse 错误响应
func ErrorResponse(c *gin.Context, httpStatus int, code int, message string) {
	if message == "" {
		if defaultMsg, exists := codeMessages[code]; exists {
			message = defaultMsg
		} else {
			message = "未知错误"
		}
	}

	response := ApiResponse{
		Code:     code,
		Messages: message,
		Data:     nil,
	}
	c.JSON(httpStatus, response)
}

// BadRequestResponse 400错误响应
func BadRequestResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusBadRequest, CODE_BAD_REQUEST, message)
}

// UnauthorizedResponse 401错误响应
func UnauthorizedResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, CODE_UNAUTHORIZED, message)
}

// ForbiddenResponse 403错误响应
func ForbiddenResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusForbidden, CODE_FORBIDDEN, message)
}

// NotFoundResponse 404错误响应
func NotFoundResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, CODE_NOT_FOUND, message)
}

// ValidationErrorResponse 数据验证错误响应
func ValidationErrorResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusBadRequest, CODE_VALIDATION_ERROR, message)
}

// InternalErrorResponse 500错误响应
func InternalErrorResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusInternalServerError, CODE_INTERNAL_ERROR, message)
}

// DatabaseErrorResponse 数据库错误响应
func DatabaseErrorResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusInternalServerError, CODE_DATABASE_ERROR, message)
}