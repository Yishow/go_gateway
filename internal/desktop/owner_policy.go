package desktop

import "errors"

// ErrHardlinkedDatabase rejects ambiguous pathname-based SQLite WAL recovery.
// Keep the original database and its journal together; never move/delete sidecars
// automatically to make an alias appear usable.
var ErrHardlinkedDatabase = errors.New("database has multiple hardlinks; preserve the original database and its journal, and do not move or delete sidecars")

func checkWritableLinkCount(count uint32) error {
	if count > 1 {
		return ErrHardlinkedDatabase
	}
	if count != 1 {
		return errors.New("database link count is unavailable")
	}
	return nil
}
