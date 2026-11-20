package rpc

import (
	"github.com/harluo/di"
	_ "github.com/harluo/http"
)

func init() {
	di.New().Instance().Put(
		newWallet,
		newRate,
	).Build().Apply()
}
