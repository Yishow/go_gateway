package pollinggroup

import "database/sql"

// UsesConfigDB reports whether the group store participates in the point store's database.
func (r *SQLRepository) UsesConfigDB(db *sql.DB) bool { return r != nil && r.db == db }

// IsInMemoryPollingGroups identifies the memory-only polling group store.
func (r *MemoryRepository) IsInMemoryPollingGroups() bool { return r != nil }
