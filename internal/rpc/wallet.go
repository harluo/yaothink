package rpc

import (
	"context"

	"gitea.com/wallety/protocol/core/api/token"
	"gitea.com/wallety/protocol/server"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
	"github.com/goexl/id"
	"github.com/goexl/log"
	"github.com/harluo/grpc"
	"github.com/harluo/yaothink/internal/core"
	"github.com/harluo/yaothink/internal/rpc/internal"
)

type Wallet struct {
	client server.TokenApiClient
	id     id.Generator
	logger log.Logger
}

func newWallet(client *grpc.Client, id id.Generator, logger log.Logger) *Wallet {
	return &Wallet{
		client: server.NewTokenApiClient(client.Connection(internal.NameBilling)),
		id:     id,
		logger: logger,
	}
}

func (w *Wallet) Token(
	ctx context.Context,
	account, module uint64, model string,
	input, completion, read, creation uint32,
) (result *core.BillingResult, err error) {
	total := input + completion + read + creation
	if total == 0 {
		w.logger.Debug("无需计费", field.New("token", 0))
	}
	if total == 0 {
		return
	}

	fields := gox.Fields[any]{
		field.New("account", account),
		field.New("module", module),
		field.New("model", model),
		field.New("input", input),
		field.New("completion", completion),
		field.New("read", read),
		field.New("creation", creation),
	}
	w.logger.Debug("开始计费", fields[0], fields[1:]...)

	req := new(token.TokenCallReq)
	req.Account = account
	req.Module = module
	req.Model = model
	req.TokenInput = input
	req.TokenCompletion = completion
	req.TokenCacheRead = read
	req.TokenCacheCreation = creation
	if next, sie := w.setId(req); sie != nil {
		err = sie
	} else if rsp, ce := w.client.RecordTokenCall(ctx, req); ce != nil {
		err = ce
	} else {
		result = w.setResult(next, rsp, fields)
	}

	return
}

func (w *Wallet) setResult(id uint64, rsp *token.TokenCallResp, fields gox.Fields[any]) (result *core.BillingResult) {
	if !rsp.Success {
		w.logger.Warn("计费失败", fields[0], fields[1:]...)
	} else {
		w.logger.Debug("计费成功", fields[0], fields[1:]...)
	}

	result = new(core.BillingResult)
	result.Id = id
	result.Success = rsp.Success
	result.Status = rsp.BalanceStatus

	return
}

func (w *Wallet) setId(req *token.TokenCallReq) (id uint64, err error) {
	if value, ne := w.id.Next(); ne != nil {
		err = ne
	} else {
		id = value.Get()
		req.RequestId = id
	}

	return
}
