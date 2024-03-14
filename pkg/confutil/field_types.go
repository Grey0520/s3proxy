package confutil

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

// SecretString 加密字符串
type SecretString string

var (
	_ fmt.Stringer   = SecretString("")
	_ yaml.Marshaler = SecretString("")
	_ json.Marshaler = SecretString("")
)

// Raw 原文
func (ss SecretString) Raw() string {
	return string(ss)
}

// String 打印的字符串
func (ss SecretString) String() string {
	return fmt.Sprintf("****** (sha256:%x)", sha256.Sum256([]byte(ss.Raw())))
}

// MarshalJSON 序列化为 JSON
func (ss SecretString) MarshalJSON() ([]byte, error) {
	return []byte("\"" + ss.String() + "\""), nil
}

// MarshalYAML 序列化为 YAML
func (ss SecretString) MarshalYAML() (interface{}, error) {
	return ss.String(), nil
}

// SentryConfig Sentry 配置
type SentryConfig struct {
	// Sentry DSN
	DSN SecretString `json:"dsn" env:"DSN"`
	// 是否调试模式
	Debug bool `json:"debug,omitempty" env:"DEBUG"`
	// 环境
	Environment string `json:"environment,omitempty" env:"ENVIRONMENT"`
	// 采样率 0 - 1
	TracesSampleRate float64 `json:"tracesSampleRate,omitempty" env:"TRACES_SAMPLE_RATE" default:"1.0"`
}
