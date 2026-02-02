package constants

// API 响应消息常量
const (
	// 通用消息
	MsgSuccess = "success"
	MsgFailed  = "failed"

	// 错误消息
	MsgInterfaceNotFound  = "接口不存在"
	MsgMethodNotAllowed   = "请求方法不允许"
	MsgBadRequest         = "请求参数错误"
	MsgUnauthorized       = "未授权"
	MsgForbidden          = "禁止访问"
	MsgNotFound           = "资源不存在"
	MsgInternalError      = "服务器内部错误"
	MsgServiceUnavailable = "服务暂时不可用"
)
