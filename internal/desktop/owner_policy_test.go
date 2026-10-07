package desktop

import "testing"

func TestWritableDatabaseRequiresSingleLinkForWALRecovery(t *testing.T) {
	for _, count := range []uint32{0, 1, 2, 10} {
		err := checkWritableLinkCount(count)
		if (err == nil) != (count == 1) {
			t.Errorf("link count %d: %v", count, err)
		}
	}
}
