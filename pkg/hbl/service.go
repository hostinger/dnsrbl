package hbl

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net"

	"github.com/hostinger/dnsrbl/pkg/abuseipdb"
	"github.com/hostinger/dnsrbl/pkg/cloudflare"
)

type Service struct {
	abuseipdbClient *abuseipdb.Client
	cfClient        *cloudflare.Client
	Config          *Config
	Store           *Store
}

func NewService(cfg *Config, store *Store,
	cfClient *cloudflare.Client, abuseipdbClient *abuseipdb.Client) *Service {
	return &Service{
		abuseipdbClient: abuseipdbClient,
		cfClient:        cfClient,
		Store:           store,
		Config:          cfg,
	}
}

func (s *Service) BlockAddress(address Address) error {
	if s.IsAddressInAllowList(address.IP) {
		return errors.New("That IP address is in allow list.")
	}
	if err := s.Store.CreateAddress(context.Background(), address); err != nil {
		return err
	}
	return nil
}

func (s *Service) UnblockAddress(ip string) error {
	_, err := s.GetAddress(ip)
	if err != nil && err == sql.ErrNoRows {
		return errors.New("That IP address isn't blocked.")
	}
	if err := s.Store.DeleteAddress(context.Background(), ip); err != nil {
		return err
	}
	if err := s.cfClient.UnblockIPAddress(ip); err != nil {
		log.Printf("Failed to execute UnblockIPAddressInCloudflare: %s", err)
	}
	return nil
}

func (s *Service) GetAddress(ip string) (Address, error) {
	address, err := s.Store.GetAddress(context.Background(), ip)
	if err != nil {
		return Address{}, err
	}
	return address, nil
}

func (s *Service) GetAddresses() ([]Address, error) {
	addresses, err := s.Store.GetAddresses(context.Background())
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

func (s *Service) IsAddressInAllowList(ip string) bool {
	if net.ParseIP(ip) == nil {
		return false
	}
	for _, item := range s.Config.AllowList {
		_, network, _ := net.ParseCIDR(item)
		if network.Contains(net.ParseIP(ip)) {
			return true
		}
	}
	return false
}
