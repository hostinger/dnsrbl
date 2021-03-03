package dnsrbl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/cloudflare/cloudflare-go"
)

type Service struct {
	cfAccount string
	cfClient  *cloudflare.API
	Config    *Config
	Store     *Store
}

func NewService(cfg *Config, store *Store, cfClient *cloudflare.API, cfAccount string) *Service {
	return &Service{
		cfAccount: cfAccount,
		cfClient:  cfClient,
		Store:     store,
		Config:    cfg,
	}
}

func (s *Service) BlockAddress(address Address) error {
	if s.IsAddressInAllowList(address.Address) {
		return errors.New("That IP address is in allow list.")
	}
	if err := s.Store.CreateAddress(context.Background(), address); err != nil {
		return err
	}
	if err := s.BlockIPAddressInCloudflare(address.Address); err != nil {
		log.Printf("Failed to execute BlockIPAddressInCloudflare: %s", err)
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
	if err := s.UnblockIPAddressInCloudflare(ip); err != nil {
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

func (s *Service) BlockIPAddressInCloudflare(ip string) error {
	if net.ParseIP(ip) == nil {
		return fmt.Errorf("Address '%s' is not a valid IP address.", ip)
	}
	rule := cloudflare.AccessRule{
		Mode: "block",
		Configuration: cloudflare.AccessRuleConfiguration{
			Target: "ip",
			Value:  ip,
		},
		Notes: "Created automatically by HBL API.",
	}
	response, err := s.cfClient.CreateAccountAccessRule(context.Background(), s.cfAccount, rule)
	if err != nil || response.Success == false {
		return err
	}
	return nil
}

func (s *Service) UnblockIPAddressInCloudflare(ip string) error {
	if net.ParseIP(ip) == nil {
		return fmt.Errorf("Address '%s' is not a valid IP address.", ip)
	}
	rule := cloudflare.AccessRule{
		Mode: "block",
		Configuration: cloudflare.AccessRuleConfiguration{
			Target: "ip",
			Value:  ip,
		},
		Notes: "Created automatically by HBL API.",
	}
	rules, err := s.cfClient.ListAccountAccessRules(context.Background(), s.cfAccount, rule, 1)
	if err != nil {
		return err
	}
	if rules.Count <= 0 || rules.Count > 1 {
		return fmt.Errorf("AccessRule for IP address '%s' was not found. ", ip)
	}
	response, err := s.cfClient.DeleteAccountAccessRule(context.Background(), s.cfAccount, rules.Result[0].ID)
	if err != nil || response.Success == false {
		return err
	}
	return nil
}
