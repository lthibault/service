package service_test

import (
	"testing"

	service "github.com/lthibault/service/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSet(t *testing.T) {
	log := new(ctrlog)
	svc := service.Go(
		log.WithCtr(1, -1),
		log.WithCtr(1, -1),
		log.WithCtr(1, -1),
	)

	require.NoError(t, svc.Start())
	assert.Equal(t, ctrlog(3), *log)

	require.NoError(t, svc.Stop())
	assert.Equal(t, ctrlog(0), *log)

}
