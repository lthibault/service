package service_test

import (
	"errors"
	"testing"

	service "github.com/lthibault/service/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSet(t *testing.T) {
	log := new(ctrlog)
	svc := service.Go(
		log.WithIncr(1, -1),
		log.WithIncr(1, -1),
		log.WithIncr(1, -1),
	)

	checkctr(t, log, svc, 3, 0)
}

func TestSetError(t *testing.T) {
	t.Run("Start", func(t *testing.T) {
		log := new(ctrlog)

		svc := service.Go(
			log.WithIncr(1, -1),
			log.WithIncr(1, -1),
			log.WithError(errors.New("fail"), nil), // should not appear in log
			log.WithIncr(1, -1),                    // should not appear in log
		)

		require.EqualError(t, svc.Start(), "fail")
		assert.Equal(t, ctrlog(0), *log)
	})

	t.Run("Stop", func(t *testing.T) {
		log := new(ctrlog)

		svc := service.Go(
			log.WithIncr(1, -1),
			log.WithIncr(1, -1),
			log.WithError(nil, errors.New("fail")), // should not appear in log
			log.WithIncr(1, -1),                    // should not appear in log
		)

		require.NoError(t, svc.Start())
		assert.Equal(t, ctrlog(3), *log)

		// N.B.:  we check that deferred-ordering is enforced
		require.EqualError(t, svc.Stop(), "fail")
		assert.Equal(t, ctrlog(0), *log)
	})
}

func checkctr(t *testing.T, log *ctrlog, svc service.Service, started, stopped int32) {
	require.NoError(t, svc.Start())
	assert.Equal(t, started, int32(*log))

	require.NoError(t, svc.Stop())
	assert.Equal(t, stopped, int32(*log))
}
