package confutil

// Field 配置字段
type Field interface {
	// Name 返回字段名
	Name() string
	// Env 返回字段绑定的环境变量
	Env() []string
	// DefaultValue 返回字段默认值
	DefaultValue() interface{}
}

// defaultField 是 Field 的一个实现
type defaultField struct {
	name         string
	env          []string
	defaultValue interface{}
}

// NewField 创建一个 Field
func NewField(name string, defaultValue interface{}, env ...string) Field {
	return &defaultField{
		name:         name,
		env:          env,
		defaultValue: defaultValue,
	}
}

// Name 返回字段名
func (field *defaultField) Name() string {
	return field.name
}

// Env 返回字段绑定的环境变量
func (field *defaultField) Env() []string {
	var ret []string
	ret = append(ret, field.env...)
	return ret
}

// DefaultValue 返回字段默认值
func (field *defaultField) DefaultValue() interface{} {
	return field.defaultValue
}
