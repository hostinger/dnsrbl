package checkers

import (
	"context"
	"sync"
)

type Checker interface {
	Name() string
	Check(ctx context.Context, ip string) (interface{}, error)
}

var (
	checkersMu = new(sync.Mutex)
	checkers   = map[string]Checker{}
)

func CheckOnOne(ctx context.Context, ip, name string) (interface{}, error) {
	if checker, ok := checkers[name]; ok {
		result, err := checker.Check(ctx, ip)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, nil
}

func Register(checker Checker) {
	checkersMu.Lock()
	defer checkersMu.Unlock()
	if _, ok := checkers[checker.Name()]; !ok {
		checkers[checker.Name()] = checker
	}
}
