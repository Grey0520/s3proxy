package confutil

import (
	"fmt"
	"os"
	"testing"
)

// TestConfigManager 测试 ConfigManager
// TODO: 有待完善
func TestConfigManager(t *testing.T) {
	mgr := &ConfigManager{
		Items: []ConfigItem{{
			Name:       "field1",
			BindingEnv: "FIELD1",
			DefaultVal: []string{"val1", "val2"},
		}},
	}
	settings := struct {
		Field1 []string
	}{}
	_ = os.Setenv("FIELD1", "ggg hhh")
	_ = mgr.Init()
	_ = mgr.Unmarshal(&settings)
	fmt.Printf("%#v\n", mgr.Viper().Get("field1"))
	fmt.Printf("%#v\n", settings)
}
