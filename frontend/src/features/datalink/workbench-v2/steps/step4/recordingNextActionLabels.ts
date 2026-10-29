import { useTranslation } from 'react-i18next';
import type { RecordingNextActionLabels } from './RecordingNextAction';

export function useRecordingNextActionLabels(): RecordingNextActionLabels {
  const { t } = useTranslation('workbench-v2-recording-guided');
  const actions = ['choose_sources', 'resolve_semantics', 'save_prepare', 'review_schema', 'preview_schema', 'confirm_schema',
    'review_scoped_start', 'query_operation', 'recover_operation', 'refresh_evidence', 'monitor', 'inspect'] as const;
  const reasons = ['sources_empty', 'counter_confirmation', 'save_required', 'offline', 'operation_recovery', 'evidence_loading',
    'evidence_error', 'evidence_unavailable', 'schema_review', 'schema_confirmation', 'schema_conflict', 'configuration_blocked',
    'readiness_blocked', 'lifecycle_unavailable', 'lifecycle_blocked', 'start_review', 'monitor_available', 'scope_unavailable'] as const;
  return { heading: t('next.heading'), scopes: { draft: t('next.scopes.draft'), group: t('next.scopes.group'), managed_device: t('next.scopes.managed_device') },
    actions: Object.fromEntries(actions.map(key => [key, t(`next.actions.${key}`)])) as RecordingNextActionLabels['actions'],
    reasons: Object.fromEntries(reasons.map(key => [key, t(`next.reasons.${key}`)])) as RecordingNextActionLabels['reasons'] };
}
