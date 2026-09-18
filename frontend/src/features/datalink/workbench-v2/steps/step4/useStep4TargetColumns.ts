import { useMemo } from 'react';
import { useStudioV2DatabaseMetadataQuery } from '@/hooks/datalink/useStudioV2WorkspaceDatabase';
import type { DbConnector } from '../../state/types';
import type { DbColumn } from '../../state/dbSchemas';
import { isMetadataForScope, metadataColumnsForScope } from './step4DatabaseHelpers';

export type Step4TargetMetadataStatus = 'not_checked' | 'checking' | 'exists' | 'missing' | 'forbidden' | 'failed';

export interface Step4TargetColumns {
  columns: DbColumn[];
  status: Step4TargetMetadataStatus;
  refetch: () => void;
}

/** 讀取已存目標資料表的實際欄位；回覆屬於其他範圍時視為尚未查詢，不沿用其欄位。 */
export function useStep4TargetColumns(connector: DbConnector): Step4TargetColumns {
  const saved = connector.persisted === true && connector.save_state === 'saved';
  const query = useStudioV2DatabaseMetadataQuery({
    connectorId: connector.connector_id,
    connectorRevision: connector.identity_revision,
    database: connector.database,
    schema: connector.schema,
    table: connector.table,
  }, saved);
  const columns = useMemo(() => metadataColumnsForScope(query.data, connector), [query.data, connector]);

  let status: Step4TargetMetadataStatus = 'not_checked';
  if (saved && query.isFetching) {
    status = 'checking';
  } else if (saved && query.isError) {
    status = 'failed';
  } else if (query.data && isMetadataForScope(query.data, connector)) {
    status = query.data.inspection_status;
  }

  const { refetch } = query;
  return { columns, status, refetch: () => { void refetch(); } };
}
