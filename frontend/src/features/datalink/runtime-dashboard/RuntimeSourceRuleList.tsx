import { RuntimeSettingsLink } from './RuntimeSettingsLink';
import { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useRuntimeSourcePage } from '../../../hooks/datalink/useRuntimeSourcePage';
import type { EvidenceScope } from '../../../types/runtimeRuleEvidence';
import type { RuntimeSourcePage, RuntimeSourcePageRequest, RuntimeSourceSelection, StudioV2RuntimeSetupSourceRule } from '../../../types/studioV2RuntimeContext';
import { RuntimeSourceRuleDetail } from './RuntimeSourceRuleDetail';

type Kind = 'rules' | 'mappings';
type Paging = { query: string; draft: string; cursors: Array<string | undefined>; index: number };
const initial = (): Paging => ({ query: '', draft: '', cursors: [undefined], index: 0 });
const request = (kind: Kind, paging: Paging): RuntimeSourcePageRequest => ({ kind, ...(paging.query ? { query: paging.query } : {}), ...(paging.cursors[paging.index] ? { cursor: paging.cursors[paging.index] } : {}) });

/** Bounded HTTP pages only; existing source authorities may still scan their saved lists. */
export function RuntimeSourceRuleList({ scope, navigateTo, initialSelection = null, onSelection }: {
  scope: EvidenceScope; navigateTo?: (target: string) => void;
  initialSelection?: RuntimeSourceSelection | null; onSelection?: (selection: RuntimeSourceSelection | null) => void;
}) {
  const { t } = useTranslation('runtime-dashboard');
  const [rulePaging, setRulePaging] = useState(initial);
  const [mappingPaging, setMappingPaging] = useState(initial);
  const [filter, setFilter] = useState<{ rule_id: string; rule_revision: string } | null>(null);
  const [selection, setSelection] = useState<RuntimeSourceSelection | null>(initialSelection);
  const [restored, setRestored] = useState(Boolean(initialSelection));
  const opener = useRef<HTMLButtonElement | null>(null);
  const ruleQuery = useRuntimeSourcePage(scope, request('rules', rulePaging));
  const mappingQuery = useRuntimeSourcePage(scope, { ...request('mappings', mappingPaging), ...filter });
  const rulePage = !ruleQuery.isFetching && ruleQuery.data?.ok ? ruleQuery.data.value : null;
  const mappingPage = !mappingQuery.isFetching && mappingQuery.data?.ok ? mappingQuery.data.value : null;
  // Keep a selected detail alive while its list is being revalidated. A missing
  // page during isFetching is a transient transport state, not proof that the
  // selected source was removed. Once the request settles, a missing row or a
  // typed page failure is allowed to clear the stale selection.
  const mappingPageSettled = !mappingQuery.isFetching && (mappingQuery.data !== undefined || mappingQuery.isError);
  const selectedExists = mappingPage?.mappings.some((row) => row.rule_id === selection?.rule_id && row.rule_revision === selection?.rule_revision && row.point_id === selection?.point_id && (row.id || undefined) === selection?.mapping_id);
  const removed = Boolean(selection && !restored && mappingPageSettled && !selectedExists);
  useEffect(() => { if (removed) { setSelection(null); onSelection?.(null); } }, [removed, onSelection]);
  // A restored explicit selection may be outside the current page. The detail reader revalidates it.
  const select = (value: RuntimeSourceSelection | null) => { setRestored(false); setSelection(value); onSelection?.(value); };
  const clearSelection = () => select(null);
  const chooseRule = (rule: StudioV2RuntimeSetupSourceRule) => {
    if (!rule.revision_id) return;
    setFilter({ rule_id: rule.id, rule_revision: rule.revision_id });
    setMappingPaging(initial()); clearSelection();
  };
  const search = (kind: Kind) => {
    const paging = kind === 'rules' ? rulePaging : mappingPaging;
    const next = { ...initial(), draft: paging.draft, query: paging.draft.trim() };
    if (kind === 'rules') { setRulePaging(next); setFilter(null); setMappingPaging(initial()); }
    else setMappingPaging(next);
    clearSelection();
  };
  const move = (kind: Kind, direction: 'previous' | 'next') => {
    const paging = kind === 'rules' ? rulePaging : mappingPaging;
    const page = kind === 'rules' ? rulePage : mappingPage;
    const set = kind === 'rules' ? setRulePaging : setMappingPaging;
    if (direction === 'previous' && paging.index > 0) set({ ...paging, index: paging.index - 1 });
    if (direction === 'next' && page?.next_cursor) set({ ...paging, cursors: [...paging.cursors.slice(0, paging.index + 1), page.next_cursor], index: paging.index + 1 });
    clearSelection();
  };
  const restart = (kind: Kind) => {
    const paging = kind === 'rules' ? rulePaging : mappingPaging;
    const set = kind === 'rules' ? setRulePaging : setMappingPaging;
    clearSelection();
    if (paging.index === 0) void (kind === 'rules' ? ruleQuery : mappingQuery).refetch();
    else set({ ...paging, cursors: [undefined], index: 0 });
  };
  return <>
    {(['rules', 'mappings'] as const).map((kind) => {
      const paging = kind === 'rules' ? rulePaging : mappingPaging;
      const query = kind === 'rules' ? ruleQuery : mappingQuery;
      const page = kind === 'rules' ? rulePage : mappingPage;
      const reason = query.data && !query.data.ok ? query.data.reason : 'unavailable';
      return <section key={kind} aria-label={t(`sourceList.${kind}.title`)} className="rounded-2xl border border-slate-800 bg-slate-950/70 p-4 text-xs text-slate-200">
        <h3 className="font-semibold">{t(`sourceList.${kind}.title`)}</h3>
        <p className="mt-1 text-slate-400">{t('sourceList.savedOnly')}</p>
        {kind === 'mappings' && filter && <div className="mt-2 flex items-center gap-2">
          <span>{t('sourceList.filtered', { rule: filter.rule_id })}</span>
          <button type="button" onClick={() => { setFilter(null); setMappingPaging(initial()); clearSelection(); }} className="rounded border border-slate-700 px-2 py-1">{t('sourceList.clearFilter')}</button>
        </div>}
        <form className="mt-3 flex flex-wrap gap-2" onSubmit={(event) => { event.preventDefault(); search(kind); }}>
          <label className="flex grow flex-col gap-1">
            {t(`sourceList.${kind}.search`)}
            <input value={paging.draft} maxLength={128} onChange={(event) => (kind === 'rules' ? setRulePaging : setMappingPaging)({ ...paging, draft: event.target.value })} className="rounded border border-slate-700 bg-slate-900 px-2 py-1" />
          </label>
          <button type="submit" className="self-end rounded border border-cyan-700 px-3 py-1">{t(`sourceList.${kind}.submit`)}</button>
        </form>
        {query.isPending || query.isFetching ? <p role="status" className="mt-3">{t('sourceList.loading')}</p> : !page ? <div role="alert" className="mt-3 text-amber-100">
          <p>{t(`sourceDetail.errors.${reason}`)}</p>
          <button type="button" onClick={() => restart(kind)} className="mt-2 rounded border border-amber-700 px-2 py-1">{t('sourceList.restart')}</button>
        </div> : <SourceRows page={page} navigateTo={navigateTo} onRule={chooseRule} onSource={(value, button) => { opener.current = button; select(value); }} />}
        <div className="mt-3 flex items-center justify-between gap-2">
          <button type="button" aria-label={t(`sourceList.${kind}.previous`)} disabled={paging.index === 0 || query.isFetching} onClick={() => move(kind, 'previous')} className="rounded border border-slate-700 px-2 py-1 disabled:opacity-50">{t('sourceList.previous')}</button>
          <span>{t('sourceList.page', { page: paging.index + 1 })}</span>
          <button type="button" aria-label={t(`sourceList.${kind}.next`)} disabled={!page?.next_cursor || query.isFetching} onClick={() => move(kind, 'next')} className="rounded border border-slate-700 px-2 py-1 disabled:opacity-50">{t('sourceList.next')}</button>
        </div>
      </section>;
    })}
    {selection && !removed && <RuntimeSourceRuleDetail key={JSON.stringify(selection)} scope={scope} selection={selection} navigateTo={navigateTo} onClose={() => { clearSelection(); opener.current?.focus(); }} />}
  </>;
}

