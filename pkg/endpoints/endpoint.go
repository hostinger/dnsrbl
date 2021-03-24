package endpoints

import (
	"context"
	"sync"

	"github.com/pkg/errors"
)

type Endpoint interface {
	Name() string
	Sync(ctx context.Context, ip string) error
	Block(ctx context.Context, ip string) error
	Unblock(ctx context.Context, ip string) error
}

var (
	endpointsMu = new(sync.Mutex)
	endpoints   = map[string]Endpoint{}
)

func ExecuteOnAll(ctx context.Context, ip, action string) error {
	for _, endpoint := range endpoints {
		switch action {
		case "Block":
			if err := endpoint.Block(ctx, ip); err != nil {
				return errors.Wrapf(err, "Block failed on Endpoint '%s'", endpoint.Name())
			}
		case "Unblock":
			if err := endpoint.Unblock(ctx, ip); err != nil {
				return errors.Wrapf(err, "Unblock failed on Endpoint '%s'", endpoint.Name())
			}
		case "Sync":
			if err := endpoint.Sync(ctx, ip); err != nil {
				return errors.Wrapf(err, "Sync failed on Endpoint '%s'", endpoint.Name())
			}
		}
	}
	return nil
}

func ExecuteOnOne(ctx context.Context, ip, action, name string) error {
	if endpoint, ok := endpoints[name]; ok {
		switch action {
		case "Block":
			if err := endpoint.Block(ctx, ip); err != nil {
				return errors.Wrapf(err, "Block failed on Endpoint '%s'", endpoint.Name())
			}
		case "Unblock":
			if err := endpoint.Unblock(ctx, ip); err != nil {
				return errors.Wrapf(err, "Unblock failed on Endpoint '%s'", endpoint.Name())
			}
		case "Sync":
			if err := endpoint.Sync(ctx, ip); err != nil {
				return errors.Wrapf(err, "Sync failed on Endpoint '%s'", endpoint.Name())
			}
		}
	}
	return nil
}

func Register(endpoint Endpoint) {
	endpointsMu.Lock()
	defer endpointsMu.Unlock()
	if _, ok := endpoints[endpoint.Name()]; !ok {
		endpoints[endpoint.Name()] = endpoint
	}
}
