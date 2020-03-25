package serviceproc_test

import (
	"testing"
	"time"

	"github.com/jbenet/goprocess"
	"github.com/lthibault/service/pkg/serviceproc"
	"github.com/stretchr/testify/assert"
)

func TestServiceProc(t *testing.T) {
	var start, stop bool

	svc := serviceproc.Go(func(p goprocess.Process) {
		start = true
		<-p.Closing()
		stop = true
	})

	assert.NoError(t, svc.Start())
	time.Sleep(time.Millisecond)

	assert.NoError(t, svc.Stop())
	time.Sleep(time.Millisecond)

	assert.True(t, start, "did not start")
	assert.True(t, stop, "did not stop")
}