function SourceRows({ page, onRule, onSource, navigateTo }: {
  navigateTo?: (target: string) => void;
  page: RuntimeSourcePage;
  onRule: (rule: StudioV2RuntimeSetupSourceRule) => void;
  onSource: (selection: RuntimeSourceSelection, button: HTMLButtonElement) => void;
}) {
  const { t } = useTranslation('runtime-dashboard');
  const locator = { workspace_id: page.workspace_id, device_id: page.device_id, rule_id: null, mapping_id: null, group_id: null, group_revision: null };
  if (page.rules.length + page.mappings.length === 0) return <p className="mt-3 text-slate-400">{t('sourceList.empty')}</p>;
  return <ul className="mt-3 space-y-2">
    {page.rules.map((rule) => <li key={rule.id} className="rounded border border-slate-800 p-2">
      <div><span className="font-mono">{rule.id}</span> · {rule.start_address} · {rule.count} {t('setup.rules.pointUnit')} · {rule.data_type}</div>
      <div>{rule.naming_prefix}</div>
      <RuntimeSettingsLink locator={{ ...locator, step: 2, rule_id: rule.id }} navigateTo={navigateTo}>{t('settingsLink.rule')}</RuntimeSettingsLink>
      <button type="button" aria-label={t('sourceList.showMappings', { rule: rule.id })} onClick={() => onRule(rule)} className="mt-1 rounded border border-cyan-700 px-2 py-1">{t('sourceList.inspectRule')}</button>
    </li>)}
    {page.mappings.map((mapping) => <li key={`${mapping.rule_id}-${mapping.point_id}-${mapping.id ?? 'point-only'}`} className="rounded border border-slate-800 p-2">
      <div className="font-mono">{mapping.address}</div>
      <div>{mapping.display_name || mapping.tag_key || mapping.tag_id || t('setup.mappings.unbound')}</div>
      {mapping.display_name && mapping.tag_key && <div>{mapping.tag_key}</div>}
      <div>{mapping.rule_id} · {mapping.point_id}</div>
      {mapping.id && <RuntimeSettingsLink locator={{ ...locator, step: 3, mapping_id: mapping.id }} navigateTo={navigateTo}>{t('settingsLink.mapping')}</RuntimeSettingsLink>}
      <button type="button" aria-label={t('sourceDetail.inspectAddress', { address: mapping.address })} onClick={(event) => onSource({ rule_id: mapping.rule_id, rule_revision: mapping.rule_revision!, point_id: mapping.point_id, ...(mapping.id ? { mapping_id: mapping.id } : {}) }, event.currentTarget)} className="mt-1 rounded border border-cyan-700 px-2 py-1">{t('sourceDetail.inspect')}</button>
    </li>)}
  </ul>;
}
