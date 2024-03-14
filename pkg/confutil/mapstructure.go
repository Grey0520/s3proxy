package confutil

import (
	"encoding/json"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

// StringToMapHookFunc 添加将字符串解析到 map[string]interface{} 的钩子
func StringToMapHookFunc() mapstructure.DecodeHookFuncKind {
	return func(f, t reflect.Kind, data interface{}) (interface{}, error) {
		if f != reflect.String || t != reflect.Map {
			return data, nil
		}

		ret := map[string]interface{}{}
		err := json.Unmarshal([]byte(data.(string)), &ret)
		return ret, err
	}
}

// StringToSliceHookFunc 添加将字符串解析到 []interface{} 的钩子
func StringToSliceHookFunc() mapstructure.DecodeHookFuncKind {
	return func(f, t reflect.Kind, data interface{}) (interface{}, error) {
		if f != reflect.String || t != reflect.Slice {
			return data, nil
		}

		var ret []interface{}
		err := json.Unmarshal([]byte(data.(string)), &ret)
		return ret, err
	}
}
