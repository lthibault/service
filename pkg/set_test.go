package service_test

import (
	"sync/atomic"
	"testing"

	service "github.com/lthibault/service/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSet(t *testing.T) {
	log := new(ctrlog)
	svc := service.Set{}.
		Add(log.WithCtr(1, -1)).
		Add(log.WithCtr(1, -1)).
		Add(log.WithCtr(1, -1))

	require.NoError(t, svc.Start())
	assert.Equal(t, ctrlog(3), *log)

	require.NoError(t, svc.Stop())
	assert.Equal(t, ctrlog(0), *log)

}

type ctrlog int32

func (log *ctrlog) WithCtr(start, stop int32) service.Service {
	return service.New(log.incr(start), log.incr(stop))
}

func (log *ctrlog) incr(i int32) func() error {
	return func() error {
		atomic.AddInt32((*int32)(log), i)
		return nil
	}
}
