package hbl

import (
	"context"

	"github.com/hostinger/hbl/pkg/alerters"
	"github.com/hostinger/hbl/pkg/checkers"
	"github.com/hostinger/hbl/pkg/endpoints"
)

type Service interface {
	Delete(ctx context.Context, ip string) error
	Check(ctx context.Context, name, ip string) (interface{}, error)
	Block(ctx context.Context, address *Address) error
	Allow(ctx context.Context, address *Address) error
	Unblock(ctx context.Context, address *Address) error
	GetOne(ctx context.Context, ip string) (*Address, error)
	GetAll(ctx context.Context) ([]*Address, error)
}

type service struct {
	Repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		Repository: repository,
	}
}

func (s *service) Unblock(ctx context.Context, address *Address) error {
	if err := endpoints.ExecuteOnAll(ctx, address.IP, "Unblock"); err != nil {
		return err
	}
	alerters.AlertOnAll(ctx,
		&alerters.Alert{IP: address.IP,
			Action: address.Action, Comment: address.Comment},
	)
	return nil
}

func (s *service) Block(ctx context.Context, address *Address) error {
	if err := s.Repository.CreateAddress(ctx, address); err != nil {
		return err
	}
	if err := endpoints.ExecuteOnAll(ctx, address.IP, "Block"); err != nil {
		return err
	}
	alerters.AlertOnAll(ctx,
		&alerters.Alert{IP: address.IP,
			Action: address.Action, Comment: address.Comment},
	)
	return nil
}

func (s *service) Allow(ctx context.Context, address *Address) error {
	if err := s.Repository.CreateAddress(ctx, address); err != nil {
		return err
	}
	alerters.AlertOnAll(ctx,
		&alerters.Alert{IP: address.IP,
			Action: address.Action, Comment: address.Comment},
	)
	return nil
}

func (s *service) Delete(ctx context.Context, ip string) error {
	return s.Repository.DeleteAddress(context.Background(), ip)
}

func (s *service) GetOne(ctx context.Context, ip string) (*Address, error) {
	return s.Repository.GetAddress(context.Background(), ip)
}

func (s *service) GetAll(ctx context.Context) ([]*Address, error) {
	return s.Repository.GetAddresses(context.Background())
}

func (s *service) Check(ctx context.Context, name, ip string) (interface{}, error) {
	return checkers.CheckOnOne(ctx, ip, name)
}
