package service

import "golang.org/x/sync/errgroup"

// Group services into a multiservice that starts/stops each member concurrently, and
// without synchronization.
//
// It is useful for large numbers of independent services, such as those constituting a
// worker pool.
type Group []Service

// Start each service in its own goroutine.  There is no synchronization.
func (g Group) Start() error {
	var xg errgroup.Group
	for _, service := range g {
		xg.Go(service.Start)
	}
	return xg.Wait()
}

// Stop each service in its own goroutine.  There is no synchronization.
func (g Group) Stop() error {
	var xg errgroup.Group
	for _, service := range g {
		xg.Go(service.Stop)
	}
	return xg.Wait()
}

// Add a service
func (g *Group) Add(s Service) *Group {
	*g = append(*g, s)
	return g
}
