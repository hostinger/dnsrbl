package hbl

import (
	"context"
)

type Service interface {
	Delete(ctx context.Context, ip string) error
	Check(ctx context.Context, name, ip string) (interface{}, error)
	Block(ctx context.Context, address *Address) error
	Allow(ctx context.Context, address *Address) error
	Unblock(ctx context.Context, address *Address) error
	GetOne(ctx context.Context, ip string) (*Address, error)
	GetAll(ctx context.Context) ([]*Address, error)
	SyncOne(ctx context.Context, ip string) error
	SyncAll(ctx context.Context) error
}
