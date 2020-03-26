package service

import "go.uber.org/multierr"

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
func (array Array) Append(ss ...Service) Array {
	return append(array, Array(ss)...)
}

// Go .
func (array Array) Go(ss ...Service) Set {
	return Go(array, Set(ss))
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
	var log = make(txlog, 0, len(array))

	for _, service := range array.reverse() {
		if err = multierr.Append(err, service.Stop()); err == nil {
			log = append(log, service)
		}
	}

	return log.MaybeRollback(err)
}

func (array Array) reverse() Array {
	for i := len(array)/2 - 1; i >= 0; i-- {
		opp := len(array) - 1 - i
		array[i], array[opp] = array[opp], array[i]
	}

	return array
}

type txlog []Service

func (log txlog) MaybeRollback(err error) error {
	if err != nil {
		return Array(log).Stop()
	}

	return nil
}
