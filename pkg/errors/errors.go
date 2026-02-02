package errors

import (
	"github.com/cockroachdb/errors"
)

// ========== 业务错误定义 ==========

var (
	// 通用错误
	ErrInternal = errors.New("内部服务错误")
	ErrNotFound = errors.New("资源不存在")

	// 认证相关错误
	ErrUnauthorized  = errors.New("未授权")
	ErrInvalidToken  = errors.New("无效的 token")
	ErrTokenNotFound = errors.New("token 不存在或已过期")
	ErrTokenExpired  = errors.New("token 已过期")

	// CheckSum 鉴权错误
	ErrInvalidCheckSum   = errors.New("签名错误")
	ErrInvalidTimestamp  = errors.New("时间戳无效")
	ErrInvalidAppKey     = errors.New("应用 KEY 无效")
	ErrAppNotFound       = errors.New("应用不存在")
	ErrAppRevoked        = errors.New("应用已注销")
	ErrAppExpired        = errors.New("应用已过期")
	ErrMissingAuthParams = errors.New("缺少必要的鉴权参数")

	// 数据库错误
	ErrDatabaseQuery  = errors.New("数据库查询失败")
	ErrDatabaseUpdate = errors.New("数据库更新失败")

	// 缓存错误
	ErrCacheGet = errors.New("缓存获取失败")
	ErrCacheSet = errors.New("缓存设置失败")

	// 参数错误
	ErrInvalidParams = errors.New("参数无效")
	ErrMissingParams = errors.New("缺少必要参数")
)

// ========== 错误包装函数 ==========

// New 创建一个新的错误（带堆栈）
func New(msg string) error {
	return errors.New(msg)
}

// Newf 创建一个格式化的新错误（带堆栈）
func Newf(format string, args ...interface{}) error {
	return errors.Newf(format, args...)
}

// Wrap 包装错误，添加上下文信息（保留堆栈）
// 如果 err 为 nil，返回 nil
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return errors.Wrap(err, msg)
}

// Wrapf 格式化包装错误，添加上下文信息（保留堆栈）
// 如果 err 为 nil，返回 nil
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return errors.Wrapf(err, format, args...)
}

// WithStack 为错误添加堆栈信息
// 如果错误已经有堆栈，不会重复添加
func WithStack(err error) error {
	if err == nil {
		return nil
	}
	return errors.WithStack(err)
}

// WithMessage 为错误添加消息（不添加堆栈）
func WithMessage(err error, msg string) error {
	if err == nil {
		return nil
	}
	return errors.WithMessage(err, msg)
}

// WithMessagef 为错误添加格式化消息（不添加堆栈）
func WithMessagef(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return errors.WithMessagef(err, format, args...)
}

// ========== 错误判断函数 ==========

// Is 判断错误是否匹配目标错误
// 兼容 Go 1.13+ errors.Is
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As 将错误转换为特定类型
// 兼容 Go 1.13+ errors.As
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// ========== 错误增强函数 ==========

// WithHint 为错误添加提示信息
// 提示信息可以帮助用户了解如何解决问题
func WithHint(err error, hint string) error {
	if err == nil {
		return nil
	}
	return errors.WithHint(err, hint)
}

// WithHintf 为错误添加格式化提示信息
func WithHintf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return errors.WithHintf(err, format, args...)
}

// WithDetail 为错误添加详细信息
// 详细信息通常包含技术细节，用于调试
func WithDetail(err error, detail string) error {
	if err == nil {
		return nil
	}
	return errors.WithDetail(err, detail)
}

// WithDetailf 为错误添加格式化详细信息
func WithDetailf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return errors.WithDetailf(err, format, args...)
}

// ========== 错误格式化函数 ==========

// GetMessage 获取错误的完整消息（不包含堆栈）
func GetMessage(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// GetAllHints 获取错误链中的所有提示信息
func GetAllHints(err error) []string {
	if err == nil {
		return nil
	}
	return errors.GetAllHints(err)
}

// GetAllDetails 获取错误链中的所有详细信息
func GetAllDetails(err error) []string {
	if err == nil {
		return nil
	}
	return errors.GetAllDetails(err)
}

// ========== 辅助函数 ==========

// Cause 获取错误的根本原因（最底层的错误）
func Cause(err error) error {
	return errors.UnwrapAll(err)
}

// Unwrap 解包一层错误
func Unwrap(err error) error {
	return errors.Unwrap(err)
}
