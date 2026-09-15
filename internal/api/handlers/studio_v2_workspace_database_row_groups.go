package handlers

import (
	"strings"

	"go-gateway/internal/datalink/workspace"
)

func workspaceDatabaseTargetRefsByPoint(refs []workspace.DatabaseTargetRef, groups []workspace.DatabaseRowGroup) map[string]string {
	out := make(map[string]string, len(refs))
	for _, ref := range refs {
		pointID := strings.TrimSpace(ref.PointID)
		rowGroupID := strings.TrimSpace(ref.RowGroupID)
		if workspaceDatabaseRowGroupContainsPoint(groups, rowGroupID, pointID) {
			out[pointID] = rowGroupID
		}
	}
	return out
}

func workspaceDatabaseRowGroupExists(groups []workspace.DatabaseRowGroup, id string) bool {
	id = strings.TrimSpace(id)
	if id == "" {
		return false
	}
	for _, group := range groups {
		if group.ID == id {
			return true
		}
	}
	return false
}

func workspaceDatabaseRowGroupContainsPoint(groups []workspace.DatabaseRowGroup, id, pointID string) bool {
	id = strings.TrimSpace(id)
	pointID = strings.TrimSpace(pointID)
	if id == "" || pointID == "" {
		return false
	}
	for _, group := range groups {
		if group.ID != id {
			continue
		}
		for _, memberPointID := range group.MemberPointIDs {
			if strings.TrimSpace(memberPointID) == pointID {
				return true
			}
		}
		return false
	}
	return false
}

func rowGroupTargetKey(rowGroupID, pointID string) string {
	rowGroupID = strings.TrimSpace(rowGroupID)
	pointID = strings.TrimSpace(pointID)
	if rowGroupID == "" || pointID == "" {
		return ""
	}
	return rowGroupID + ":" + pointID
}
