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
func (array Array) Append(ss ...Service) Service {
	return append(array, Array(ss))
}

// Go .
func (array Array) Go(ss ...Service) Service {
	return append(array, Set(ss))
}

// Start the service by running each hook's OnStart method.
func (array Array) Start() (err error) {
	for _, service := range array {
		if err = service.Start(); err != nil {
			break
		}
	}

	return
}

// Stop the service by running each hook's OnStop method.
func (array Array) Stop() (err error) {
	for _, service := range array.reverse() {
		if err = service.Stop(); err != nil {
			break
		}
	}

	return
}

func (array Array) reverse() Array {
	for i := len(array)/2 - 1; i >= 0; i-- {
		opp := len(array) - 1 - i
		array[i], array[opp] = array[opp], array[i]
	}

	return array
}
