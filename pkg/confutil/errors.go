package confutil

import "errors"

// ErrInitMultiTimes 配置管理器被初始化了多次
var ErrInitMultiTimes = errors.New("config manager initialized multiple times")
