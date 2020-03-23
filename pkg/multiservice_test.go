package service_test

import (
	"testing"

	service "github.com/lthibault/service/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMultiService(t *testing.T) {
	res := []int{}

	svc := new(service.MultiService).
		Append(service.Hook{
			OnStart: func() error {
				res = append(res, 1)
				return nil
			},
			OnStop: func() error {
				res = append(res, -1)
				return nil
			},
		}).
		Append(service.Hook{
			OnStart: func() error {
				res = append(res, 2)
				return nil
			},
			OnStop: func() error {
				res = append(res, -2)
				return nil
			},
		})

	require.NoError(t, svc.Start())

	assert.Equal(t, []int{1, 2}, res)

	require.NoError(t, svc.Stop())

	// N.B.:  we check that deferred-ordering is enforced
	assert.Equal(t, []int{1, 2, -2, -1}, res)
}

func TestHierarchicalMultiService(t *testing.T) {
	res := []int{}

	svc := new(service.MultiService).
		Append(service.MultiService{
			service.Hook{
				OnStart: func() error {
					res = append(res, 1)
					return nil
				},
				OnStop: func() error {
					res = append(res, -1)
					return nil
				},
			},
			service.Hook{
				OnStart: func() error {
					res = append(res, 2)
					return nil
				},
				OnStop: func() error {
					res = append(res, -2)
					return nil
				},
			},
		}).
		Append(service.Hook{
			OnStart: func() error {
				res = append(res, 3)
				return nil
			},
			OnStop: func() error {
				res = append(res, -3)
				return nil
			},
		})

	require.NoError(t, svc.Start())

	assert.Equal(t, []int{1, 2, 3}, res)

	require.NoError(t, svc.Stop())

	// N.B.:  we check that deferred-ordering is enforced
	assert.Equal(t, []int{1, 2, 3, -3, -2, -1}, res)
}
