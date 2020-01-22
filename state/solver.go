package state

// Solver
func (s *ApisixCombination)Solver() (bool, error){
	quit := make(chan Quit)
	// 1.route workers
	rwg := NewRouteWorkers(s.Routes, quit)
	// 2.service workers
	swg := NewServiceWorkers(s.Services, quit, &rwg)
	// 3.upstream workers
	SolverUpstream(s.Upstreams, swg)
	return true, nil
}
