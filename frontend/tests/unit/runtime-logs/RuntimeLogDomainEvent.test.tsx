import { cleanup, fireEvent, render, screen, within } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import RuntimeLogDomainEvent from '@/features/runtime-logs/RuntimeLogDomainEvent';
import { parseStudioEntityLocator } from '@/services/studioEntityLocator';
import i18n from '@/i18n/config';
import { domainCases, domainInstant, domainRecord } from './domain-event-fixtures';

const first = '2026-10-29T09:00:00.000000001Z';
function mount(record = domainRecord(), navigateTo?: (target: string) => void, router = false) {
  const table = <table><tbody><RuntimeLogDomainEvent record={record} navigateTo={navigateTo} /></tbody></table>;
  return render(router ? <MemoryRouter>{table}</MemoryRouter> : table);
}
function details() {
  const summary = screen.getByText('Event details'); fireEvent.click(summary);
  expect(summary.tagName).toBe('SUMMARY'); expect(summary.parentElement).toHaveAttribute('open');
  return summary.parentElement;
}
function value(label: string) {
  const term = screen.getByText(label); const description = term.nextElementSibling;
  if (!description || term.tagName !== 'DT' || description.tagName !== 'DD') throw new Error(`Missing semantic description for ${label}`);
  return description;
}
const meanings: Record<string, RegExp> = {
  'domain.source.applied': /source.*applied/i, 'domain.source.rejected': /source.*reject/i,
  'domain.collector.failed': /acquisition.*fail/i, 'domain.collector.recovered': /acquisition.*recover/i,
  'domain.mapping.failed': /mapping.*fail/i, 'domain.mapping.recovered': /mapping.*recover/i,
  'domain.group.applied': /intake.*active|active.*intake/i, 'domain.group.paused': /intake.*paus/i,
  'domain.group.retiring': /retiring/i, 'domain.delivery.retrying': /delivery.*retry/i,
  'domain.delivery.blocked': /delivery.*block/i, 'domain.delivery.quarantined': /quarantin/i,
  'domain.delivery.unknown': /delivery.*unknown/i, 'domain.delivery.first-confirmed': /first.*commit/i,
  'domain.delivery.recovered': /delivery.*recover/i, 'domain.delivery.summary': /commit.*summary|summary.*commit/i,
  'domain.counter-replay.quarantined': /input.*quarantin|quarantin.*input/i,
  'domain.recording.refused': /recording.*refus/i, 'domain.quota.refused': /quota.*refus/i,
  'domain.quota.recovered': /capacity/i, 'domain.repair.completed': /repair.*complet/i,
  'domain.repair.blocked': /repair.*block/i, 'domain.cutover.applied': /cutover.*applied/i,
  'domain.cutover.rejected': /cutover.*reject/i,
};
beforeEach(async () => { await i18n.changeLanguage('en'); vi.stubGlobal('fetch', vi.fn()); });
afterEach(() => { cleanup(); vi.unstubAllGlobals(); });

