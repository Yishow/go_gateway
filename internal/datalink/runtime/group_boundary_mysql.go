package runtime

import (
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
)

// validateDestinationSample refuses a good value the frozen destination cannot
// represent before journaling or offering it. Bad or missing samples still
// contribute the shared NULL/status record.
func (b *GroupBoundary) validateDestinationSample(sample snapshot.Sample) error {
	if sample.Quality != schema.QualityGood || sample.Value.Type() == "" {
		return nil
	}
	layout, ok := b.layout.(*dbtarget.LongRowLayout)
	if !ok {
		return nil
	}
	if err := layout.ValidateSampleValue(sample.Value); err != nil {
		return &GroupSampleError{
			Outcome: snapshot.OfferRejected,
			Reason:  "exact_value_unrepresentable",
			Cause:   err,
		}
	}
	return nil
}
