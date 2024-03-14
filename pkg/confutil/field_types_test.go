package confutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"gopkg.in/yaml.v3"
	k8syaml "sigs.k8s.io/yaml"
)

// TestSecretString 测试 SecretString
func TestSecretString(t *testing.T) {
	raw := "abc"
	ss := SecretString("abc")

	// 检查原始值
	//goland:noinspection GoBoolExpressions
	if string(ss) != raw {
		t.Errorf("unexpected raw value (string(obj)): \"%s\" != \"%s\"", string(ss), raw)
		return
	}
	if ss.Raw() != raw {
		t.Errorf("unexpected raw value (obj.Raw()): \"%s\" != \"%s\"", ss.Raw(), raw)
		return
	}

	expectedShowingValue := "****** (sha256:ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad)"

	// 检查打印输出
	if fmt.Sprintf("%s", ss) != expectedShowingValue { //nolint:gosimple
		t.Errorf("unexpected showing value (%%s): \"%s\" != \"%s\"", ss, expectedShowingValue)
	}
	if fmt.Sprintf("%v", ss) != expectedShowingValue {
		t.Errorf("unexpected showing value (%%v): \"%v\" != \"%s\"", ss, expectedShowingValue)
	}

	// 检查序列化结果
	k8syamlVal, err := k8syaml.Marshal(ss)
	if err != nil {
		t.Errorf("k8s yaml marshal error: %s", err)
	}
	k8syamlVal = bytes.Trim(bytes.TrimSpace(k8syamlVal), "'\"")
	if string(k8syamlVal) != expectedShowingValue {
		t.Errorf("unexpected marshal result (k8s yaml): \"%s\" != \"%s\"", k8syamlVal, expectedShowingValue)
	}
	yamlVal, err := yaml.Marshal(ss)
	if err != nil {
		t.Errorf("yaml marshal error: %s", err)
	}
	yamlVal = bytes.Trim(bytes.TrimSpace(yamlVal), "'\"")
	if string(yamlVal) != expectedShowingValue {
		t.Errorf("unexpected marshal result (yaml): \"%s\" != \"%s\"", yamlVal, expectedShowingValue)
	}
	jsonVal, err := json.Marshal(ss)
	if err != nil {
		t.Errorf("json marshal error: %s", err)
	}
	if string(jsonVal) != "\""+expectedShowingValue+"\"" {
		t.Errorf("unexpected marshal result (json): \"%s\" != \"%s\"", jsonVal, expectedShowingValue)
	}
}
