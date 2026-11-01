import type * as React from 'react';
import { useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { useStudioV2MySQLRecordingCapabilityQuery, useStudioV2SQLServerRecordingCapabilityQuery } from '../../../../../../hooks/datalink/useStudioV2WorkspaceDatabase';
import type { StudioV2MySQLRecordingCapability } from '../../../../../../types/studioV2MySQLRecordingCapability';
import type { StudioV2SQLServerRecordingCapability } from '../../../../../../types/studioV2SQLServerRecordingCapability';
import type { DbConnector } from '../../../state/types';
import { mixedEditorAvailable } from '../../../recording-editor/RecordingMixedEditor';
import type { GroupEditorProps } from './GroupEditor';
export function destinationSaved(connector: DbConnector): boolean {
  return connector.persisted === true && connector.save_state === 'saved' &&
    Boolean(connector.connector_id && connector.identity_revision);
}
type CapabilityProps = GroupEditorProps & { Content: React.ComponentType<GroupEditorContentProps> };
/** Dialect capabilities remain canonical, independently queried only for saved destinations. */
export const GroupEditorWithCapabilities: React.FC<CapabilityProps> = ({ Content, ...props }) => {
  const { t } = useTranslation('workbench-v2');
  const unavailable = props.group?.recording_profile?.kind === 'mixed' && !mixedEditorAvailable(props.calculationCapabilities, props.connector.kind, props.group.layout_version);
  const admitted = useRef(!unavailable);
  if (!unavailable) admitted.current = true;
  // Mixed requires the explicit canonical consumer capability; absence never falls through to legacy conversion.
  if (!admitted.current) return (
    <section data-testid="mixed-group-editor-readonly" className="space-y-3 rounded border border-slate-700 p-4">
      <h3>{props.group?.name}</h3>
      <p role="status">{t('step4.group_kind_unsupported')}</p>
      {props.presentation !== 'simple' && <button type="button" onClick={props.onClose}>{t('step4.group.editor.back')}</button>}
    </section>
  );
  return <>{unavailable && <p role="status">{t('step4.group_kind_unsupported')}</p>}<DialectCapabilities {...props} Content={Content} editorReadonly={props.editorReadonly || unavailable} /></>;
};

/** Keep one Content instance while destination authority changes; disabled queries grant nothing. */
const DialectCapabilities: React.FC<CapabilityProps> = ({ Content, ...props }) => {
  const { connector } = props;
  const saved = destinationSaved(connector);
  const mysqlEnabled = connector.kind === 'mysql' && saved;
  const sqlserverEnabled = connector.kind === 'sqlserver' && saved;
  const mysql = useStudioV2MySQLRecordingCapabilityQuery(connector.connector_id ?? '', connector.identity_revision ?? '', mysqlEnabled, props.locatorScope);
  const sqlserver = useStudioV2SQLServerRecordingCapabilityQuery(connector.connector_id ?? '', connector.identity_revision ?? '', sqlserverEnabled, props.locatorScope);
  return <Content {...props}
    mysqlCapability={mysqlEnabled && !mysql.isError && !mysql.isFetching ? mysql.data : undefined}
    sqlserverCapability={sqlserverEnabled && !sqlserver.isError && !sqlserver.isFetching ? sqlserver.data : undefined} />;
};

export type GroupEditorContentProps = GroupEditorProps & {
  mysqlCapability?: StudioV2MySQLRecordingCapability;
  sqlserverCapability?: StudioV2SQLServerRecordingCapability;
};

