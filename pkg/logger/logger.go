package logger

import (
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// Logger 全局日志实例
	Logger *zap.Logger
	// Sugar 全局 SugaredLogger 实例（更方便的 API）
	Sugar *zap.SugaredLogger
)

// Field 日志字段类型（隔离 zap 依赖）
type Field = zapcore.Field

// Config 日志配置
type Config struct {
	Level      string // 日志级别：debug, info, warn, error
	Filename   string // 日志文件路径
	MaxSize    int    // 单个日志文件最大大小（MB）
	MaxBackups int    // 保留的旧日志文件数量
	MaxAge     int    // 保留旧日志文件的最大天数
	Compress   bool   // 是否压缩旧日志文件
	Console    bool   // 是否同时输出到控制台
}

// NewLogger 创建日志实例
func NewLogger(cfg *Config) (*zap.Logger, error) {
	// 设置日志级别
	level := zapcore.InfoLevel
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	}

	// 创建日志目录
	if cfg.Filename != "" {
		dir := filepath.Dir(cfg.Filename)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}

	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 创建 Core
	var cores []zapcore.Core

	// 文件输出
	if cfg.Filename != "" {
		// 配置 lumberjack 实现日志滚动
		fileWriter := &lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}

		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(fileWriter),
			level,
		)
		cores = append(cores, fileCore)
	}

	// 控制台输出
	if cfg.Console {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		consoleCore := zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			level,
		)
		cores = append(cores, consoleCore)
	}

	// 创建 logger
	core := zapcore.NewTee(cores...)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	// 设置全局实例
	Logger = logger
	Sugar = logger.Sugar()

	return logger, nil
}

// Close 关闭日志
func Close() error {
	if Logger != nil {
		return Logger.Sync()
	}
	return nil
}

// 便捷方法

// Debug 调试日志
func Debug(msg string, fields ...Field) {
	Logger.Debug(msg, fields...)
}

// Info 信息日志
func Info(msg string, fields ...Field) {
	Logger.Info(msg, fields...)
}

// Warn 警告日志
func Warn(msg string, fields ...Field) {
	Logger.Warn(msg, fields...)
}

// Error 错误日志
func Error(msg string, fields ...Field) {
	Logger.Error(msg, fields...)
}

// Fatal 致命错误日志
func Fatal(msg string, fields ...Field) {
	Logger.Fatal(msg, fields...)
}

// Debugf 格式化调试日志
func Debugf(format string, args ...interface{}) {
	Sugar.Debugf(format, args...)
}

// Infof 格式化信息日志
func Infof(format string, args ...interface{}) {
	Sugar.Infof(format, args...)
}

// Warnf 格式化警告日志
func Warnf(format string, args ...interface{}) {
	Sugar.Warnf(format, args...)
}

// Errorf 格式化错误日志
func Errorf(format string, args ...interface{}) {
	Sugar.Errorf(format, args...)
}

// Fatalf 格式化致命错误日志
func Fatalf(format string, args ...interface{}) {
	Sugar.Fatalf(format, args...)
}

// With 添加字段
func With(fields ...Field) *zap.Logger {
	return Logger.With(fields...)
}

// ============================================
// Field 构造函数（隔离 zap 依赖）
// ============================================

// String 字符串字段
func String(key, val string) Field {
	return zap.String(key, val)
}

// Int 整数字段
func Int(key string, val int) Field {
	return zap.Int(key, val)
}

// Int64 64位整数字段
func Int64(key string, val int64) Field {
	return zap.Int64(key, val)
}

// Uint 无符号整数字段
func Uint(key string, val uint) Field {
	return zap.Uint(key, val)
}

// Uint64 64位无符号整数字段
func Uint64(key string, val uint64) Field {
	return zap.Uint64(key, val)
}

// Float64 浮点数字段
func Float64(key string, val float64) Field {
	return zap.Float64(key, val)
}

// Bool 布尔字段
func Bool(key string, val bool) Field {
	return zap.Bool(key, val)
}

// Time 时间字段
func Time(key string, val time.Time) Field {
	return zap.Time(key, val)
}

// Duration 时长字段
func Duration(key string, val time.Duration) Field {
	return zap.Duration(key, val)
}

// Err 错误字段（key 为 "error"）
func Err(err error) Field {
	return zap.Error(err)
}

// NamedErr 自定义 key 的错误字段
func NamedErr(key string, err error) Field {
	return zap.NamedError(key, err)
}

// Any 任意类型字段（使用反射，性能较低）
func Any(key string, val interface{}) Field {
	return zap.Any(key, val)
}

// Strings 字符串数组字段
func Strings(key string, val []string) Field {
	return zap.Strings(key, val)
}

// Ints 整数数组字段
func Ints(key string, val []int) Field {
	return zap.Ints(key, val)
}

// Skip 跳过字段（用于条件性添加字段）
func Skip() Field {
	return zap.Skip()
}
