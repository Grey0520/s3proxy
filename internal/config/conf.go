package config

import (
	"fmt"

	"github.com/Grey0520/s3proxy/pkg/confutil"
)

var Cfg Config
var Manager confutil.Manager

func init() {
	Manager = confutil.NewManager()

	if err := Manager.RegisterStruct(&Config{}); err != nil {
		panic(fmt.Errorf("register config error: %w", err))
	}
}

func LoadConfig(configPath string) error {
	if configPath != "" {
		if err := Manager.ReadInFile(configPath, "yaml"); err != nil {
			return fmt.Errorf("read config from file error: %w", err)
		}
	} else {
		// 没有输入地址，则在程序目录搜索 config.yaml
		if err := Manager.SearchInFile(
			"config", "yaml", "./"); err != nil {
			return fmt.Errorf("read config from file error: %w", err)
		}
	}

	if err := Manager.Unmarshal(&Cfg); err != nil {
		return fmt.Errorf("load config error: %w", err)
	}

	return nil
}
