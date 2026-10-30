import { createElement } from 'react';
import { cleanup, render, screen } from '@testing-library/react';
import { MemoryRouter, useLocation } from 'react-router-dom';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { AppRoutes } from '@/App';
const received = vi.hoisted(() => ({ props: {} as Record<string, unknown> }));
vi.mock('@/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page', () => ({ default: function Receiver(props: Record<string, unknown>) {
  received.props = props; const location = useLocation();
  return createElement('div', { 'data-testid': 'receiver' }, location.pathname + location.search + location.hash);
} }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
afterEach(cleanup);
const locator = '?entry=locator&workspace_id=W&step=4&group_id=G&group_revision=R';
async function open(target: string) {
  render(createElement(MemoryRouter, { initialEntries: [target] }, createElement(AppRoutes)));
  return await screen.findByTestId('receiver');
}
describe('Studio presentation routing does not weaken viewing identity', () => {
  it('opts into recordings without rewriting query or hash', async () => {
    const target = '/studio/v2?view=recordings&section=database#saved';
    expect(await open(target)).toHaveTextContent(target);
    expect(received.props.presentation).toBe('recordings');
    expect(received.props.locatorIntent).toEqual({ kind: 'absent' });
  });
  it.each(['recordings', 'advanced'])('valid locator precedes %s view', async view => {
    const target = `/studio/v2${locator}&view=${view}#saved`;
    expect(await open(target)).toHaveTextContent(target);
    expect(received.props.locatorIntent).toMatchObject({ kind: 'locator', locator: { workspace_id: 'W', group_id: 'G', group_revision: 'R' } });
    expect(received.props.presentation).toBe('advanced');
  });
  it.each([
    `${locator}&view=recordings&view=recordings`, `${locator}&view=recordings&%76iew=advanced`,
    `${locator}&view=unknown`, `${locator}&view=recordings&bad=1`, `${locator}&view=%ZZ`,
    `${locator}&view=recordings&device_id=%E0%A4%A`, '?view=recordings&workspace_id=W',
    '?view=recordings&view=advanced', '?view=recordings&section=%ZZ', '?view=recordings&section=x&section=x',
  ])('invalid request stays invalid rather than opening home: %s', async search => {
    await open(`/studio/v2${search}`);
    expect(received.props.locatorIntent).toEqual({ kind: 'invalid' });
    expect(received.props.presentation).toBe('advanced'); expect(received.props.runtimeReturnFocus).toBeNull();
  });
  it('runtime readiness focus precedes home without losing the repair target', async () => {
    await open('/studio/v2?focus=readiness&step=3&issue=mapping-required&view=recordings#repair');
    expect(received.props.runtimeReturnFocus).toEqual({ step: 3, issueCode: 'mapping-required' });
    expect(received.props.presentation).toBe('advanced');
  });
  it.each(['', '?view=advanced', '?view=unknown'])('default and unsupported non-locator view retain advanced entry: %s', async search => {
    await open(`/studio/v2${search}`); expect(received.props.presentation).toBe('advanced');
    expect(received.props.locatorIntent).toEqual({ kind: 'absent' });
  });
  it.each(['/', '/datalink', '/datalink/workbench', '/unrecognized', '/studio'])('generic redirect preserves search/hash from %s', async path => {
    expect(await open(`${path}?view=recordings#saved`)).toHaveTextContent('/studio/v2?view=recordings#saved');
    expect(received.props.presentation).toBe('recordings');
  });
});

import { buildStudioPresentationTarget, parseStudioPresentation } from '@/services/studioRecordingView';
describe('presentation-only target builder', () => {
  it('preserves unrelated raw encoding and hash while replacing only the view', () => {
    expect(buildStudioPresentationTarget('?section=a%20b&%76iew=recordings', '#return%20here', 'advanced'))
      .toBe('/studio/v2?section=a%20b&view=advanced#return%20here');
  });
  it.each(['?view=x&view=y', '?section=%ZZ', '?entry=locator&workspace_id=W&view=recordings'])('cannot repair invalid intent through navigation: %s', search => {
    expect(buildStudioPresentationTarget(search, '#saved', 'advanced')).toBeNull();
  });
  it('does not expand the standalone locator parser permissions', () => {
    expect(parseStudioPresentation(`${locator}&view=advanced`).locatorIntent.kind).toBe('locator');
    expect(parseStudioPresentation(`${locator}&view=unknown`).locatorIntent.kind).toBe('invalid');
  });
});
