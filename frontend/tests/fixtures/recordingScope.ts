export const savedRecordingScope = {
  workspaceId: 'workspace-a', fingerprint: 'saved-membership-1',
  members: [{ member_id: 'member-1', measurement_id: 'measurement-1', equipment_id: 'equipment-1', name: 'Fixture measurement' }],
};

export const savedRecordingPlanScopeFields = {
  workspace_id: savedRecordingScope.workspaceId,
  members: savedRecordingScope.members,
  streams: [{ stream_id: 'stream-1', measurement_id: 'measurement-1', mode: 'raw_history' }],
  destinations: [{ connector_id: 'saved-sqlite', connector_revision: 'identity-1' }],
};
