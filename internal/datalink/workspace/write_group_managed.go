package workspace

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

func managedGroupTableName(id string) string {
	digest := sha256.Sum256([]byte(id))
	return "gw_group_" + hex.EncodeToString(digest[:16])
}

func managedGroupDeviceID(group *WriteGroup) string {
	if len(group.Members) == 0 {
		return ""
	}
	device := group.Members[0].DeviceID
	for _, member := range group.Members {
		if member.DeviceID != device {
			return ""
		}
	}
	return device
}

// prepareManagedWriteGroup fixes the mapping at save time. Display labels and
// member ordering never participate in column identity.
func prepareManagedWriteGroup(group *WriteGroup) error {
	if group.Destination.StorageStrategy != WriteGroupStorageStrategyManaged {
		return nil
	}
	for i := range group.Members {
		member := &group.Members[i]
		identity, err := json.Marshal([]string{member.DeviceID, member.PointID, member.TagID})
		if err != nil {
			return fmt.Errorf("encode managed member identity: %w", err)
		}
		digest := sha256.Sum256(identity)
		member.TargetColumn = "v_" + hex.EncodeToString(digest[:12])
	}
	group.RowPolicy.RecordKeyColumn = "record_id"
	group.RowPolicy.GroupIDColumn = "group_id"
	group.RowPolicy.BucketStartColumn = "bucket_start"
	group.RowPolicy.ProvenanceColumn = "provenance"
	group.RowPolicy.DeviceIDColumn = ""
	if len(group.Members) > 0 {
		device := group.Members[0].DeviceID
		singleDevice := true
		for _, member := range group.Members {
			singleDevice = singleDevice && member.DeviceID == device
		}
		if singleDevice {
			group.RowPolicy.DeviceIDColumn = "device_id"
		}
	}
	group.WritePolicy.Mode = "append"
	group.WritePolicy.DedupeCapability = "receipt"
	return nil
}
