import { useTranslation } from 'react-i18next';
import type { ProtocolType, SourceRuleRecord } from '../../../../types/datalink';
import type { WorkbenchV2State } from '../state/types';
import type { HomeSourceTarget } from './RecordingHomeSourceHost';
import { describeMemoryRule, type MemoryRuleInput } from './recordingHomeMemoryMap';

/** Inline range choice inside the one settings frame. Labels come from the scoped
 * range list; choosing one still opens it through its own scoped locator read. */
export function RecordingHomeRangeOptions({ protocol, listed, inputs, state, target, ready, failed, onRetry, onSelect }: {
  protocol: ProtocolType; listed?: readonly SourceRuleRecord[]; inputs: readonly MemoryRuleInput[]; state: WorkbenchV2State; target: HomeSourceTarget;
  /** Fail closed: only a settled, successful scoped list may start another range or a new one. */
  ready: boolean; failed: boolean; onRetry: () => void; onSelect: (ruleId?: string, newRange?: boolean) => void;
}) {
  const { t } = useTranslation('recording-home');
  return <div className="gw-home-range-options-wrap">
    <div role="group" aria-label={t('source.ranges')} className="gw-home-range-options" data-testid="home-range-options">
      {listed?.map(item => {
        const rule = state.rules.find(saved => saved.id === item.id && saved.workspace_id === target.workspaceId && saved.device_id === target.deviceId && saved.save_state === 'saved') ?? item;
        const input = inputs.find(value => value.id === item.id), span = input ? describeMemoryRule(protocol, { ...input, label: rule.naming_prefix }) : undefined;
        const address = span?.end && span.end !== span.start ? `${span.start} – ${span.end}` : span?.start ?? rule.start_address;
        const footprint = span?.cells !== undefined ? t(`source.memory.${span.unit}`, { count: span.cells }) : t('source.memory.footprintUnknown');
        const facts = [t('source.memory.items', { count: rule.count }), rule.data_type, footprint,
          ...(rule.skipped_addresses.length ? [t('source.memory.skipped', { count: rule.skipped_addresses.length })] : [])];
        const selected = item.id === target.ruleId;
        return <button type="button" key={item.id} data-testid={`home-range-option-${item.id}`} aria-pressed={selected} disabled={!ready && !selected} onClick={() => { if (ready && !selected) onSelect(item.id); }}>
          <span className="gw-range-option-name"><strong>{rule.naming_prefix}</strong>{!rule.enabled && <span className="gw-range-flag">{t('source.memory.disabled')}</span>}</span>
          <span className="gw-range-option-facts"><span className="gw-range-address">{address}</span> · {facts.join(' · ')}</span></button>;
      })}
      <button type="button" data-testid="home-range-option-new" aria-pressed={Boolean(target.newRange)} disabled={!ready && !target.newRange} onClick={() => { if (ready && !target.newRange) onSelect(undefined, true); }}>{t('source.newRange')}</button>
    </div>
    {!ready && !failed && <p role="status">{t('loading')}</p>}
    {failed && <p role="alert">{t('source.rangesUnavailable')} <button type="button" onClick={onRetry}>{t('source.memory.retry')}</button></p>}
  </div>;
}
