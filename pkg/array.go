package service

// Array of services whose start/stop hooks are run in deferred-order, i.e.:
//
// * Start hooks are run left-to-right
//
// * Stop hooks are run right-to-left
//
// This emulates the native Go `defer` semantics with respect to startup and shutdown
// methods.
type Array []Service

// Append a hook to the Array
func (ms Array) Append(s Service) Array {
	return append(ms, s)
}

// Start the service by running each hook's OnStart method.
func (ms Array) Start() (err error) {
	for _, service := range ms {
		if err = service.Start(); err != nil {
			break
		}
	}

	return
}

// Stop the service by running each hook's OnStop method.
func (ms Array) Stop() (err error) {
	for _, service := range ms.reverse() {
		if err = service.Stop(); err != nil {
			break
		}
	}

	return
}

func (ms Array) reverse() Array {
	for i := len(ms)/2 - 1; i >= 0; i-- {
		opp := len(ms) - 1 - i
		ms[i], ms[opp] = ms[opp], ms[i]
	}

	return ms
}
