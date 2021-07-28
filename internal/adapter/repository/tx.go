package repository

import "context"

type noopTransaction struct {
}

func NewNoopTransaction() *noopTransaction {
	return &noopTransaction{}
}

func (*noopTransaction) Transaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	return fn(ctx)
}
