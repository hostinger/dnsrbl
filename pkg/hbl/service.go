package hbl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/hostinger/dnsrbl/pkg/abuseipdb"
	"github.com/hostinger/dnsrbl/pkg/cloudflare"
	"github.com/hostinger/dnsrbl/pkg/dns"
)

var (
	ErrAddressExists    = errors.New("the IP address exists")
	ErrAddressNotExists = errors.New("the IP address doesn't exist")
	ErrAddressIsAllowed = errors.New("the IP address exists and is marked as allowed")
	ErrAddressIsBlocked = errors.New("the IP address exists and is marked as blocked")
)

type Service struct {
	abuseipdbClient *abuseipdb.Client
	dnsClient       *dns.Client
	cfClient        *cloudflare.Client
	AddressStore    *AddressStore
	MetadataStore   *MetadataStore
}

func NewService(addressStore *AddressStore, metadataStore *MetadataStore,
	cfClient *cloudflare.Client, abuseipdbClient *abuseipdb.Client, dnsClient *dns.Client) *Service {
	return &Service{
		abuseipdbClient: abuseipdbClient,
		MetadataStore:   metadataStore,
		AddressStore:    addressStore,
		dnsClient:       dnsClient,
		cfClient:        cfClient,
	}
}

func (s *Service) BlockAddress(address *Address) error {
	a, err := s.AddressStore.GetOne(context.Background(), address.IP)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == nil {
		if a.Action == ActionAllow {
			return ErrAddressIsAllowed
		}
		if a.Action == ActionBlock {
			return ErrAddressIsBlocked
		}
	}
	if err = s.onBlock(address); err != nil {
		return fmt.Errorf("failed to execute onBlock() actions: %s", err)
	}
	if err = s.AddressStore.Create(context.Background(), address); err != nil {
		return fmt.Errorf("failed to execute AddressStore.Create(): %s", err)
	}
	report, err := s.abuseipdbClient.Check(address.IP)
	if err != nil {
		log.Printf("failed to fetch AbuseIPDB metadata: %s", err)
	}
	metadata := AbuseIPDBMetadata{
		IP:                   address.IP,
		ISP:                  report.Data.Isp,
		UsageType:            report.Data.UsageType,
		CountryCode:          report.Data.CountryCode,
		TotalReports:         report.Data.TotalReports,
		LastReportedAt:       report.Data.LastReportedAt,
		NumDistinctUsers:     report.Data.NumDistinctUsers,
		AbuseConfidenceScore: report.Data.AbuseConfidenceScore,
	}
	if err := s.MetadataStore.Create(context.Background(), &metadata); err != nil {
		return err
	}
	return nil
}

func (s *Service) AllowAddress(address *Address) error {
	a, err := s.AddressStore.GetOne(context.Background(), address.IP)
	if err == nil {
		if a.Action == ActionAllow {
			return ErrAddressIsAllowed
		}
		if a.Action == ActionBlock {
			return ErrAddressIsBlocked
		}
	}
	if err := s.AddressStore.Create(context.Background(), address); err != nil {
		return fmt.Errorf("Failed to execute AddressStore.Create(): %s", err)
	}
	return nil
}

func (s *Service) DeleteAddress(ip string) error {
	a, err := s.AddressStore.GetOne(context.Background(), ip)
	if err != nil && err == sql.ErrNoRows {
		return ErrAddressNotExists
	}
	if a.Action == "Block" {
		if os.Getenv("ENVIRONMENT") != "DEV" {
			if err := s.cfClient.Unblock(ip); err != nil {
				return fmt.Errorf("Failed to execute CloudflareClient.Unblock(): %s", err)
			}
		}
		if err := s.dnsClient.Unblock(context.Background(), ip, "hostinger.rbl"); err != nil {
			return fmt.Errorf("Failed to execute DNSClient.Unblock(): %s", err)
		}
		if err := s.MetadataStore.Delete(context.Background(), ip); err != nil {
			return fmt.Errorf("Failed to execute MetadataStore.Delete(): %s", err)
		}
	}
	if err := s.AddressStore.Delete(context.Background(), ip); err != nil {
		return fmt.Errorf("Failed to execute AddressStore.Delete(): %s", err)
	}
	return nil
}

func (s *Service) GetAddress(ip string) (Address, error) {
	address, err := s.AddressStore.GetOne(context.Background(), ip)
	if err != nil && err == sql.ErrNoRows {
		return Address{}, ErrAddressNotExists
	}
	metadata, err := s.MetadataStore.GetOne(context.Background(), ip)
	if err != nil && err != sql.ErrNoRows {
		return Address{}, err
	}
	address.Metadata.AbuseIPDBMetadata = metadata
	return address, nil
}

func (s *Service) GetAddresses() ([]Address, error) {
	addresses, err := s.AddressStore.GetAll(context.Background())
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

func (s *Service) onBlock(address *Address) error {
	if os.Getenv("ENVIRONMENT") != "DEV" {
		if err := s.cfClient.Block(address.IP); err != nil {
			return fmt.Errorf("failed to execute CloudflareClient.Block(): %s", err)
		}
	}
	if err := s.dnsClient.Block(context.Background(), address.IP, "hostinger.rbl"); err != nil {
		return fmt.Errorf("failed to execute DNSClient.Block(): %s", err)
	}
	return nil
}
