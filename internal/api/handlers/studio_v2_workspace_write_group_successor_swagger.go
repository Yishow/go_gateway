package handlers

import (
	"time"

	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"
)

// These documentation-only request models describe the existing strict decoder;
// they are not used to decode, authorize, or execute mutations. Keeping each
// action separate avoids suggesting that fields from different stages can mix.

// WriteGroupSuccessorAction is only a discriminator, NOT a complete request.
// @Description Swagger 2 cannot express this action-dependent union. Select exactly one full request from the operation's x-action-contracts; fields from other actions are rejected. Sending only action is invalid.
type WriteGroupSuccessorAction struct {
	Action string `json:"action" binding:"required" enums:"reserve,schema-preview,schema-confirm,activation-review,activation-readiness,activation-confirm,first-start-confirm,operation-status"`
}

// WriteGroupSuccessorReviewRequest documents the local advisory review input.
// @Description All scope and original CAS fields are required. Path id is G1. operation_id is the original root request ID, never a schema operation or first-start result ID. group is a candidate, not authority; boundary is an explicit timestamp, not permission to activate.
type WriteGroupSuccessorReviewRequest struct {
	WorkspaceID                          string                `json:"workspace_id" binding:"required"`
	ExpectedWorkspaceRevision            string                `json:"expected_workspace_revision" binding:"required"`
	LineageID                            string                `json:"lineage_id" binding:"required"`
	ExpectedLineageRevision              int64                 `json:"expected_lineage_revision" binding:"required" minimum:"1"`
	ExpectedPredecessorRevision          string                `json:"expected_predecessor_revision" binding:"required"`
	ExpectedPredecessorConnectorRevision string                `json:"expected_predecessor_connector_revision" binding:"required"`
	ExpectedTargetConnectorRevision      string                `json:"expected_target_connector_revision" binding:"required"`
	OperationID                          string                `json:"operation_id" binding:"required"`
	Boundary                             time.Time             `json:"boundary" binding:"required"`
	Group                                *workspace.WriteGroup `json:"group" binding:"required"`
}

// WriteGroupSuccessorCurrentCAS documents separately refreshed stage revisions.
// @Description Current CAS is refreshed for a new stage; it never replaces the original immutable envelope or root operation ID. group_id/group_revision name reserved G2, predecessor_revision names path G1. An immutable reserve result is not current CAS.
type WriteGroupSuccessorCurrentCAS struct {
	WorkspaceRevision            string `json:"workspace_revision" binding:"required"`
	LineageRevision              int64  `json:"lineage_revision" binding:"required" minimum:"1"`
	PredecessorRevision          string `json:"predecessor_revision" binding:"required"`
	PredecessorConnectorRevision string `json:"predecessor_connector_revision" binding:"required"`
	GroupID                      string `json:"group_id" binding:"required"`
	GroupRevision                string `json:"group_revision" binding:"required"`
	TargetConnectorRevision      string `json:"target_connector_revision" binding:"required"`
}

// WriteGroupSuccessorStageRequest binds current CAS to the original review.
// @Description Return the complete original server envelope and review_digest unchanged, plus separately refreshed current CAS. Path id must match envelope.predecessor_group_id; workspace and lineage remain the original scope.
type WriteGroupSuccessorStageRequest struct {
	Envelope     workspace.WriteGroupSuccessorEnvelope `json:"envelope" binding:"required"`
	ReviewDigest string                                `json:"review_digest" binding:"required"`
	Current      WriteGroupSuccessorCurrentCAS         `json:"current" binding:"required"`
}

// WriteGroupSuccessorReserveRequest records the exact original reviewed intent.
// @Description reserve revalidates original envelope scope/CAS and may probe required source/target capabilities. Exact replay returns the immutable original reservation even after head/revisions advance; changed intent under the same root conflicts. It does not run DDL or activate intake.
type WriteGroupSuccessorReserveRequest struct {
	Action       string                                `json:"action" binding:"required" enums:"reserve"`
	Envelope     workspace.WriteGroupSuccessorEnvelope `json:"envelope" binding:"required"`
	ReviewDigest string                                `json:"review_digest" binding:"required"`
}

// WriteGroupSuccessorPreviewRequest selects a side-effect-free target inspection.
// @Description schema-preview may perform external read-only target probes and stores a preview token locally; it never executes DDL or activates intake. The returned operation_id is a schema operation ID, not the original root ID.
type WriteGroupSuccessorPreviewRequest struct {
	Action string `json:"action" binding:"required" enums:"schema-preview"`
	WriteGroupSuccessorStageRequest
}

// WriteGroupSuccessorConfirmRequest explicitly approves the stored schema plan.
// @Description schema-confirm requires the preview's token, 64-character digest and schema_operation_id, plus original envelope and current CAS. Only the server-stored plan can execute DDL. Exact confirmation replay keeps its recorded operation. HTTP 200 is a recorded stage result, not necessarily succeeded; HTTP 202 means in progress. Never substitute schema_operation_id for the root operation ID.
type WriteGroupSuccessorConfirmRequest struct {
	Action string `json:"action" binding:"required" enums:"schema-confirm"`
	WriteGroupSuccessorStageRequest
	Token             string `json:"token" binding:"required"`
	Digest            string `json:"digest" binding:"required" minLength:"64" maxLength:"64"`
	SchemaOperationID string `json:"schema_operation_id" binding:"required"`
}

// WriteGroupSuccessorActivationReviewRequest selects an explicit future boundary.
// @Description activation-review is local read only after schema preparation, using original stage envelope plus current CAS and an explicit future boundary B. The returned activation_kind selects cutover versus first_start; it grants no activation permission.
type WriteGroupSuccessorActivationReviewRequest struct {
	Action   string                          `json:"action" binding:"required" enums:"activation-review"`
	Stage    WriteGroupSuccessorStageRequest `json:"stage" binding:"required"`
	Boundary time.Time                       `json:"boundary" binding:"required"`
}

