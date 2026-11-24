package config

import "github.com/harluo/di"

func init() {
	di.New().Instance().Put(
		newWrapper,
		newYaothink,
		newRate,
	).Build().Apply()
}
