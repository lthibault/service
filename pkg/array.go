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
func (ss Array) Append(s Service) Array {
	return append(ss, s)
}

// Start the service by running each hook's OnStart method.
func (ss Array) Start() (err error) {
	for _, service := range ss {
		if err = service.Start(); err != nil {
			break
		}
	}

	return
}

// Stop the service by running each hook's OnStop method.
func (ss Array) Stop() (err error) {
	for _, service := range ss.reverse() {
		if err = service.Stop(); err != nil {
			break
		}
	}

	return
}

func (ss Array) reverse() Array {
	for i := len(ss)/2 - 1; i >= 0; i-- {
		opp := len(ss) - 1 - i
		ss[i], ss[opp] = ss[opp], ss[i]
	}

	return ss
}
