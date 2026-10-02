// Package grouptestwrite performs explicit, operator-confirmed test writes for
// one saved WriteGroup. A preview binds the exact test row, the group and
// connector revisions and an operation in the shared durable operation ledger;
// a confirmation claims it once, writes one operation-owned row through the
// production row layout and insert path, reads it back and compares typed
// values, then removes only the rows that operation owns. Write verification
// and cleanup are reported independently, and an uncertain step is "unknown",
// never success.
package grouptestwrite
