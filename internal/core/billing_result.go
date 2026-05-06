package core

import (
	"gitea.com/wallety/protocol/core/api"
)

type BillingResult struct {
	Id      uint64            `json:"id,omitempty"`
	Success bool              `json:"success,omitempty"`
	Status  api.BalanceStatus `json:"status,omitempty"`
}
