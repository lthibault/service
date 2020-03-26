package service_test

import (
	"errors"
	"testing"

	service "github.com/lthibault/service/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArray(t *testing.T) {
	t.Run("Flat", func(t *testing.T) {
		log := new(intlog)

		svc := service.With(
			log.WithEntry(1, -1),
			log.WithEntry(2, -2),
		)

		check(t, log, svc,
			intlog{1, 2},
			intlog{1, 2, -2, -1})
	})

	t.Run("Tree", func(t *testing.T) {
		log := new(intlog)

		svc := service.With(
			service.With(
				log.WithEntry(1, -1),
				log.WithEntry(2, -2),
			),
			log.WithEntry(3, -3),
		)

		check(t, log, svc,
			intlog{1, 2, 3},
			intlog{1, 2, 3, -3, -2, -1})

	})
}

func TestArrayError(t *testing.T) {
	t.Run("Start", func(t *testing.T) {
		log := new(intlog)

		svc := service.With(
			log.WithEntry(1, -1),
			log.WithEntry(2, -2),
			log.WithError(errors.New("fail"), nil), // should not appear in log
			log.WithEntry(3, -3),                   // should not appear in log
		)

		require.EqualError(t, svc.Start(), "fail")
		assert.Equal(t, intlog{1, 2, -2, -1}, *log)
	})

	t.Run("Stop", func(t *testing.T) {
		log := new(intlog)

		svc := service.With(
			log.WithEntry(1, -1),
			log.WithEntry(2, -2),
			log.WithError(nil, errors.New("fail")), // should not appear in log
			log.WithEntry(3, -3),                   // should not appear in log
		)

		require.NoError(t, svc.Start())
		assert.Equal(t, intlog{1, 2, 3}, *log)

		// N.B.:  we check that deferred-ordering is enforced
		require.EqualError(t, svc.Stop(), "fail")
		assert.Equal(t, intlog{1, 2, 3, -3, -2, -1}, *log)
	})
}

func check(t *testing.T, log *intlog, svc service.Service, running, stopped intlog) {
	require.NoError(t, svc.Start())
	assert.Equal(t, running, *log)

	// N.B.:  we check that deferred-ordering is enforced
	require.NoError(t, svc.Stop())
	assert.Equal(t, stopped, *log)
}
