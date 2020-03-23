package serviceproc

import "github.com/jbenet/goprocess"

// GoChild .
func GoChild(parent goprocess.Process, f goprocess.ProcessFunc) *Service {

}

// Service .
type Service struct {
	f goprocess.ProcessFunc
	p goprocess.Process
}

// Start the process
func (s *Service) Start() error {
	s.p = s.f()
	return nil
}

// Stop the process
func (s *Service) Stop() error {
	return s.p.Close()
}
