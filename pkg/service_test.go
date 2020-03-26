package service_test

import (
	"sync/atomic"

	service "github.com/lthibault/service/pkg"
)

type mockService struct {
	start, stop func() error
}

func (m mockService) Start() error {
	if m.start == nil {
		return nil
	}

	return m.start()
}

func (m mockService) Stop() error {
	if m.stop == nil {
		return nil
	}

	return m.stop()
}

func (m mockService) Append(ss ...service.Service) service.Service {
	return service.With(m).Append(ss...)
}

func (m mockService) Go(ss ...service.Service) service.Service {
	return service.With(m).Go(ss...)
}

type intlog []int

func (log *intlog) WithEntry(start, stop int) service.Service {
	return mockService{
		start: log.append(start),
		stop:  log.append(stop),
	}
}

func (log intlog) WithError(start, stop error) service.Service {
	return mockService{
		start: func() error { return start },
		stop:  func() error { return stop },
	}
}

func (log *intlog) append(i int) func() error {
	return func() error {
		*log = append(*log, i)
		return nil
	}
}

type ctrlog int32

func (log *ctrlog) WithIncr(start, stop int32) service.Service {
	return mockService{
		start: log.incr(start),
		stop:  log.incr(stop),
	}
}

func (log ctrlog) WithError(start, stop error) service.Service {
	return mockService{
		start: func() error { return start },
		stop:  func() error { return stop },
	}
}

func (log *ctrlog) incr(i int32) func() error {
	return func() error {
		atomic.AddInt32((*int32)(log), i)
		return nil
	}
}
