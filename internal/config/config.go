package config

type S3ProxyConfig struct {
	Endpoint       string `env:"ENDPOINT"`
	SecureEndpoint string `env:"SECURE_ENDPOINT"`
}

// tip: Provider 对应的环境变量为: CLOUD_PROVIDER , 从 Config 里继承的 CLOUD
type CloudsConfig struct {
	Provider   string           `env:"PROVIDER"`
	Identity   string           `env:"IDENTITY"`
	Endpoint   string           `env:"ENDPOINT"`
	Appid      string           `env:"APPID"`
	Key        string           `env:"KEY"`
	Credential string           `env:"CREDENTIAL"`
	Filesystem FilesystemConfig `envPrefix:"FILESYSTEM_"`
}

type FilesystemConfig struct {
	Basedir string `env:"BASEDIR"`
}

// Config 是配置文件的最顶级
type Config struct {
	S3Proxy S3ProxyConfig `envPrefix:"S3PROXY_"`
	Cloud   CloudsConfig  `envPrefix:"CLOUD_"`
}
