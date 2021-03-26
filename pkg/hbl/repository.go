package hbl

import (
	"context"
)

type Repository interface {
	GetAddress(ctx context.Context, ip string) (*Address, error)
	CreateAddress(ctx context.Context, address *Address) error
	GetAddresses(ctx context.Context) ([]*Address, error)
	DeleteAddress(ctx context.Context, ip string) error
}
