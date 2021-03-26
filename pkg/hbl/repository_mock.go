// +build !codeanalysis

package hbl

import (
	"context"
	"errors"
)

type mockRepository struct {
	db map[string]*Address
}

func NewMockRepository() Repository {
	return &mockRepository{
		db: make(map[string]*Address),
	}
}

func (r *mockRepository) CreateAddress(ctx context.Context, address *Address) error {
	if _, ok := r.db[address.IP]; !ok {
		r.db[address.IP] = address
		return nil
	}
	return errors.New("Address already exists")
}

func (r *mockRepository) DeleteAddress(ctx context.Context, ip string) error {
	if _, ok := r.db[ip]; !ok {
		return errors.New("Address doesn't exist")
	}
	delete(r.db, ip)
	return nil
}

func (r *mockRepository) GetAddress(ctx context.Context, ip string) (*Address, error) {
	if _, ok := r.db[ip]; !ok {
		return nil, errors.New("Address doesn't exist")
	}
	return r.db[ip], nil
}

func (r *mockRepository) GetAddresses(ctx context.Context) ([]*Address, error) {
	var addresses []*Address
	for _, address := range r.db {
		addresses = append(addresses, address)
	}
	return addresses, nil
}
