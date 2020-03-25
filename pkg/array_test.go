package service_test

import (
	"testing"

	service "github.com/lthibault/service/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMultiService(t *testing.T) {
	log := new(intlog)

	svc := service.Array{}.
		Append(log.WithCtr(1, -1)).
		Append(log.WithCtr(2, -2))

	require.NoError(t, svc.Start())
	assert.Equal(t, intlog{1, 2}, *log)

	// Check that deferred-ordering of services is enforced
	require.NoError(t, svc.Stop())
	assert.Equal(t, intlog{1, 2, -2, -1}, *log)
}

func TestHierarchicalMultiService(t *testing.T) {
	log := new(intlog)

	svc := service.Array{}.
		Append(service.Array{
			log.WithCtr(1, -1),
			log.WithCtr(2, -2),
		}).
		Append(
			log.WithCtr(3, -3),
		)

	require.NoError(t, svc.Start())
	assert.Equal(t, intlog{1, 2, 3}, *log)

	// N.B.:  we check that deferred-ordering is enforced
	require.NoError(t, svc.Stop())
	assert.Equal(t, intlog{1, 2, 3, -3, -2, -1}, *log)
}

// func TestDefer(t *testing.T) {
// 	log := new(intlog)

// 	svc := service.Array{}.
// 		Append(log.WithCtr(1, -1)).
// 		Defer(log.WithCtr(2, -2)).
// 		Append(log.WithCtr(3, -3))

// 	require.NoError(t, svc.Start())
// 	assert.Equal(t, intlog{1, 3, 2}, *log)

// 	// N.B.:  we check that deferred-ordering is enforced
// 	require.NoError(t, svc.Stop())
// 	assert.Equal(t, intlog{1, 3, 2, -2, -3, -1}, *log)
// }

type intlog []int

func (log *intlog) WithCtr(start, stop int) service.Hook {
	return service.New(func() error {
		*log = append(*log, start)
		return nil
	}, func() error {
		*log = append(*log, stop)
		return nil
	})
}
