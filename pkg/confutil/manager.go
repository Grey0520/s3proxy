package confutil

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/viper"
)

// Manager 配置管理器
type Manager interface {
	// Viper 获取管理配置的 Viper 实例
	Viper() *viper.Viper
	// Register 注册字段
	Register(fields ...Field) error
	// RegisterStruct 注册结构体
	RegisterStruct(ins interface{}) error
	// SearchInFile 搜索并读取配置文件
	SearchInFile(name, format string, searchPaths ...string) error
	// ReadInFile 从指定文件读取
	ReadInFile(configPath, format string) error
	// Unmarshal 将配置反序列化到配置对象中
	Unmarshal(obj interface{}, opts ...viper.DecoderConfigOption) error
}

// defaultManager 是 Manager 的一个实现
type defaultManager struct {
	v *viper.Viper

	initOnce sync.Once
}

// NewManager 创建一个 Manager
func NewManager() Manager {
	mgr := &defaultManager{}
	mgr.initViper()
	return mgr
}

// initViper 初始化 Viper
func (mgr *defaultManager) initViper() {
	mgr.initOnce.Do(func() {
		mgr.v = viper.New()
		// 设置变量类型与默认值相同
		mgr.v.SetTypeByDefaultValue(true)
	})
}

// Viper 获取管理配置的 Viper 实例
func (mgr *defaultManager) Viper() *viper.Viper {
	return mgr.v
}

// Register 注册字段
func (mgr *defaultManager) Register(fields ...Field) error {
	for _, f := range fields {
		// 设置默认值
		mgr.v.SetDefault(f.Name(), f.DefaultValue())
		// 绑定环境变量
		if err := mgr.v.BindEnv(append([]string{f.Name()}, f.Env()...)...); err != nil {
			return fmt.Errorf("bind env %v to field \"%s\" error: %w", f.Env(), f.Name(), err)
		}
	}
	return nil
}

// RegisterStruct 注册结构体
func (mgr *defaultManager) RegisterStruct(ins interface{}) error {
	// 解析结构体
	fields, err := ExpendConfigStruct(ins)
	if err != nil {
		return fmt.Errorf("parse config struct error: %w", err)
	}

	// 注册字段
	return mgr.Register(fields...)
}

// SearchInFile 从配置文件读取
func (mgr *defaultManager) SearchInFile(name, format string, searchPaths ...string) error {
	// 设置配置文件读取路径
	mgr.v.SetConfigName(name)
	mgr.v.SetConfigType(format)
	for _, p := range searchPaths {
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

// ReadInFile 从指定配置文件读取
func (mgr *defaultManager) ReadInFile(configPath, format string) error {
	mgr.v.SetConfigType(format)

	file, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("open config file error: %w", err)
	}

	return mgr.v.ReadConfig(file)
}

// Unmarshal 将配置反序列化到配置对象中
func (mgr *defaultManager) Unmarshal(obj interface{}, opts ...viper.DecoderConfigOption) error {
	if opts == nil {
		opts = []viper.DecoderConfigOption{DefaultViperDecoderConfigOption}
	}
	return mgr.v.Unmarshal(obj, opts...)
}
