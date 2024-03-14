package confutil

import (
	"reflect"
	"testing"
	"time"
)

// TestExpendConfigStruct 测试 ExpendConfigStruct 方法
func TestExpendConfigStruct(t *testing.T) {
	type sysConf struct {
		OS   string
		Arch string `default:"amd64" env:"ARCH"`
	}
	type config struct {
		Debug bool `env:"DEBUG"`
		User  struct {
			Name  string  `env:"NAME,^USERNAME" default:"anonymous"`
			Email string  `env:"EMAIL"`
			Age   uint8   `env:"AGE" default:"23"`
			Phone *string `env:"PHONE"`
		} `envPrefix:"USER_"`
		System struct {
			sysConf `envPrefix:"PLATFORM_" mapstructure:",squash"`
			Root    *bool `env:"ROOT" default:"true"`
		} `envPrefix:"SYSTEM_"`
		Admins     []string `env:"ADMINS" default:"[\"admin1\", \"admin2\"]"`
		SampleRate float64  `default:"0.233"`
		Version    map[string]string
		Timeout    time.Duration `default:"1h"`
	}

	fields, err := ExpendConfigStruct(&config{})
	if err != nil {
		t.Errorf("expend struct of config error: %s", err)
		return
	}

	trueVar := true
	expectRet := []Field{
		NewField("Debug", false, "DEBUG"),
		NewField("User.Name", "anonymous", "USER_NAME", "USERNAME"),
		NewField("User.Email", "", "USER_EMAIL"),
		NewField("User.Age", uint8(23), "USER_AGE"),
		NewField("User.Phone", (*string)(nil), "USER_PHONE"),
		NewField("System.OS", ""),
		NewField("System.Arch", "amd64", "SYSTEM_PLATFORM_ARCH"),
		NewField("System.Root", &trueVar, "SYSTEM_ROOT"),
		NewField("Admins", []string{"admin1", "admin2"}, "ADMINS"),
		NewField("SampleRate", 0.233),
		NewField("Version", map[string]string(nil)),
		NewField("Timeout", time.Hour),
	}
	if len(fields) != len(expectRet) {
		t.Errorf("unexpected len of ret: %d != %d", len(fields), len(expectRet))
		return
	}
	for i, field := range fields {
		if !reflect.DeepEqual(field, expectRet[i]) {
			t.Errorf("unexpected ret for field %d: %v != %v", i, field, expectRet[i])
		}
	}
}
