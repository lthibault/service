package service_test

import (
	"sync/atomic"

	service "github.com/lthibault/service/pkg"
)

type intlog []int

func (log *intlog) WithEntry(start, stop int) service.Service {
	return service.Hook{
		OnStart: log.append(start),
		OnStop:  log.append(stop),
	}
}

func (log intlog) WithError(start, stop error) service.Service {
	return service.Hook{
		OnStart: func() error { return start },
		OnStop:  func() error { return stop },
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
	return service.Hook{
		OnStart: log.incr(start),
		OnStop:  log.incr(stop),
	}
}

func (log ctrlog) WithError(start, stop error) service.Service {
	return service.Hook{
		OnStart: func() error { return start },
		OnStop:  func() error { return stop },
	}
}

func (log *ctrlog) incr(i int32) func() error {
	return func() error {
		atomic.AddInt32((*int32)(log), i)
		return nil
	}
}
