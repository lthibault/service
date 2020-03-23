package serviceproc

import "github.com/jbenet/goprocess"

// Go .
func Go(f goprocess.ProcessFunc) *Service {
	return &Service{f: f}
}

// Service .
type Service struct {
	f goprocess.ProcessFunc
	p goprocess.Process
}

// Start the process
func (s *Service) Start() error {
	s.p = goprocess.Go(s.f)
	return nil
}

// Stop the process
func (s *Service) Stop() error {
	return s.p.Close()
}
