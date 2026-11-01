package main

import (
	"slices"
	"time"
)

func acceptancePercentile(values []time.Duration, percentile int) time.Duration {
	ordered := slices.Clone(values)
	slices.Sort(ordered)
	return ordered[(len(ordered)*percentile+99)/100-1]
}
