package constants

// HTTP Header 常量
const (
	// 认证相关 Header
	HeaderRequestID = "X-Request-ID" // 请求 ID

	// CheckSum 鉴权 Header
	HeaderAppKey    = "app_key"   // 应用 KEY
	HeaderTimestamp = "timestamp" // 时间戳
	HeaderNonce     = "nonce"     // 随机字符串
	HeaderCheckSum  = "checksum"  // 签名
)
