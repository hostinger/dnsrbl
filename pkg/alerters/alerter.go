package alerters

import (
	"context"
	"sync"
)

type Alert struct {
	IP      string
	Action  string
	Author  string
	Comment string
}

type Alerter interface {
	Name() string
	Alert(ctx context.Context, alert *Alert)
}

var (
	alertersMu = new(sync.Mutex)
	alerters   = map[string]Alerter{}
)

func AlertOnAll(ctx context.Context, alert *Alert) {
	for _, alerter := range alerters {
		alerter.Alert(ctx, alert)
	}
}

func AlertOnOne(ctx context.Context, alert *Alert, name string) {
	if _, ok := alerters[name]; ok {
		alerters[name].Alert(ctx, alert)
	}
}

func Register(alerter Alerter) {
	alertersMu.Lock()
	defer alertersMu.Unlock()
	if _, ok := alerters[alerter.Name()]; !ok {
		alerters[alerter.Name()] = alerter
	}
}
