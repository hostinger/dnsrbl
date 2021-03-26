package hbl

import (
	"context"

	"github.com/hostinger/hbl/pkg/alerters"
	"github.com/hostinger/hbl/pkg/checkers"
	"github.com/hostinger/hbl/pkg/endpoints"
	"github.com/hostinger/hbl/pkg/logger"
	"go.uber.org/zap"
)

type service struct {
	logger     logger.Logger
	repository Repository
}

func NewDefaultService(l logger.Logger, r Repository) Service {
	return &service{
		repository: r,
		logger:     l,
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
	if err := s.repository.CreateAddress(ctx, address); err != nil {
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
	if err := s.repository.CreateAddress(ctx, address); err != nil {
		return err
	}
	alerters.AlertOnAll(ctx,
		&alerters.Alert{IP: address.IP,
			Action: address.Action, Comment: address.Comment},
	)
	return nil
}

func (s *service) Delete(ctx context.Context, ip string) error {
	return s.repository.DeleteAddress(ctx, ip)
}

func (s *service) GetOne(ctx context.Context, ip string) (*Address, error) {
	return s.repository.GetAddress(ctx, ip)
}

func (s *service) GetAll(ctx context.Context) ([]*Address, error) {
	return s.repository.GetAddresses(ctx)
}

func (s *service) Check(ctx context.Context, name, ip string) (interface{}, error) {
	return checkers.CheckOnOne(ctx, ip, name)
}

func (s *service) SyncOne(ctx context.Context, ip string) error {
	address, err := s.repository.GetAddress(ctx, ip)
	if err != nil {
		return err
	}
	if err := endpoints.ExecuteOnAll(ctx, address.IP, "Sync"); err != nil {
		return err
	}
	s.logger.Info("Synced address with all endpoints", zap.String("address", address.IP))
	return nil
}

func (s *service) SyncAll(ctx context.Context) error {
	addresses, err := s.repository.GetAddresses(ctx)
	if err != nil {
		return err
	}
	for _, address := range addresses {
		if err := endpoints.ExecuteOnAll(ctx, address.IP, "Sync"); err != nil {
			return err
		}
		s.logger.Info("Synced address with all endpoints", zap.String("address", address.IP))
	}
	return nil
}
