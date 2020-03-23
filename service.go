package service

import "context"

// Service .
type Service interface {
	Start(context.Context) error
	Stop(context.Context) error
}

// Hook encapsulates startup and shutdown functions.
type Hook struct {
	OnStart, OnStop func(context.Context) error
}

// Start call OnStart if it is defined.
func (h Hook) Start(ctx context.Context) error {
	if h.OnStart == nil {
		return nil
	}

	return h.OnStart(ctx)
}

// Stop calls OnStop if it is defined.
func (h Hook) Stop(ctx context.Context) error {
	if h.OnStop == nil {
		return nil
	}

	return h.OnStop(ctx)
}

// MultiService is a collection of hooks that are run in-order on both startup and shutdown.
type MultiService []Service

// Append a hook to the MultiService
func (ms *MultiService) Append(s Service) {
	*ms = append(*ms, s)
}

// Start the service by running each hook's OnStart method.
func (ms MultiService) Start(ctx context.Context) error {
	return ms.foreach(ctx, func(ctx context.Context, s Service) error {
		return s.Start(ctx)
	})
}

// Stop the service by running each hook's OnStop method.
func (ms MultiService) Stop(ctx context.Context) error {
	return ms.foreach(ctx, func(ctx context.Context, s Service) error {
		return s.Stop(ctx)
	})
}

func (ms MultiService) foreach(ctx context.Context, f func(context.Context, Service) error) (err error) {
	for _, service := range ms {
		if err = f(ctx, service); err != nil {
			break
		}
	}

	return
}
