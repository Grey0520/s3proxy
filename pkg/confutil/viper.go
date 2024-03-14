package confutil

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// ConfigItem 配置项
type ConfigItem struct {
	Name       string
	BindingEnv string
	DefaultVal interface{}
}

// ConfigManager 配置管理器
type ConfigManager struct {
	v *viper.Viper

	// 注册的配置项
	Items []ConfigItem

	// 配置文件名（不含文件扩展名）
	ConfigFileName string

	// 配置文件类型
	ConfigFileType string

	// 配置文件搜索目录
	ConfigFileSearchPaths []string
}

// Init 初始化配置管理器
func (mgr *ConfigManager) Init() error {
	if mgr.v != nil {
		return ErrInitMultiTimes
	}

	mgr.v = viper.New()

	// 设置默认值
	if err := mgr.setDefaults(); err != nil {
		return fmt.Errorf("set default value for config items error: %w", err)
	}

	// 读取配置文件
	if err := mgr.readInFiles(); err != nil {
		return fmt.Errorf("read config file error: %w", err)
	}

	// 绑定环境变量
	if err := mgr.bindEnvs(); err != nil {
		return fmt.Errorf("bind env for config items error: %w", err)
	}

	return nil
}

// Viper 获取管理的 Viper 对象
func (mgr *ConfigManager) Viper() *viper.Viper {
	return mgr.v
}

// Unmarshal 将配置反序列化到结构体
func (mgr *ConfigManager) Unmarshal(rawVal interface{}, opts ...viper.DecoderConfigOption) error {
	if opts == nil {
		opts = []viper.DecoderConfigOption{DefaultViperDecoderConfigOption}
	}
	return mgr.v.Unmarshal(rawVal, opts...)
}

// DefaultViperDecoderConfigOption 默认 Viper 解析配置
func DefaultViperDecoderConfigOption(config *mapstructure.DecoderConfig) {
	config.DecodeHook = mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
		StringToMapHookFunc(),
		StringToSliceHookFunc(),
	)
}

// setDefaults 设置默认值
func (mgr *ConfigManager) setDefaults() error {
	// 设置变量类型与默认值相同
	mgr.v.SetTypeByDefaultValue(true)

	// 设置默认值
	for _, item := range mgr.Items {
		mgr.v.SetDefault(item.Name, item.DefaultVal)
	}

	return nil
}

// readInFiles 从配置文件读取
func (mgr *ConfigManager) readInFiles() error {
	// 无配置文件
	if mgr.ConfigFileName == "" || len(mgr.ConfigFileSearchPaths) == 0 {
		return nil
	}

	// 设置配置文件搜索设置
	mgr.v.SetConfigName(mgr.ConfigFileName)
	mgr.v.SetConfigType(mgr.ConfigFileType)
	for _, p := range mgr.ConfigFileSearchPaths {
		mgr.v.AddConfigPath(p)
	}

	// 读取配置文件
	if err := mgr.v.ReadInConfig(); err != nil {
		// 忽略文件找不到错误
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil
		}
		return err
	}
	return nil
}

// bindEnv 绑定环境变量到配置项
func (mgr *ConfigManager) bindEnvs() error {
	for _, item := range mgr.Items {
		// 跳过无需绑定的
		if item.BindingEnv == "" {
			continue
		}

		// 绑定环境变量
		if err := mgr.v.BindEnv(item.Name, item.BindingEnv); err != nil {
			return fmt.Errorf(
				"bind env to config item error: %w, name: %s, env: %s",
				err,
				item.Name,
				item.BindingEnv,
			)
		}
	}
	return nil
}
