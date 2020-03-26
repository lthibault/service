package service

import "golang.org/x/sync/errgroup"

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
	var g errgroup.Group
	for _, service := range set {
		g.Go(service.Start)
	}
	return g.Wait()
}

// Stop each service in its own goroutine.  There is no synchronization.
func (set Set) Stop() error {
	var g errgroup.Group
	for _, service := range set {
		g.Go(service.Stop)
	}
	return g.Wait()
}
