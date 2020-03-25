package service

import "github.com/jbenet/goprocess"

type (
	// Proc .
	Proc interface {
		goprocess.Process
	}

	// ProcFunc .
	ProcFunc func(Proc)
)

// Service .
type Service interface {
	Start() error
	Stop() error
}

type process struct {
	start, stop func() error
}

func (p process) Append(ss ...Service) Array {
	return append(Array{p}, Array(ss)...)
}

func (p process) Go(ss ...Service) Set {
	return append(Set{p}, Set(ss)...)
}

func (p process) Start() error {
	if p.start == nil {
		return nil
	}

	return p.start()
}

func (p process) Stop() error {
	if p.start == nil {
		return nil
	}

	return p.stop()
}

// New service
func New(f ProcFunc) Service {
	var p Proc
	return process{
		start: func() error {
			p = goprocess.Background().Go(func(p goprocess.Process) {
				f(p)
			})

			return nil
		},
		stop: func() error {
			return p.Close()
		},
	}
}

// With bundles individual services into a group.
func With(ss ...Service) Array {
	return Array(ss)
}

// Go bundles individual services into a set.
func Go(ss ...Service) Set {
	return Set(ss)
}
