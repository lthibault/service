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

// With bundles individual services into a group.
func With(ss ...Service) Array {
	return Array(ss)
}

// Go bundles individual services into a set.
func Go(ss ...Service) Set {
	return Set(ss)
}
