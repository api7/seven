package state

// Solver
func (s *ApisixCombination)Solver() (bool, error){
	// 1.route workers
	rwg := NewRouteWorkers(s.Routes)
	// 2.service workers
	swg := NewServiceWorkers(s.Services, &rwg)
	// 3.upstream workers
	SolverUpstream(s.Upstreams, swg)
	return true, nil
}
