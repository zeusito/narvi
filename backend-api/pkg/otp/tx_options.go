package otp

import "github.com/uptrace/bun"

// TxOptions is used to support transactions, leave empty if not using transactions
type TxOptions struct {
	BunTx bun.IDB
}
