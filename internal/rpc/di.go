package rpc

import (
	"github.com/harluo/di"
)

func init() {
	di.New().Instance().Put(
		newWallet,
	).Build().Apply()
}
