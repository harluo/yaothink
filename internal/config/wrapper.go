package config

import (
	"github.com/harluo/config"
)

// 组件配置包装器。
type Wrapper struct {
	Yaothink Yaothink `json:"yaothink"`
}

func newWrapper(getter config.Getter) (wrapper *Wrapper, err error) {
	wrapper = new(Wrapper)
	err = getter.Get(wrapper)

	return
}
