package config

import "testing"

func TestLoadConfig(t *testing.T) {
	type args struct {
		configPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// 测试用例 1：传入空路径
		{
			name: "empty config path",
			args: args{
				configPath: "",
			},
			wantErr: false,
		},
		// 测试用例 2：传入不存在的文件路径
		{
			name: "non-existent config file",
			args: args{
				configPath: "./non-existent.yaml",
			},
			wantErr: true,
		},
		// 测试用例 3：传入错误格式的文件
		{
			name: "invalid config file format",
			args: args{
				configPath: "./invalid.json",
			},
			wantErr: true,
		},
		// 测试用例 4：传入正确的配置文件
		{
			name: "valid config file",
			args: args{
				configPath: "../../configs/config.yaml",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LoadConfig(tt.args.configPath); (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
