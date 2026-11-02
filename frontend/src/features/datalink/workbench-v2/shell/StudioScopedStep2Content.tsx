import { useRef, useState, type ReactNode } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { readHomeDeviceRanges } from '../recording-home/recordingHomeSourceList';
import { buildStudioEntityLocatorTarget } from '../../../../services/studioEntityLocator';
import type { EvidenceScope } from '../../../../types/runtimeRuleEvidence';

/** Mount keyed by the accepted W/revision/D only. This existing list reader
 * brackets workspace membership; its rows never hydrate or authorize an editor. */
export function StudioScopedStep2Content({ scope, children, navigateTo, hasLocalDraft }: {
  scope: EvidenceScope; children: ReactNode; navigateTo: (target: string) => void;
  hasLocalDraft: boolean;
}) {
  const { t } = useTranslation('recording-home');
  const [adding, setAdding] = useState(false);
  const query = useQuery({
    queryKey: ['studio-scoped-step2-ranges', scope.workspace_id, scope.workspace_revision, scope.device_id],
    queryFn: () => readHomeDeviceRanges(scope.workspace_id, scope.device_id),
    retry: false, refetchOnWindowFocus: false, gcTime: 0,
  });
  const current = useRef(query); current.current = query;
  const ready = query.isSuccess && !query.isFetching;
  const rows = ready ? query.data : undefined;
  const open = (id: string) => {
    const fresh = current.current;
    if (!fresh.isSuccess || fresh.isFetching || fresh.data !== rows ||
        !rows?.some(rule => rule.id === id && rule.workspace_id === scope.workspace_id && rule.device_id === scope.device_id)) return;
    const target = buildStudioEntityLocatorTarget({ step: 2, workspace_id: scope.workspace_id,
      device_id: scope.device_id, rule_id: id, mapping_id: null, group_id: null, group_revision: null });
    if (target) navigateTo(target);
  };
  // Transport refresh/failure cannot hide an already-created local rule form.
  if (!rows) return <><section data-testid="studio-scoped-step2-list" role={query.isPending || query.isFetching ? 'status' : 'alert'}>
    {t(query.isPending || query.isFetching ? 'loading' : 'source.unavailable')}
  </section>{(adding || hasLocalDraft) && children}</>;
  return <section data-testid="studio-scoped-step2-list">
    <h3>{t('source.ranges')}</h3>
    <ul>{rows.map(rule => <li key={rule.id}>
      <span>{rule.naming_prefix} · {rule.start_address} · {rule.count} · {rule.data_type}</span>
      <button type="button" data-testid={`studio-scoped-rule-${rule.id}`} onClick={() => open(rule.id)}>{t('source.edit')}</button>
    </li>)}</ul>
    {!adding && rows.length > 0 && <button type="button" onClick={() => setAdding(true)}>{t('source.newRange')}</button>}
    {(adding || rows.length === 0 || hasLocalDraft) && children}
  </section>;
}
