import { useTranslation } from 'react-i18next';
import type { WriteGroupTestWriteOperation } from '../../../../../types/studioV2WriteGroupTestWrite';

interface Props {
  operation: WriteGroupTestWriteOperation;
  testIdPrefix?: string;
  manualGuidance?: boolean;
}
/** Display only: a committed but unverified write is not an unknown write or successful readback. */
export function RecordingTestWriteOperationResult({ operation, testIdPrefix = 'group-test-write', manualGuidance = false }: Props) {
  const { t } = useTranslation('workbench-v2');
  const key = (suffix: string) => t(`step4.group.test_write.${suffix}`);
  const active = operation.status === 'pending' || operation.status === 'running';
  const verified = operation.write_outcome === 'written_verified';
  const committed = verified || operation.write_outcome === 'written_unverified';
  const notWritten = operation.write_outcome === 'failed';
  const cleanupClaim = operation.status === 'running' && verified &&
    (operation.cleanup_status === 'failed' || operation.cleanup_status === 'unknown') && Boolean(operation.completed_at);
  const write = committed ? 'write_committed' : notWritten ? 'write_not_written' : 'write_unconfirmed';
  const readback = verified ? 'readback_verified' : committed ? 'readback_unverified' : notWritten ? 'readback_not_attempted' : 'readback_unconfirmed';
  return <div data-testid={`${testIdPrefix}-result`}>
    <dl className="grid grid-cols-[auto_1fr] gap-x-3 gap-y-1 text-xs text-slate-300">
      <dt>{key('operation')}</dt><dd className="font-mono">{operation.operation_id}</dd>
      {active && <><dt>{key('status')}</dt><dd data-testid={`${testIdPrefix}-running`}>{key('running')}</dd></>}
      {operation.write_outcome && <><dt>{key('recorded_outcome')}</dt><dd data-testid={`${testIdPrefix}-outcome`}>{key(`outcome_${operation.write_outcome}`)}</dd></>}
      <dt>{key('write_outcome')}</dt><dd data-testid={`${testIdPrefix}-write-phase`}>{key(write)}</dd>
      <dt>{key('readback')}</dt><dd data-testid={`${testIdPrefix}-readback`}>{key(readback)}</dd>
      <dt>{key('cleanup_status')}</dt>
      <dd data-testid={`${testIdPrefix}-cleanup`}>{key(cleanupClaim ? 'cleanup_repair_running' : `cleanup_${operation.cleanup_status ?? 'unconfirmed'}`)}</dd>
      {cleanupClaim && <><dt>{key('prior_cleanup')}</dt><dd data-testid={`${testIdPrefix}-prior-cleanup`}>{key(`cleanup_${operation.cleanup_status}`)}</dd></>}
    </dl>
    {manualGuidance && <p className="mt-2 text-xs text-amber-200" data-testid={`${testIdPrefix}-manual-guidance`}>{key('recovery_manual_guidance')}</p>}
  </div>;
}
