//go:build f_write_group_fixture

package collector

// EmitFixtureValue injects one collected value into the real scheduler output
// channel for the tagged runtime consumer boundary test.
func (s *Scheduler) EmitFixtureValue(value CollectedValue) bool {
	if s == nil {
		return false
	}
	select {
	case s.valueChan <- value:
		return true
	default:
		return false
	}
}
