package service

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
	for _, service := range ms {
		if err = service.Start(); err != nil {
			break
		}
	}

	return
}

// Stop the service by running each hook's OnStop method.
func (ms MultiService) Stop() (err error) {
	for _, service := range ms.reverse() {
		if err = service.Stop(); err != nil {
			break
		}
	}

	return
}

func (ms MultiService) reverse() MultiService {
	for i := len(ms)/2 - 1; i >= 0; i-- {
		opp := len(ms) - 1 - i
		ms[i], ms[opp] = ms[opp], ms[i]
	}

	return ms
}
