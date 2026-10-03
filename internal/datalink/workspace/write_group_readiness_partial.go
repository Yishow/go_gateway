package workspace

import (
	"strings"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
)

// Partial readiness uses the same inspected row layout as the production
// boundary. A saved opt-in alone never proves nullable or provenance storage.
func evaluateWriteGroupPartialInspection(group *WriteGroup, tagTypes map[string]schema.DataType, inspection *dbtarget.TableInspection, connectorKind string) []ReadinessIssue {
	if strings.ToLower(strings.TrimSpace(group.RowPolicy.IncompletePolicy)) != "partial" || inspection == nil || inspection.Status != dbtarget.TableInspectionExists {
		return nil
	}
	dialect := dbtarget.SQLDialectSQLite
	if connectorKind == string(schema.DatabaseConnectorKindPostgres) {
		dialect = dbtarget.SQLDialectPostgres
	}
	members := make([]dbtarget.GroupRowMember, 0, len(group.Members))
	entityKeyed := false
	for _, member := range group.Members {
		entityKeyed = entityKeyed || member.EntityKey != ""
		kind, _ := measurement.ExactTypeForTag(tagTypes[member.TagID])
		members = append(members, dbtarget.GroupRowMember{
			MemberKey: member.TagID, EntityKey: member.EntityKey, Column: member.TargetColumn,
			Type: kind, Required: member.Required,
		})
	}
	_, layoutIssues := dbtarget.NewGroupRowLayout(dbtarget.GroupRowSpec{
		Dialect: dialect, Columns: inspection.Columns, Members: members, Partial: true,
		ProvenanceColumn: group.RowPolicy.ProvenanceColumn, EntityKeyed: entityKeyed,
		EntityKeyColumn: group.RowPolicy.EntityKeyColumn,
	})
	issues := make([]ReadinessIssue, 0, len(layoutIssues))
	for _, issue := range layoutIssues {
		issues = append(issues, writeGroupReadinessIssue(
			"partial-"+issue.Code, "partial rows require verified nullable values and member provenance", group.ID,
		))
	}
	return issues
}
