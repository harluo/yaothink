package config

// 组件配置分组。
type Yaothink struct {
	Rate *Rate `json:"rate,omitempty" validate:"required"`
}

func newYaothink(wrapper *Wrapper) *Yaothink{
	return &wrapper.Yaothink
}
