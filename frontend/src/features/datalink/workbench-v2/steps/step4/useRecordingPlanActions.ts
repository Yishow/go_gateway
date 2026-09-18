import { useEffect, useMemo, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useApplySchemaMutation, useCreateRecordingPlanMutation, usePreviewSchemaMutation, useTestWritePlanMutation } from '@/hooks/datalink/useStudioV2WorkspaceRecordingPlans';
import { getSafeErrorMessage, type SafeErrorMessage } from '@/utils/typedErrors';
import { getRecordingOperationId, isRecordingOperationUnconfirmed, parseRecordingTestWriteResult, RecordingPlanResponseError } from '@/utils/recordingPlanJson';
import type { RecordingPlan, SchemaOperation, SchemaPreviewToken, StreamMode, TestWriteResult } from '@/types/recordingPlan';
import type { DbConnector } from '../../state/types';
import type { RecordingPlanScope } from './recordingPlanMembership';

interface Options {
  activePlan?: RecordingPlan;
  connector?: DbConnector;
  scope?: RecordingPlanScope;
  disabled: boolean;
  managedSchemaSupported: boolean;
  testWritesSupported: boolean;
  refetch: () => unknown;
}
interface UnconfirmedOperationState { operationId?: string }

/**
 * Managed recording table prefix sent as `table_prefix`; the backend accepts and
 * validates it against defaultManagedTablePrefix
 * (internal/datalink/recordingplan/schema_preview.go — kept in sync both ways).
 * Locked by tests/unit/workbench-v2/managed-table-prefix.test.ts.
 */
export const MANAGED_TABLE_PREFIX = 'gw_record_';
const DEFAULT_PLAN_TIMEZONE = 'Asia/Taipei';

