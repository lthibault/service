package service

// Service .
type Service interface {
	Start() error
	Stop() error
}

// With bundles individual services into a group.
func With(ss ...Service) Array {
	return Array(ss)
}

// Go bundles individual services into a set.
func Go(ss ...Service) Set {
	return Set(ss)
}

// Hook is a basic service.
type Hook struct {
	OnStart, OnStop func() error
}

// Start is called if OnStart is set.
func (h Hook) Start() (err error) {
	if h.OnStart != nil {
		err = h.OnStart()
	}

	return
}

// Stop is called if OnStop is set.
func (h Hook) Stop() (err error) {
	if h.OnStop != nil {
		err = h.OnStop()
	}

	return
}

// Append .
func (h Hook) Append(ss ...Service) Array {
	return With(h).Append(ss...)
}

// Go .
func (h Hook) Go(ss ...Service) Set {
	return With(h).Go(ss...)
}
