// Package groupdelivery makes write-group rows durable: a sample is ACKed only
// after its journal row commits, closed buckets (row, outbox intent, bucket
// outcome and checkpoint) commit in one local transaction, and senders confirm
// destination effects by a destination-scoped effect key.
//
// It deliberately does not promise exactly-once delivery: when a destination
// cannot prove an effect, an uncertain outcome is blocked for reconciliation
// instead of being blindly retried.
package groupdelivery
