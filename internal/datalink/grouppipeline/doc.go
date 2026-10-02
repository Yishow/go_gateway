// Package grouppipeline is the production owner of applied write groups: it
// reconciles the applied group revisions into runtime boundaries, fans typed
// samples out to them, drives their time-based closure, runs the durable
// delivery worker, and tells the legacy writer which outputs a group owns so
// there is exactly one writer per output.
package grouppipeline
