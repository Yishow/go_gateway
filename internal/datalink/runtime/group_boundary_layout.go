package runtime

import (
	"cmp"
	"strings"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/snapshot"
)

func boundaryRowSpec(cfg GroupBoundaryConfig, members []dbtarget.GroupRowMember, entityKeyed bool) dbtarget.GroupRowSpec {
	group := cfg.Group
	deviceID := ""
	if len(group.Members) > 0 {
		deviceID = group.Members[0].DeviceID
		for _, member := range group.Members {
			if member.DeviceID != deviceID {
				deviceID = ""
				break
			}
		}
	}
	return dbtarget.GroupRowSpec{
		Dialect: cfg.Dialect, Columns: cfg.Columns, Members: members,
		Partial:     snapshot.IncompletePolicy(strings.ToLower(strings.TrimSpace(group.RowPolicy.IncompletePolicy))) == snapshot.IncompletePartial,
		EntityKeyed: entityKeyed, EntityKeyColumn: group.RowPolicy.EntityKeyColumn,
		RecordKeyColumn:   cmp.Or(group.RowPolicy.RecordKeyColumn, cfg.RecordKeyColumn),
		BucketStartColumn: cmp.Or(group.RowPolicy.BucketStartColumn, cfg.BucketStartColumn),
		ProvenanceColumn:  group.RowPolicy.ProvenanceColumn,
		GroupIDColumn:     group.RowPolicy.GroupIDColumn, GroupID: group.ID,
		DeviceIDColumn: group.RowPolicy.DeviceIDColumn, DeviceID: deviceID,
	}
}
