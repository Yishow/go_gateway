import type * as React from 'react';
import { useStudioV2MySQLRecordingCapabilityQuery, useStudioV2SQLServerRecordingCapabilityQuery } from '../../../../../hooks/datalink/useStudioV2WorkspaceDatabase';
import type { StudioV2MySQLRecordingCapability } from '../../../../../types/studioV2MySQLRecordingCapability';
import type { StudioV2SQLServerRecordingCapability } from '../../../../../types/studioV2SQLServerRecordingCapability';
import { destinationIsSaved } from './basicRecordingPanelHelpers';
import type { CircuitProfileBasicPanelProps } from './CircuitProfileBasicPanel';
interface MySQLCapabilityQueryState {
  data?: StudioV2MySQLRecordingCapability;
  isLoading: boolean;
  isFetching: boolean;
  isError: boolean;
  refetch: () => Promise<unknown>;
}

const EMPTY_MYSQL_CAPABILITY_QUERY: MySQLCapabilityQueryState = {
  data: undefined, isLoading: false, isFetching: false, isError: false, refetch: async () => undefined,
};

interface SQLServerCapabilityQueryState {
  data?: StudioV2SQLServerRecordingCapability;
  isLoading: boolean;
  isFetching: boolean;
  isError: boolean;
  refetch: () => Promise<unknown>;
}

const EMPTY_SQLSERVER_CAPABILITY_QUERY: SQLServerCapabilityQueryState = {
  data: undefined, isLoading: false, isFetching: false, isError: false, refetch: async () => undefined,
};


export type CircuitProfileBasicPanelContentProps = CircuitProfileBasicPanelProps & {
  mysqlCapabilityQuery: MySQLCapabilityQueryState;
  sqlserverCapabilityQuery: SQLServerCapabilityQueryState;
};
/** Same editor owner across destination authority changes; disabled queries grant no capability. */
export function CircuitProfileCapabilities({ Content, ...props }: CircuitProfileBasicPanelProps & { Content: React.ComponentType<CircuitProfileBasicPanelContentProps> }) {
  const connector = props.state.db.connector;
  const mysqlEnabled = connector.kind === 'mysql' && destinationIsSaved(props.state);
  const sqlserverEnabled = connector.kind === 'sqlserver' && destinationIsSaved(props.state);
  const mysql = useStudioV2MySQLRecordingCapabilityQuery(connector.connector_id ?? '', connector.identity_revision ?? '', mysqlEnabled);
  const sqlserver = useStudioV2SQLServerRecordingCapabilityQuery(connector.connector_id ?? '', connector.identity_revision ?? '', sqlserverEnabled);
  return <Content {...props} mysqlCapabilityQuery={mysqlEnabled ? mysql : EMPTY_MYSQL_CAPABILITY_QUERY}
    sqlserverCapabilityQuery={sqlserverEnabled ? sqlserver : EMPTY_SQLSERVER_CAPABILITY_QUERY} />;
}
