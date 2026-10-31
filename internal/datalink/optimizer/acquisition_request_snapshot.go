package optimizer

// Request returns one detached request in O(request units + point slices), not
// an O(plan size) copy for each wire read. Invalid indexes have no projection.
func (p AcquisitionPlan) Request(index int) (AcquisitionRequest, bool) {
	if index < 0 || index >= len(p.requests) {
		return AcquisitionRequest{}, false
	}
	request := p.requests[index]
	request.Points = append([]PointSlice(nil), request.Points...)
	request.Units = append([]NativeUnit(nil), request.Units...)
	return request, true
}
