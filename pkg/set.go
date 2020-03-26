package service

import (
	"sync"

	"go.uber.org/multierr"
	"golang.org/x/sync/errgroup"
)

// Set is a collection of services that are collectively started and stopped as a group.
// Set differs from Array in that services are started and stopped in parallel, without
// any synchronization.
//
// Use set when you have many independent services.  This is commonly found in worker
// pools.
type Set []Service

// Append .
func (set Set) Append(ss ...Service) Array {
	return With(set, Array(ss))
}

// Go runs a service in a Set.
func (set Set) Go(ss ...Service) Set {
	return append(set, Set(ss)...)
}

// Start each service in its own goroutine.  There is no synchronization.
func (set Set) Start() error {
	var mu sync.Mutex
	log := make(txlog, len(set))

	var g errgroup.Group
	for i, s := range set {
		g.Go(set.startService(&mu, log, i, s))
	}

	return log.MaybeRollbackAsync(g.Wait())
}

// Stop each service in its own goroutine.  There is no synchronization.
func (set Set) Stop() (err error) {
	var wg sync.WaitGroup
	wg.Add(len(set))

	var mu sync.Mutex
	for _, s := range set {
		go func(service Service) {
			defer wg.Done()

			if e := service.Stop(); e != nil {
				mu.Lock()
				defer mu.Unlock()

				err = multierr.Append(err, e)
			}
		}(s)
	}

	wg.Wait()
	return
}

func (Set) startService(mu *sync.Mutex, log txlog, i int, s Service) func() error {
	return func() (err error) {
		if err = s.Start(); err != nil {
			mu.Lock()
			defer mu.Unlock()

			log[i] = s
		}

		return
	}
}

func (log txlog) MaybeRollbackAsync(err error) error {
	if err != nil {
		started := log[:0]
		for _, service := range log {
			if service != nil {
				started = append(started, service)
			}
		}

		return Set(started).Stop()
	}

	return nil
}
