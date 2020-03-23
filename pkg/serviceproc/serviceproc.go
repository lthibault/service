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

// Process .
func (s *Service) Process() goprocess.Process {
	return s.p
}

// Go .
func (s *Service) Go(f goprocess.ProcessFunc) goprocess.Process {
	return s.p.Go(f)
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
