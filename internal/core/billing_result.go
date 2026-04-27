package core

type BillingResult struct {
	Id      uint64 `json:"id,omitempty"`
	Balance int64  `json:"balance,omitempty"`
	Amount  uint64 `json:"amount,omitempty"`
}