describe('safe localized full domain rows preserve historical evidence', () => {
  it.each(domainCases)('renders the fixed meaning for %s without using candidate message or raw JSON', code => {
    const record = domainRecord(code, { message: 'password_fixture <script>private.example</script>' });
    const { container } = mount(record); const row = screen.getByRole('row');
    expect(row.textContent).toMatch(meanings[code]);
    expect(row).toHaveTextContent('Historical event');
    expect(row.textContent).not.toContain(record.message); expect(container.querySelector('script,img')).toBeNull();
    expect(container.textContent).not.toContain('password_fixture'); expect(container.textContent).not.toContain('private.example');
    expect(row.textContent).not.toContain('{"'); expect(fetch).not.toHaveBeenCalled();
    if (record.fields.action === 'none') expect(within(row).queryByRole('link')).not.toBeInTheDocument();
  });
  it('keeps original identities/revisions, exact uint64 occurrence count and nanosecond times distinct from transport time', () => {
    mount(domainRecord('domain.delivery.blocked', { fields: { first_seen_at: first, last_seen_at: domainInstant, occurrences: '9007199254740993' } }));
    details();
    for (const [label, expected] of [['Workspace ID', 'W'], ['Recording group ID', 'G'], ['Original entity revision', 'R1'],
      ['Observed workspace revision', 'W1'], ['Original group revision', 'R1'], ['Original connector revision', 'C1'],
      ['First seen', first], ['Last seen', domainInstant], ['Reference ID', 'REF1']]) expect(value(label)).toHaveTextContent(expected);
    expect(screen.getAllByText('9007199254740993').length).toBeGreaterThan(0);
    expect(screen.queryByText('9007199254740992')).not.toBeInTheDocument();
    expect(screen.getAllByText(domainInstant).length).toBeGreaterThan(0);
    const link = screen.getByRole('link', { name: 'Open group settings' });
    const href = link.getAttribute('href'); if (!href) throw new Error('Missing settings target');
    expect(parseStudioEntityLocator(new URL(href, 'http://localhost').search)).toMatchObject({ kind: 'locator', locator: {
      workspace_id: 'W', group_id: 'G', group_revision: 'R1', device_id: null,
    } });
    expect(href).not.toMatch(/connector|reference|effect|receipt|payload|return/); expect(fetch).not.toHaveBeenCalled();
  });
  it('accepts max uint64 as exact text and does not invent missing summary from transport timestamp', () => {
    const view = mount(domainRecord()); details();
    expect(value('First seen')).toHaveTextContent('Unavailable'); expect(value('Last seen')).toHaveTextContent('Unavailable');
    expect(screen.getAllByText('Unavailable').length).toBeGreaterThanOrEqual(3);
    view.unmount(); mount(domainRecord('domain.delivery.blocked', { fields: { first_seen_at: first, last_seen_at: domainInstant, occurrences: '18446744073709551615' } }));
    expect(screen.getAllByText('18446744073709551615').length).toBeGreaterThan(0);
  });
  it('group-wide evidence never invents an affected device, rule, point or mapping', () => {
    mount(); details();
    expect(screen.queryByText('Device ID')).not.toBeInTheDocument(); expect(screen.queryByText('Source rule ID')).not.toBeInTheDocument();
    expect(screen.queryByText('Mapping ID')).not.toBeInTheDocument();
    expect(screen.getByRole('link').getAttribute('href')).not.toContain('device_id');
  });
  it.each(['domain.source.rejected', 'domain.collector.failed', 'domain.mapping.failed'] as const)('%s opens only current exact source settings; no invented revision URL hint', code => {
    mount(domainRecord(code)); details();
    expect(value('Original entity revision')).toHaveTextContent('R1');
    const link = screen.getByRole('link', { name: 'Open current settings' }); const href = link.getAttribute('href');
    expect(href).not.toMatch(/revision|group_id/); expect(href).toContain('device_id=D');
    expect(screen.getByText(/historical.*snapshot|snapshot.*historical/i)).toBeVisible();
  });
  it('receipt details prove their own effect only and summary count is not aggregate SQL proof', () => {
    mount(domainRecord('domain.delivery.summary', { fields: { first_seen_at: first, last_seen_at: domainInstant, occurrences: '3' } }));
    details(); expect(value('Receipt record ID')).toHaveTextContent('REC1'); expect(value('Receipt effect key')).toHaveTextContent('REF1');
    expect(value('Receipt committed at')).toHaveTextContent(domainInstant);
    expect(screen.getByText(/receipt.*only.*own effect/i)).toBeVisible();
    expect(screen.getByText(/observed.*transition|observed.*commit/i)).toBeVisible();
    expect(screen.queryByRole('link')).not.toBeInTheDocument();
  });
  it('quota recovery is capacity only, repair completion is not data restoration, and paused intake may retain backlog', () => {
    const view = mount(domainRecord('domain.quota.recovered'));
    expect(screen.getByText(/not.*intake.*SQL|not.*SQL.*intake/i)).toBeVisible(); expect(screen.queryByRole('link')).not.toBeInTheDocument();
    view.unmount(); const repair = mount(domainRecord('domain.repair.completed'));
    expect(screen.getByText(/data.*recover.*not|not.*data.*recover/i)).toBeVisible(); repair.unmount();
    mount(domainRecord('domain.group.paused'));
    expect(screen.getByText(/historical.*pending|backlog|delivery.*pending/i)).toBeVisible();
  });
  it.each([
    ['unsafe ID', domainRecord('domain.delivery.blocked', { fields: { group_id: 'db.private' } })],
    ['unknown field', domainRecord('domain.delivery.blocked', { fields: { raw_error: 'password_fixture' } })],
    ['wrong receipt', domainRecord('domain.delivery.summary', { fields: { effect_key: 'OTHER' } })],
    ['wrong stage', domainRecord('domain.delivery.blocked', { fields: { stage: 'sql-committed' } })],
    ['unsafe outer timestamp', domainRecord('domain.delivery.blocked', { timestamp: '2026-10-29 (private.example)' })],
    ['unsafe outer source', domainRecord('domain.delivery.blocked', { source: 'private.example' })],
    ['unknown domain', { ...domainRecord(), code: 'domain.password_fixture' }],
    ['partial summary', domainRecord('domain.delivery.blocked', { fields: { first_seen_at: first } })],
  ])('renders %s as a complete unavailable row without raw outer cells or href', (_name, record) => {
    const { container } = mount(record); expect(screen.getByRole('row')).toHaveTextContent(/unavailable/i);
    expect(screen.queryByRole('link')).not.toBeInTheDocument();
    expect(container.textContent).not.toMatch(/private\.example|db\.private|password_fixture|raw_error|OTHER/);
    expect(container.querySelector('script,img')).toBeNull(); expect(fetch).not.toHaveBeenCalled();
  });
  it('fixed identity-free unavailable never borrows candidate IDs or invents a settings destination', () => {
    mount(domainRecord('domain.unavailable'));
    expect(screen.getByRole('row')).toHaveTextContent(/unavailable/i); expect(screen.queryByRole('link')).not.toBeInTheDocument();
    expect(screen.queryByText('Recording group ID')).not.toBeInTheDocument(); expect(fetch).not.toHaveBeenCalled();
  });
  it('standalone navigation callback receives only the validated same-origin target, never an operation', () => {
    const navigate = vi.fn(); mount(domainRecord(), navigate);
    const link = screen.getByRole('link', { name: 'Open group settings' }); fireEvent.click(link, { button: 0 });
    expect(navigate).toHaveBeenCalledExactlyOnceWith(link.getAttribute('href')); expect(fetch).not.toHaveBeenCalled();
  });
  it.each(['en', 'zh-TW'])('offers native keyboard-focusable disclosure and localized guidance in %s', async language => {
    await i18n.changeLanguage(language); const { container } = mount(domainRecord(), undefined, true);
    const summary = container.querySelector('summary'); if (!(summary instanceof HTMLElement)) throw new Error('Missing native disclosure');
    summary.focus(); expect(summary).toHaveFocus(); fireEvent.click(summary, { detail: 0 });
    expect(summary.parentElement).toHaveAttribute('open');
    const link = screen.getByRole('link'); link.focus(); expect(link).toHaveFocus();
    expect(link).toHaveAccessibleName(language === 'en' ? 'Open group settings' : /群組/);
    expect(summary).toHaveTextContent(language === 'en' ? 'Event details' : /事件.*詳細/);
    expect(screen.getByRole('row')).toHaveTextContent(language === 'en' ? 'Delivery blocked' : '資料交付受阻');
    expect(fetch).not.toHaveBeenCalled();
  });
});
