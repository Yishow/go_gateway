package handlers

import (
	"slices"
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
	return slices.ContainsFunc(groups, func(group workspace.DatabaseRowGroup) bool {
		return group.ID == id
	})
}

func workspaceDatabaseRowGroupContainsPoint(groups []workspace.DatabaseRowGroup, id, pointID string) bool {
	id = strings.TrimSpace(id)
	pointID = strings.TrimSpace(pointID)
	if id == "" || pointID == "" {
		return false
	}
	for _, group := range groups {
		if group.ID == id {
			return slices.ContainsFunc(group.MemberPointIDs, func(memberPointID string) bool {
				return strings.TrimSpace(memberPointID) == pointID
			})
		}
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
