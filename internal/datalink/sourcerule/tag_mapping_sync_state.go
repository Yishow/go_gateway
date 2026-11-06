package sourcerule

import "go-gateway/internal/datalink/schema"

type mappingRollbackState struct {
	enabled              bool
	transformPipeline    []schema.TransformStep
	status               schema.MappingStatus
	ruleCandidateID      string
	proposedSignature    string
	lastAppliedSignature string
	blockingReason       string
}

type tagRollbackState struct {
	dataType schema.DataType
}

type tagIdentityMove struct{ before, after *schema.Tag }
type mappingIdentityMove struct{ before, after *schema.Mapping }

type tagMappingSyncResult struct {
	createdTagIDs     []string
	createdMappingIDs []string
	updatedMappings   map[string]mappingRollbackState
	updatedTags       map[string]tagRollbackState
	tagIdentityMoves  []tagIdentityMove
	mapIdentityMoves  []mappingIdentityMove
}
