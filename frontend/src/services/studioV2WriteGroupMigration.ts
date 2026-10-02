import type {
  RecordingPlanMigrationPreview,
  WriteGroupMigrationPreview,
  WriteGroupMigrationPreviewRequest,
  WriteGroupMigrationReviewRequest,
  WriteGroupMigrationReviewResponse,
} from '../types/studioV2WriteGroupMigration';
import {
  parseWriteGroupMigrationAPIData,
  parseWriteGroupMigrationPreviewData,
  parseWriteGroupMigrationReviewData,
} from '../utils/studioV2WriteGroupMigrationJson';
import {
  parseRecordingPlanMigrationPreviewData,
  parseRecordingPlanMigrationReviewData,
} from '../utils/studioV2WriteGroupRecordingPlanMigrationJson';
import { studioV2DatalinkApi } from './studioV2Workspace';

export const studioV2WriteGroupMigrationAPI = {
  async previewRecordingPlanMigration(
    request: WriteGroupMigrationPreviewRequest,
  ): Promise<RecordingPlanMigrationPreview> {
    const payload: WriteGroupMigrationPreviewRequest = {
      workspace_id: request.workspace_id,
      source_ids: [...request.source_ids],
    };
    const res = await studioV2DatalinkApi.post<unknown>(
      '/studio-v2/workspace/write-groups/migrations/recording-plans/preview',
      payload,
    );
    return parseWriteGroupMigrationAPIData(
      res.data,
      'write-group recording-plan migration preview',
      parseRecordingPlanMigrationPreviewData,
    );
  },

  async reviewRecordingPlanMigration(
    request: WriteGroupMigrationReviewRequest,
  ): Promise<WriteGroupMigrationReviewResponse> {
    const payload: WriteGroupMigrationReviewRequest = {
      workspace_id: request.workspace_id,
      expected_workspace_revision: request.expected_workspace_revision,
      expected_connector_revision: request.expected_connector_revision,
      review_digest: request.review_digest,
      source_ids: [...request.source_ids],
      confirm_snapshot_conversion: request.confirm_snapshot_conversion,
    };
    const res = await studioV2DatalinkApi.post<unknown>(
      '/studio-v2/workspace/write-groups/migrations/recording-plans/review',
      payload,
    );
    return parseWriteGroupMigrationAPIData(
      res.data,
      'write-group recording-plan migration review',
      parseRecordingPlanMigrationReviewData,
    );
  },

  async previewSingleMappingMigration(
    request: WriteGroupMigrationPreviewRequest,
  ): Promise<WriteGroupMigrationPreview> {
    const payload: WriteGroupMigrationPreviewRequest = {
      workspace_id: request.workspace_id,
      source_ids: [...request.source_ids],
    };
    const res = await studioV2DatalinkApi.post<unknown>(
      '/studio-v2/workspace/write-groups/migrations/single-mappings/preview',
      payload,
    );
    return parseWriteGroupMigrationAPIData(
      res.data,
      'write-group migration preview',
      parseWriteGroupMigrationPreviewData,
    );
  },

  async reviewSingleMappingMigration(
    request: WriteGroupMigrationReviewRequest,
  ): Promise<WriteGroupMigrationReviewResponse> {
    const payload: WriteGroupMigrationReviewRequest = {
      workspace_id: request.workspace_id,
      expected_workspace_revision: request.expected_workspace_revision,
      expected_connector_revision: request.expected_connector_revision,
      review_digest: request.review_digest,
      source_ids: [...request.source_ids],
      confirm_snapshot_conversion: request.confirm_snapshot_conversion,
    };
    const res = await studioV2DatalinkApi.post<unknown>(
      '/studio-v2/workspace/write-groups/migrations/single-mappings/review',
      payload,
    );
    return parseWriteGroupMigrationAPIData(
      res.data,
      'write-group migration review',
      parseWriteGroupMigrationReviewData,
    );
  },

  async previewRowGroupMigration(
    request: WriteGroupMigrationPreviewRequest,
  ): Promise<WriteGroupMigrationPreview> {
    const payload: WriteGroupMigrationPreviewRequest = {
      workspace_id: request.workspace_id,
      source_ids: [...request.source_ids],
    };
    const res = await studioV2DatalinkApi.post<unknown>(
      '/studio-v2/workspace/write-groups/migrations/row-groups/preview',
      payload,
    );
    return parseWriteGroupMigrationAPIData(
      res.data,
      'write-group row-group migration preview',
      parseWriteGroupMigrationPreviewData,
    );
  },

  async reviewRowGroupMigration(
    request: WriteGroupMigrationReviewRequest,
  ): Promise<WriteGroupMigrationReviewResponse> {
    const payload: WriteGroupMigrationReviewRequest = {
      workspace_id: request.workspace_id,
      expected_workspace_revision: request.expected_workspace_revision,
      expected_connector_revision: request.expected_connector_revision,
      review_digest: request.review_digest,
      source_ids: [...request.source_ids],
      confirm_snapshot_conversion: request.confirm_snapshot_conversion,
    };
    const res = await studioV2DatalinkApi.post<unknown>(
      '/studio-v2/workspace/write-groups/migrations/row-groups/review',
      payload,
    );
    return parseWriteGroupMigrationAPIData(
      res.data,
      'write-group row-group migration review',
      parseWriteGroupMigrationReviewData,
    );
  },
};
