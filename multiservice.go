package service

import "go.uber.org/multierr"

// MultiService is a collection of hooks that are run in-order on startup, and in
// reverse order (deferred-order) on shutdown.
type MultiService []Service

// Append a hook to the MultiService
func (ms *MultiService) Append(s Service) *MultiService {
	*ms = append(*ms, s)
	return ms
}

// Start the service by running each hook's OnStart method.
func (ms MultiService) Start() (err error) {
	started := make(MultiService, 0, len(ms))
	for _, service := range ms {
		if err = service.Start(); err != nil {
			return multierr.Combine(err, rollback(started))
		}

		started = append(started, service)
	}

	return
}

// Stop the service by running each hook's OnStop method.
func (ms MultiService) Stop() error {
	return rollback(ms)
}

func (ms MultiService) reverse() MultiService {
	for i := len(ms)/2 - 1; i >= 0; i-- {
		opp := len(ms) - 1 - i
		ms[i], ms[opp] = ms[opp], ms[i]
	}

	return ms
}

func rollback(ms MultiService) (err error) {
	for _, service := range ms.reverse() {
		err = multierr.Append(err, service.Stop())
	}

	return
}
