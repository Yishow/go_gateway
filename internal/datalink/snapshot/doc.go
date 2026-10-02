// Package snapshot assembles periodic, deterministic write-group rows from
// typed acquisition samples. It is pure in-memory logic with injected time:
// no database, journal or sender. Durable acceptance and delivery belong to
// the delivery change; nothing here may be described as restart-safe.
//
// A basic snapshot is NOT an interval average, usage figure or every-sample
// history: each member contributes the sample with the greatest observed_at
// inside a UTC half-open bucket [start,end).
package snapshot
