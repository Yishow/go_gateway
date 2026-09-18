import { useTranslation } from 'react-i18next';
import { useStudioV2WorkspaceMeasurementsQuery } from '@/hooks/datalink/useStudioV2WorkspaceMeasurements';
import { SafeQueryBoundary } from '@/utils/SafeQueryBoundary';
import type { WorkbenchV2State } from '../../state/types';
import { RecordingPlanSetupSection } from './RecordingPlanSetupSection';
import { resolveRecordingPlanScope } from './recordingPlanMembership';

interface Props {
  state: WorkbenchV2State;
  workspaceId?: string;
  disabled?: boolean;
}

/** Resolve recording membership from saved definitions across the workspace. */
function WorkspaceRecordingPlanContent({ state, workspaceId, disabled }: Props) {
  const { t } = useTranslation('workbench-v2');
  const id = workspaceId ?? state.db.connector.workspace_id ?? '';
  const measurements = useStudioV2WorkspaceMeasurementsQuery(Boolean(id));

  // A failed or pending membership read must not look like an empty workspace
  // that is ready for a new plan.
  if (id && measurements.isError) {
    return (
      <div role="alert" className="text-xs text-amber-300" data-testid="recording-members-load-failed">
        {t('step4.recording_members_load_failed')}
        <button type="button" className="ml-2 underline" onClick={() => void measurements.refetch()}>
          {t('step4.recording_retry')}
        </button>
      </div>
    );
  }
  if (id && !measurements.isSuccess) {
    return (
      <p role="status" className="text-xs text-slate-400" data-testid="recording-members-loading">
        {t('step4.recording_members_loading')}
      </p>
    );
  }

  const scope = id ? resolveRecordingPlanScope(state, id, measurements.data ?? []) : undefined;
  return <RecordingPlanSetupSection connector={state.db.connector} scope={scope} disabled={disabled || !id} />;
}

export function WorkspaceRecordingPlanSetupSection(props: Props) {
  return <SafeQueryBoundary><WorkspaceRecordingPlanContent {...props} /></SafeQueryBoundary>;
}