// WriteGroupSuccessorReadinessRequest selects inspection, not confirmation.
// @Description activation-readiness requires the exact activation-review envelope/digest, including its stage current CAS and future B. It may inspect the target but performs no DDL or activation; readiness_digest is evidence, not a persisted grant.
type WriteGroupSuccessorReadinessRequest struct {
	Action       string                                          `json:"action" binding:"required" enums:"activation-readiness"`
	Envelope     workspace.WriteGroupSuccessorActivationEnvelope `json:"envelope" binding:"required"`
	ReviewDigest string                                          `json:"review_digest" binding:"required" minLength:"64" maxLength:"64"`
}

// WriteGroupSuccessorCutoverConfirmRequest confirms a reviewed cutover only.
// @Description activation-confirm requires the unchanged activation envelope, review_digest and readiness_digest; it rechecks future B, scope/CAS, schema and live readiness. activation_kind must be cutover; only a valid legacy v1 envelope may omit it. first_start envelopes are rejected. Success preserves the existing WriteGroupCutover response, not proof of runtime or target writes.
type WriteGroupSuccessorCutoverConfirmRequest struct {
	Action          string                                          `json:"action" binding:"required" enums:"activation-confirm"`
	Envelope        workspace.WriteGroupSuccessorActivationEnvelope `json:"envelope" binding:"required"`
	ReviewDigest    string                                          `json:"review_digest" binding:"required" minLength:"64" maxLength:"64"`
	ReadinessDigest string                                          `json:"readiness_digest" binding:"required" minLength:"64" maxLength:"64"`
}

// WriteGroupSuccessorFirstStartConfirmRequest confirms a never-applied lineage.
// @Description first-start-confirm requires activation_kind=first_start in the unchanged activation envelope and both digests. It rechecks whole-lineage absence of applied/intake/history evidence, current CAS, schema/readiness and explicit future B. It returns WriteGroupSuccessorFirstStartResult, never an invented predecessor end/cutover/fence. Cutover envelopes are rejected.
type WriteGroupSuccessorFirstStartConfirmRequest struct {
	Action          string                                          `json:"action" binding:"required" enums:"first-start-confirm"`
	Envelope        workspace.WriteGroupSuccessorActivationEnvelope `json:"envelope" binding:"required"`
	ReviewDigest    string                                          `json:"review_digest" binding:"required" minLength:"64" maxLength:"64"`
	ReadinessDigest string                                          `json:"readiness_digest" binding:"required" minLength:"64" maxLength:"64"`
}

// WriteGroupSuccessorStatusRequest identifies the immutable original operation.
// @Description operation-status is a coherent local read with exact workspace_id, lineage_id and original request operation_id; path id is original G1. No current revision, envelope, target probe or token is required or accepted. Schema/first-start result IDs, current head and latest preview cannot replace the root. Response contains all related schema outcomes, not just the latest.
type WriteGroupSuccessorStatusRequest struct {
	Action      string `json:"action" binding:"required" enums:"operation-status"`
	WorkspaceID string `json:"workspace_id" binding:"required"`
	LineageID   string `json:"lineage_id" binding:"required"`
	OperationID string `json:"operation_id" binding:"required"`
}

// WriteGroupSuccessorResponse documents the common wrapper, not a data union.
// @Description data is exactly the selected action's response_data schema in x-action-contracts, not an arbitrary combination. A recorded stage success is not runtime running and not SQL committed.
type WriteGroupSuccessorResponse struct {
	Success bool           `json:"success" example:"true"`
	Data    map[string]any `json:"data"`
}

// WriteGroupSuccessorSchemaResponse is the schema-confirm in-progress wrapper.
type WriteGroupSuccessorSchemaResponse struct {
	Success bool                           `json:"success" example:"true"`
	Data    *recordingplan.SchemaOperation `json:"data"`
}

// writeGroupSuccessorSwaggerModels is a route-less model catalog for Swag 2.
// Swag parses these references into definitions but no operation without Router.
// This is not an endpoint, request body, new stage URL, or runtime authority.
// The real dispatcher links these exact models through x-action-contracts.
// @Param reserve body WriteGroupSuccessorReserveRequest true "reserve only"
// @Param preview body WriteGroupSuccessorPreviewRequest true "schema-preview only"
// @Param confirm body WriteGroupSuccessorConfirmRequest true "schema-confirm only"
// @Param activationReview body WriteGroupSuccessorActivationReviewRequest true "activation-review only"
// @Param readiness body WriteGroupSuccessorReadinessRequest true "activation-readiness only"
// @Param cutover body WriteGroupSuccessorCutoverConfirmRequest true "activation-confirm only"
// @Param firstStart body WriteGroupSuccessorFirstStartConfirmRequest true "first-start-confirm only"
// @Param status body WriteGroupSuccessorStatusRequest true "operation-status only"
// @Success 200 {object} workspace.WriteGroupSuccessorPreparation
// @Success 200 {object} recordingplan.SchemaPreviewToken
// @Success 200 {object} recordingplan.SchemaOperation
// @Success 200 {object} workspace.WriteGroupSuccessorActivationReview
// @Success 200 {object} workspace.WriteGroupSuccessorActivationReadiness
// @Success 200 {object} workspace.WriteGroupCutover
// @Success 200 {object} workspace.WriteGroupSuccessorFirstStartResult
// @Success 200 {object} workspace.WriteGroupSuccessorOperationStatus
func writeGroupSuccessorSwaggerModels() {}
