package service

import "golang.org/x/sync/errgroup"

// Parallel multiservice starts/stops each service in a separate goroutine, without
// synchronization.  It is useful for large numbers of independent services, such as
// those constituting a worker pool.
type Parallel []Service

// Start each service in its own goroutine.  There is no synchronization.
func (p Parallel) Start() error {
	var g errgroup.Group
	for _, service := range p {
		g.Go(service.Start)
	}
	return g.Wait()
}

// Stop each service in its own goroutine.  There is no synchronization.
func (p Parallel) Stop() error {
	var g errgroup.Group
	for _, service := range p {
		g.Go(service.Stop)
	}
	return g.Wait()
}
