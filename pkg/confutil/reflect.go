package confutil

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
)

const (
	confutilTag = "confutil"
)

// getStructFieldEnv 从结构体字段获取绑定的环境变量
func getStructFieldEnv(field reflect.StructField, prefix []string) []string {
	envTag := field.Tag.Get("env")
	if envTag == "" {
		return nil
	}
	env := strings.Split(envTag, ",")
	prefixStr := strings.Join(prefix, "")
	for i := range env {
		if strings.HasPrefix(env[i], "^") {
			env[i] = env[i][1:]
		} else {
			env[i] = prefixStr + env[i]
		}
	}
	return env
}

// getStructFieldEnvPrefix 获取结构体字段绑定的环境变量前缀
func getStructFieldEnvPrefix(field reflect.StructField, prefix []string) []string {
	ret := make([]string, len(prefix))
	copy(ret, prefix)
	envPrefixTag := field.Tag.Get("envPrefix")
	if envPrefixTag != "" {
		ret = append(ret, envPrefixTag)
	}
	return ret
}

// getStructFieldName 获取结构体字段名
func getStructFieldName(field reflect.StructField, prefix []string) []string {
	ret := make([]string, len(prefix))
	copy(ret, prefix)
	if field.Name != "" && !strings.Contains(field.Tag.Get("mapstructure"), ",squash") {
		ret = append(ret, field.Name)
	}
	return ret
}

// getStructFieldDefault 从结构体字段获取其默认值
func getStructFieldDefault(field reflect.StructField) (interface{}, error) {
	defaultTag := field.Tag.Get("default")
	defaultValue := reflect.New(field.Type)
	if defaultTag == "" {
		return defaultValue.Elem().Interface(), nil
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           defaultValue.Interface(),
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc( // TODO(P2): 应当可配置
			mapstructure.StringToTimeDurationHookFunc(),
			StringToMapHookFunc(),
			StringToSliceHookFunc(),
		),
	})
	if err != nil {
		return nil, fmt.Errorf("new mapstructure decoder error: %w", err)
	}

	if err := decoder.Decode(defaultTag); err != nil {
		return nil, err
	}

	return defaultValue.Elem().Interface(), nil
}

// isNestedStructFieldType 判断是否嵌套字段
func isNestedStructFieldType(t reflect.Type) bool {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Kind() == reflect.Struct
}

// isNormalStructFieldType 判断是否普通字段类型
func isNormalStructFieldType(t reflect.Type) bool {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Array, reflect.Map, reflect.Slice,
		reflect.String:
		return true
	}
	return false
}

type structField struct {
	t          reflect.StructField
	namePrefix []string
	envPrefix  []string
}

// ExpendConfigStruct 展开配置结构体
func ExpendConfigStruct(ins interface{}) ([]Field, error) {
	// 获取类型
	rootFieldType := reflect.TypeOf(ins)
	if rootFieldType.Kind() == reflect.Ptr {
		rootFieldType = rootFieldType.Elem()
	}
	if rootFieldType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%#v is not a struct instance", ins)
	}

	structFields := []structField{{
		t: reflect.StructField{
			Name:      "",
			PkgPath:   "",
			Type:      rootFieldType,
			Tag:       "",
			Offset:    0,
			Index:     []int{0},
			Anonymous: true,
		},
	}}

	// 递归解析结构体
	var fields []Field
	for len(structFields) > 0 {
		// 取栈顶一个字段
		field := structFields[0]
		structFields = structFields[1:]

		fieldName := getStructFieldName(field.t, field.namePrefix)
		fieldNameStr := strings.Join(fieldName, ".")
		envPrefix := getStructFieldEnvPrefix(field.t, field.envPrefix)

		switch {
		case field.t.Tag.Get(confutilTag) == "-": // 带 tag `confutil:"-"` 的字段
			continue
		case isNestedStructFieldType(field.t.Type): // 嵌套结构体
			fieldType := field.t.Type
			if fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
			}
			nestedStructs := make([]structField, fieldType.NumField())
			for i := 0; i < fieldType.NumField(); i++ {
				nestedStructs[i] = structField{
					t:          fieldType.Field(i),
					namePrefix: fieldName,
					envPrefix:  envPrefix,
				}
			}
			// 将新的结构体字段重新压栈
			structFields = append(nestedStructs, structFields...)
		case isNormalStructFieldType(field.t.Type): // 简单值
			defaultValue, err := getStructFieldDefault(field.t)
			if err != nil {
				return nil, fmt.Errorf("parse default value for field \"%s\" error: %w", fieldNameStr, err)
			}
			fields = append(fields, NewField(fieldNameStr, defaultValue, getStructFieldEnv(field.t, envPrefix)...))
		default: // 不支持
			kind := field.t.Type.Kind().String()
			if field.t.Type.Kind() == reflect.Ptr {
				kind = "ptr of " + field.t.Type.Elem().Kind().String()
			}
			return nil, fmt.Errorf("not allow field type: %s, name: %s", kind, fieldNameStr)
		}
	}

	return fields, nil
}
