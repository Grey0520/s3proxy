# confutil - 配置加载

**import path: `tencent.com/libs/go/confutil`**

## 核心功能

- 所有配置通过一个结构体定义，配置读取后加载到该结构体的一个实例中，通过该实例使用配置
- 每个配置项的环境变量、默认值都在字段上直接声明，不需要额外的代码
- 配置项环境变量的前缀可以通过层级定义
- 基于 [viper](https://github.com/spf13/viper) ，且兼容 viper 的其它所有功能

## 示例

使用示例代码： [apps/demo/go/srcs/load-config](https://git.woa.com/CodeMatrix/CodeMatrix/tree/master/apps/demo/go/srcs/load-config)

通过 `bazel run //apps/demo/go/srcs/load-config` 即可运行该示例，该示例会打印其加载的配置内容，可以通过通过指定环境变量或配置文件看看其配置加载结果

比如：

```shell
# 修改 http 端口为 8080
SERVE_HTTP_ADDR=8080 bazel run //apps/demo/go/srcs/load-config

# 从配置文件加载
bazel run //apps/demo/go/srcs/load-config -- -config $(pwd)/apps/demo/go/srcs/load-config/samples/config.yaml
```

## 使用流程

**1\. 定义配置**

服务配置应该被定义为一个结构体，可以多个结构体组合：

```go
// Package config
// apps/fc/sample-app/srcs/something/config/config.go
package config

type Config struct {
	Field1 string `env:"FIELD1" default:"value"`
	Field2 int `env:"FIELD2"`
	Extra ExtraConfig `envPrefix:"EXTRA_"`
}

type ExtraConfig struct {
	Field3 float64 `env:"FIELD3"`
}
```

以上定义等价于以下在 viper 中的配置项：

- key: `field1` env: `FIELD1` default: `"value"`
- key: `field2` env: `FIELD2`
- key: `extra.field3` env: `EXTRA_FIELD3`

配置定义详细用法参考 [配置定义说明](#配置定义说明)

**2\. 注册配置、定义配置加载方法**

因为定义的是某服务实例化后的配置，仅对某个服务有效，因此这部分代码建议放到私有包中避免被其它模块错误引用

```go
// Package conf
// apps/fc/sample-app/srcs/something/internal/conf/conf.go
package conf

import (
	"fmt"
	"os"

	"tencent.com/apps/demo/go/srcs/load-config/config"
	"tencent.com/libs/go/confutil"
)

// Config 加载的全局配置
// 从该对象可读取加载后的配置
var Config config.Config

// Manager 配置管理器
var Manager confutil.Manager

func init() {
	// 初始化配置管理器
	Manager = confutil.NewManager()
	// 注册配置
	if err := Manager.RegisterStruct(&config.Config{}); err != nil {
    panic(fmt.Errorf("register config error: %w", err))
	}
}

// LoadConfig 加载配置
func LoadConfig(configPath string) error {
	// 读取配置文件（可选）
	if configPath != "" {
    // 从指定路径读取
    if err := Manager.ReadInFile(configPath, "yaml"); err != nil {
      return fmt.Errorf("read config from file error: %w", err)
    }
	} else {
    // 从默认目录搜索
    // - ./config.yaml
    // - ~/.codematrix/demo/config.yaml
    // 找不到也不会报错
    if err := Manager.SearchInFile(
      "config", "yaml",
      ".",
      os.ExpandEnv("${HOME}/.codematrix/demo"),
    ); err != nil {
      return fmt.Errorf("read config from file error: %w", err)
    }
	}

	// 加载到对象
	if err := Manager.Unmarshal(&Config); err != nil {
    return fmt.Errorf("load config error: %w", err)
	}

	return nil
}
```

**3\. 加载配置、使用配置**

调用上述 `LoadConfig` 方法后配置才可用，因此尽量在程序运行比较早期时调用，比如在 `main` 函数中。

该示例加载配置后直接打印其值：

```go
// Package main
// apps/fc/sample-app/srcs/something/main.go
package main

import (
	"encoding/json"
	"flag"
	"log"

	"tencent.com/apps/fc/sample-app/srcs/something/internal/conf"
)

var (
	configPath = flag.String("config", "", "Config file path")
)

func main() {
	flag.Parse()

	// 加载配置
	if err := conf.LoadConfig(*configPath); err != nil {
		log.Fatalf("load config error: %v", err)
	}

	configJSONRaw, err := json.MarshalIndent(conf.Config, "", "  ")
	if err != nil {
		log.Fatalf("marshal config to json error: %v", err)
	}
	log.Printf("Config:\n%s\n", string(configJSONRaw))
}
```

#### 配置定义说明

参考示例： [apps/demo/go/srcs/load-config](https://git.woa.com/CodeMatrix/CodeMatrix/blob/master/apps/demo/go/srcs/load-config/config/config.go)

##### 环境变量

- 结构体字段通过 `envPrefix` 标签定义该层级以下字段的环境变量的前缀（可嵌套多层）
- 叶子字段通过 `env` 标签定义其绑定的环境变量
  - 同一个字段绑定多个环境变量使用 `,` 分隔（可用于配置调整后兼容旧的环境变量）
  - 以 `^` 开头的环境变量表示不带 `envPrefix` 定义的前缀

比如

```go
package config

type Config struct {
	Middleware MiddlewareConfig `envPrefix:"MIDDLEWARE_"`
}

type MiddlewareConfig struct {
	MySQL MySQLConfig `envPrefix:"MYSQL_"`
}

type MySQLConfig struct {
	Host     string `env:"HOST,^DB_HOST"`
	Port     int    `env:"PORT"`
	User     string `env:"USER,USERNAME"`
	Password string `env:"PASSWORD,PASSWD"`
	Database string `env:"DATABASE,DB"`
}
```

该示例声明了如下环境变量

- `MIDDLEWARE_MYSQL_HOST` 或 `DB_HOST` MySQL 主机名
- `MIDDLEWARE_MYSQL_PORT` MySQL 端口号
- `MIDDLEWARE_MYSQL_USER` 或 `MIDDLEWARE_MYSQL_USERNAME` MySQL 用户名
- `MIDDLEWARE_MYSQL_PASSWORD` 或 `MIDDLEWARE_MYSQL_PASSWD` MySQL 密码
- `MIDDLEWARE_MYSQL_DATABASE` 或 `MIDDLEWARE_MYSQL_DB` MySQL 数据库名

##### 默认值

- 叶子字段通过 `default` 标签定义其默认值
- 默认的默认值是类型的“零值”，所以如果期望的默认值就是“零值”可以不用定义
- 默认值声明是使用 viper 默认的反序列化方法进行反序列化的

比如:

```go
package config

type Config struct {
	Debug bool `default:"true"`
	LogLevel string `default:"info"`
	FloatField float64 `default:"1.2"`
}
```

##### 匿名字段

对于不想多一级嵌套的匿名结构体字段，可以加标签 `mapstructure:",squash"`

比如

```go
package config

import (
	"tencent.com/libs/go/confutil"
)

// Config 配置
type Config struct {
	BaseConfig `mapstructure:",squash"`
	Field string
}

type BaseConfig struct {
	// 是否调试模式
	Debug bool
	// 日志级别
	// debug / info / warning / error
	LogLevel string
}
```

最终解析的 viper 字段为：

- `debug`
- `loglevel`
- `field`

而不会包含 `baseconfig.` 前缀，在从配置文件加载时也不需要多一级：

```yaml
# YAML 示例
debug: false
logLevel: "info"
field: "whatever"
```

##### JSON 字段类型

`[]interface{}` 和 `map[string]interface{}` 类型的字段通过 JSON 进行反序列化

比如：

```go
package config

type Config struct {
	Field []string `env:"FIELD"`
}
```

可以通过环境变量 `FIELD=["value1","value2"]` 指定该字段值

但是在从文件读取时也可以直接使用 YAML 或 JSON 的 list 来表示：

```yaml
# YAML 示例：
field:
  - value1
  - value2
```

也支持表示为 JSON 格式的字符串

```yaml
# YAML 示例：
field: '["value1","value2"]'
```

##### 敏感字符串字段类型

对于包含敏感内容的字符串字段，建议使用 `tencent.com/libs/go/confutil.SecretString` 类型。该类型修改了 string 默认的序列化方法，使其在一般情况下不会直接打印出原文。

```go
package config

import (
	"tencent.com/libs/go/confutil"
)

type MySQLConfig struct {
	Host     string                `json:"host" env:"HOST"`
	Port     int                   `json:"port" env:"PORT" default:"3306"`
	User     string                `json:"user" env:"USER,USERNAME" default:"root"`
	Password confutil.SecretString `json:"password" env:"PASSWORD,PASSWD"`
	Database string                `json:"database" env:"DATABASE,DB"`
}
```

该示例中， `Password` 字段在打印时一般会显示为 `****** (sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855)` 。不会显示原文，提供的哈希也足以校验配置是否正确。

在程序中要使用的时候可以通过 `.Raw()` 方法获取其原始内容

##### Duration 字段类型

`time.Duration` 类型的字段反序列化时使用其字符串格式表示，比如 `1s` `2m` `3h`

##### 忽略字段

默认情况下 confutil 会解析所有字段，遇到不支持的字段类型会报错，如果不希望 confutil 处理某些字段，可以加 tag `confutil:"-"`

比如：

```go
package config

type Config struct {
	Whatever interface{} `confutil:"-"`
}
```

实例中 `Whatever` 字段将会被忽略

#### 直接操作 Viper 对象

如果需要直接操作管理配置的 Viper 对象（比如需要将配置项绑定到 flag ），可以调用实例化的 `tencent.com/libs/go/confutil.Manager` 的 `.Viper()` 方法。比如：

```go
package conf

import (
	"fmt"
	"os"

	"tencent.com/apps/demo/go/srcs/load-config/config"
	"tencent.com/libs/go/confutil"
)

// Config 加载的全局配置
// 从该对象可读取加载后的配置
var Config config.Config

// Manager 配置管理器
var Manager confutil.Manager

func init() {
	// 初始化配置管理器
	Manager = confutil.NewManager()
	// 注册配置
	if err := Manager.RegisterStruct(&config.Config{}); err != nil {
		panic(fmt.Errorf("register config error: %w", err))
	}
	// 打印注册的字段
	fmt.Println(Manager.Viper().AllKeys())  // <--- 通过 .Viper 可以获取管理该配置的 Viper 对象
}
```
