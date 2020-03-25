package service_test

import (
	"testing"

	service "github.com/lthibault/service/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArray(t *testing.T) {
	log := new(intlog)

	svc := service.With(
		log.WithCtr(1, -1),
		log.WithCtr(2, -2),
	)

	require.NoError(t, svc.Start())
	assert.Equal(t, intlog{1, 2}, *log)

	// Check that deferred-ordering of services is enforced
	require.NoError(t, svc.Stop())
	assert.Equal(t, intlog{1, 2, -2, -1}, *log)
}

func TestTree(t *testing.T) {
	log := new(intlog)

	svc := service.With(
		service.With(
			log.WithCtr(1, -1),
			log.WithCtr(2, -2),
		),
		log.WithCtr(3, -3),
	)

	require.NoError(t, svc.Start())
	assert.Equal(t, intlog{1, 2, 3}, *log)

	// N.B.:  we check that deferred-ordering is enforced
	require.NoError(t, svc.Stop())
	assert.Equal(t, intlog{1, 2, 3, -3, -2, -1}, *log)
}
