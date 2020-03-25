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
func (array Array) Append(ss ...Service) Group {
	return append(array, Array(ss))
}

// Go .
func (array Array) Go(ss ...Service) Group {
	return array.Append(Set(ss))
}

// Defer .
func (array Array) Defer(ss ...Service) Group {
	return deferredArray{
		d: Array(ss),
		a: array,
	}
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

type deferredArray struct {
	a, d Array
}

func (d deferredArray) Append(ss ...Service) Group {
	d.a = append(d.a, Array(ss))
	return d
}

func (d deferredArray) Go(ss ...Service) Group {
	d.a = append(d.a, Set(ss))
	return d
}

func (d deferredArray) Defer(ss ...Service) Group {
	d.d = append(d.d, Array(ss))
	return d
}

// Start the array, then the deferred service.
func (d deferredArray) Start() (err error) {
	if err = d.a.Start(); err == nil {
		err = d.d.Start()
	}

	return
}

// Stop the array, then the deferred service.
func (d deferredArray) Stop() (err error) {
	if err = d.d.Stop(); err == nil {
		err = d.a.Stop()
	}

	return
}
