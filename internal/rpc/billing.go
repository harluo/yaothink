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

type Billing struct {
	client server.TokenApiClient
	id     id.Generator
	logger log.Logger
}

func newBilling(client *grpc.Client, id id.Generator, logger log.Logger) *Billing {
	return &Billing{
		client: server.NewTokenApiClient(client.Connection(internal.NameBilling)),
		id:     id,
		logger: logger,
	}
}

func (b *Billing) Token(
	ctx context.Context,
	account, module uint64, model string,
	input, completion, read, creation uint32,
) (result *core.BillingResult, err error) {
	total := input + completion + read + creation
	if total == 0 {
		b.logger.Debug("无需计费", field.New("token", 0))
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
	b.logger.Debug("开始计费", fields[0], fields[1:]...)

	req := new(token.TokenCallReq)
	req.Account = account
	req.Module = module
	req.Model = model
	req.TokenInput = input
	req.TokenCompletion = completion
	req.TokenCacheRead = read
	req.TokenCacheCreation = creation
	if next, sie := b.setId(req); sie != nil {
		err = sie
	} else if rsp, ce := b.client.RecordTokenCall(ctx, req); ce != nil {
		err = ce
	} else if !rsp.Success {
		b.logger.Warn("计费失败", fields[0], fields[1:]...)
	} else {
		result = new(core.BillingResult)
		result.Id = next
		result.Balance = rsp.Balance
		result.Amount = rsp.CurrentRequestFee
		b.logger.Debug("计费成功", fields[0], fields[1:]...)
	}

	return
}

func (b *Billing) setId(req *token.TokenCallReq) (id uint64, err error) {
	if value, ne := b.id.Next(); ne != nil {
		err = ne
	} else {
		id = value.Get()
		req.RequestId = id
	}

	return
}
