package collector

// OutputLossDiagnostics returns cumulative delivery loss with bounded recent identity detail.
func (s *Scheduler) OutputLossDiagnostics() OutputLossSnapshot { return s.outputLoss.snapshot() }
