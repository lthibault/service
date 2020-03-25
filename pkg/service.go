package service

// Service .
type Service interface {
	Start() error
	Stop() error
}

// Hook encapsulates startup and shutdown functions.
type Hook struct {
	OnStart, OnStop func() error
}

// Start call OnStart if it is defined.
func (h Hook) Start() error {
	if h.OnStart == nil {
		return nil
	}

	return h.OnStart()
}

// Stop calls OnStop if it is defined.
func (h Hook) Stop() error {
	if h.OnStop == nil {
		return nil
	}

	return h.OnStop()
}

// New Hook
func New(start, stop func() error) Hook {
	return Hook{
		OnStart: start,
		OnStop:  stop,
	}
}