/** Keep action results bound to the selected saved recording scope. */
export function useRecordingPlanActions({ activePlan, connector, scope, disabled, managedSchemaSupported, testWritesSupported, refetch }: Options) {
  const { t } = useTranslation('workbench-v2');
  const createPlanMutation = useCreateRecordingPlanMutation();
  const previewSchemaMutation = usePreviewSchemaMutation();
  const applySchemaMutation = useApplySchemaMutation();
  const testWriteMutation = useTestWritePlanMutation();
  const [planName, setPlanName] = useState(t('step4.recording_default_name'));
  const [streamMode, setStreamMode] = useState<StreamMode>('window_summary');
  const [retentionDays, setRetentionDays] = useState(365);
  const [previewToken, setPreviewToken] = useState<SchemaPreviewToken | null>(null);
  const [applyOperation, setApplyOperation] = useState<SchemaOperation | null>(null);
  const [testResult, setTestResult] = useState<TestWriteResult | null>(null);
  const [actionError, setActionError] = useState<SafeErrorMessage | null>(null);
  const [applyUnconfirmed, setApplyUnconfirmed] = useState<UnconfirmedOperationState | null>(null);
  const [testUnconfirmed, setTestUnconfirmed] = useState<UnconfirmedOperationState | null>(null);
  const [busy, setBusy] = useState<'create' | 'preview' | 'apply' | 'test' | null>(null);
  const pending = useRef(false);
  const scopeKey = JSON.stringify([scope?.fingerprint, activePlan?.id, activePlan?.revision, connector?.connector_id,
    connector?.identity_revision, connector?.setup_revision, disabled]);
  // A new object identity per scope key lets late responses detect that the
  // selected plan, members or connector changed while they were in flight.
  const scopeVersion = useMemo(() => ({ key: scopeKey }), [scopeKey]);
  const currentScope = useRef<{ key: string } | null>(scopeVersion);
  const resultScope = useRef<{ key: string } | null>(null);
  currentScope.current = scopeVersion;
  useEffect(() => {
    currentScope.current = scopeVersion;
    setPreviewToken(null); setTestResult(null); setActionError(null);
    setApplyOperation(null); setApplyUnconfirmed(null); setTestUnconfirmed(null);
    return () => { currentScope.current = null; };
  }, [scopeVersion]);

  const begin = (kind: NonNullable<typeof busy>) => {
    if (disabled || pending.current || !scope?.members.length || !connector?.connector_id || !connector.identity_revision) return false;
    pending.current = true; setBusy(kind); setActionError(null);
    return true;
  };
  const finish = () => { pending.current = false; setBusy(null); };
  const setSafeActionError = (error: unknown) => {
    if (currentScope.current === scopeVersion) {
      resultScope.current = scopeVersion;
      setActionError(getSafeErrorMessage(error, (key, options) => t(key, options)));
    }
  };
  const markUnconfirmed = (kind: 'apply' | 'test', error: unknown) => {
    if (currentScope.current !== scopeVersion || !isRecordingOperationUnconfirmed(error)) return;
    const value = { operationId: getRecordingOperationId(error) };
    if (kind === 'apply') setApplyUnconfirmed(value); else setTestUnconfirmed(value);
  };

  const handleCreateDefaultPlan = async () => {
    const connectorId = connector?.connector_id;
    if (!scope || !connectorId || !planName.trim() || !begin('create')) return;
    try {
      await createPlanMutation.mutateAsync({
        workspace_id: scope.workspaceId, name: planName, timezone: DEFAULT_PLAN_TIMEZONE, members: scope.members,
        retention: { raw_days: 30, summary_days: retentionDays, events_days: 90, correction_horizon_hours: 24 },
        streams: scope.members.map(member => ({ stream_id: `stream-${member.measurement_id}`, measurement_id: member.measurement_id,
          equipment_id: member.equipment_id, mode: streamMode, raw_policy: 'every_sample' })),
        destinations: [{ destination_id: 'dest-default-db', connector_id: connectorId, connector_revision: connector?.identity_revision, table_prefix: MANAGED_TABLE_PREFIX }],
      });
      if (currentScope.current === scopeVersion) void refetch();
    } catch (error) { setSafeActionError(error); } finally { finish(); }
  };

  const handlePreviewSchema = async () => {
    const connectorId = connector?.connector_id;
    const workspaceRevision = connector?.setup_revision;
    if (!activePlan || !connector || !connectorId || !workspaceRevision || !begin('preview')) return;
    try {
      // The server applies the table prefix saved on the plan destination.
      const token = await previewSchemaMutation.mutateAsync({ plan_id: activePlan.id, connector_id: connectorId,
        expected_connector_revision: connector.identity_revision, expected_workspace_revision: workspaceRevision,
        expected_plan_revision: activePlan.revision, dialect: connector.kind });
      if (currentScope.current === scopeVersion) { resultScope.current = scopeVersion; setPreviewToken(token); }
    } catch (error) { setSafeActionError(error); } finally { finish(); }
  };

  // A retry confirms the same operation the preview issued, so the backend can
  // return the recorded result instead of running the batch again.
  const handleApplySchema = async () => {
    const workspaceRevision = previewToken?.workspace_revision;
    const connectorRevision = previewToken?.connector_revision ?? connector?.identity_revision;
    if (!previewToken || !workspaceRevision || !connectorRevision || !managedSchemaSupported || applyUnconfirmed || !begin('apply')) return;
    try {
      const operation = await applySchemaMutation.mutateAsync({
        token: previewToken.token,
        operation_id: previewToken.operation_id,
        expected_workspace_revision: workspaceRevision,
        expected_plan_revision: previewToken.plan_revision,
        expected_connector_revision: connectorRevision,
      });
      if (currentScope.current !== scopeVersion) return;
      resultScope.current = scopeVersion;
      setApplyOperation(operation);
      if (operation.status === 'succeeded') { setPreviewToken(null); void refetch(); }
    } catch (error) { markUnconfirmed('apply', error); setSafeActionError(error); } finally { finish(); }
  };

  const handleTestWrite = async () => {
    if (!activePlan || !testWritesSupported || testUnconfirmed || !begin('test')) return;
    setTestResult(null);
    try {
      const result = await testWriteMutation.mutateAsync({ plan_id: activePlan.id, table_prefix: MANAGED_TABLE_PREFIX });
      const parsed = parseRecordingTestWriteResult(result);
      if (!parsed) throw new RecordingPlanResponseError('recording test write', result);
      if (currentScope.current === scopeVersion) { resultScope.current = scopeVersion; setTestResult(parsed); }
    } catch (error) { markUnconfirmed('test', error); setSafeActionError(error); } finally { finish(); }
  };

  const currentResults = resultScope.current === scopeVersion;
  return { planName, setPlanName, streamMode, setStreamMode, retentionDays, setRetentionDays,
    previewToken: currentResults ? previewToken : null, testResult: currentResults ? testResult : null,
    applyOperation: currentResults ? applyOperation : null,
    actionError: currentResults ? actionError : null, applyUnconfirmed: currentResults ? applyUnconfirmed : null,
    testUnconfirmed: currentResults ? testUnconfirmed : null, busy,
    handleCreateDefaultPlan, handlePreviewSchema, handleApplySchema, handleTestWrite };
}
