package tx

import "context"

// Manager is a transaction boundary abstraction. Infra will provide DB-backed impls.
type Manager interface {
    WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}


