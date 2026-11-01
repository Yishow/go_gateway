package optimizer

import (
	"cmp"
	"slices"
)

// AcquisitionWireUnits counts unique physical address cells for one complete
// compatibility namespace. It excludes bridge gaps and rejected scalar spans.
// AddressUnitBits is the physical cell width, not the read component width.
type AcquisitionWireUnits struct {
	Group              AcquisitionGroup
	AddressUnitBits    int
	UniqueAddressUnits int
}

// WireUnits returns detached physical-address union summaries. Canonical areas
// remain separate; the planner does not infer aliases between protocol areas.
func (p AcquisitionPlan) WireUnits() []AcquisitionWireUnits {
	return append([]AcquisitionWireUnits(nil), p.wireUnits...)
}

// The key retains declarations as well as their revisions. Reusing a revision
// with contradictory capability or coverage must not silently deduplicate cells.
type acquisitionUnitUnion struct {
	summary AcquisitionWireUnits
	spans   []NativeRange
}

func collectAcquisitionWireUnits(groups map[string]*acquisitionUnitUnion, source AcquisitionSource, scalar AcquisitionRequest) {
	key := canonicalJSON(struct {
		Group      AcquisitionGroup
		Capability AcquisitionCapability
		Coverage   ReadableCoverage
	}{scalar.Group, source.Capability, source.Coverage})
	group := groups[key]
	if group == nil {
		group = &acquisitionUnitUnion{summary: AcquisitionWireUnits{
			Group:           scalar.Group,
			AddressUnitBits: scalar.Points[0].NativeBits / scalar.AddressStep,
		}}
		groups[key] = group
	}
	// planScalar already validated the complete span and physical address step.
	// Component phase constrains packing but not overlapping physical footprints.
	group.spans = append(group.spans, NativeRange{Area: scalar.Group.Area, Start: scalar.Start, End: scalar.Start + scalar.Length*scalar.AddressStep})
}

func summarizeAcquisitionWireUnits(groups map[string]*acquisitionUnitUnion) []AcquisitionWireUnits {
	keys := make([]string, 0, len(groups))
	for key := range groups {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	result := make([]AcquisitionWireUnits, 0, len(keys))
	for _, key := range keys {
		group := groups[key]
		slices.SortFunc(group.spans, func(a, b NativeRange) int {
			if n := cmp.Compare(a.Start, b.Start); n != 0 {
				return n
			}
			return cmp.Compare(a.End, b.End)
		})
		start, end := group.spans[0].Start, group.spans[0].End
		for _, span := range group.spans[1:] {
			if span.Start > end {
				group.summary.UniqueAddressUnits += end - start
				start, end = span.Start, span.End
			} else if span.End > end {
				end = span.End
			}
		}
		group.summary.UniqueAddressUnits += end - start
		result = append(result, group.summary)
	}
	return result
}
