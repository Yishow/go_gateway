import { createElement } from 'react';
import { cleanup, render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { AppRoutes } from '@/App';
const received = vi.hoisted(() => ({ props: {} as Record<string, unknown> }));
vi.mock('@/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page', () => ({ default: (props: Record<string, unknown>) => {
  received.props = props; return createElement('div', { 'data-testid': 'locator-receiver' });
} }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
afterEach(cleanup);
describe('existing Studio route receives safe locator intent', () => {
  it('passes an exact group locator instead of silently falling back to ordinary setup', async () => {
    render(createElement(MemoryRouter, { initialEntries: ['/studio/v2?entry=locator&workspace_id=W&step=4&group_id=G1&group_revision=R1'] }, createElement(AppRoutes)));
    await screen.findByTestId('locator-receiver');
    expect(received.props.locatorIntent).toEqual({ kind: 'locator', locator: { step: 4, workspace_id: 'W', device_id: null, rule_id: null, mapping_id: null, group_id: 'G1', group_revision: 'R1' } });
    expect(received.props.runtimeReturnFocus).toBeNull();
  });
});

import { buildStudioEntityLocatorTarget, parseStudioEntityLocator } from '@/services/studioEntityLocator';
import type { EntityLocator } from '@/types/runtimeRuleEvidence';
const empty: EntityLocator = { step: 1, workspace_id: 'W', device_id: null, rule_id: null, mapping_id: null, group_id: null, group_revision: null };
describe('one strict IDs-only locator URL contract', () => {
  it.each([
    { ...empty, step: 1, device_id: 'D' }, { ...empty, step: 2, device_id: 'D', rule_id: 'R' },
    { ...empty, step: 3, device_id: 'D', mapping_id: 'M' }, { ...empty, step: 4, group_id: 'G' },
    { ...empty, step: 4, device_id: 'D', group_id: 'G', group_revision: 'OLD' },
  ] as EntityLocator[])('round-trips exact step $step without adding context', (locator) => {
    const target = buildStudioEntityLocatorTarget(locator); expect(target).toMatch(/^\/studio\/v2\?entry=locator&workspace_id=W&step=[1-4]/);
    expect(parseStudioEntityLocator(target!.slice('/studio/v2'.length))).toEqual({ kind: 'locator', locator });
  });
  it.each(['', '?focus=readiness&step=3&issue=source-required', '?section=database'])('preserves ordinary/readiness entry when locator is absent: %s', (search) => {
    expect(parseStudioEntityLocator(search)).toEqual({ kind: 'absent' });
  });
  it.each([
    '?workspace_id=W&step=1&device_id=D', '?step=1', '?entry=other&workspace_id=W&step=1&device_id=D',
    '?editor=1', '?expected_group_revision=R', '?expected_workspace_revision=', '?locator_workspace_id=W&locator_workspace_revision=W1',
    '?entry=locator&workspace_id=W&step=2&device_id=D', '?entry=locator&workspace_id=W&step=3&device_id=D&mapping_id=M&rule_id=R',
    '?entry=locator&workspace_id=W&step=1&device_id=D&group_id=G', '?entry=locator&workspace_id=W&step=4&group_revision=R',
    '?entry=locator&workspace_id=W&step=1&device_id=D&device_id=D', '?entry=locator&workspace_id=W&step=1&device_id=D&device_id=%ZZ',
    '?entry=%ZZ&workspace_id=%ZZ', '?entry=locator;workspace_id=W&step=1&device_id=D',
    '?entry=locator&workspace_id=W&step=1&device_id=D&focus=readiness', '?entry=locator&workspace_id=W&step=1&device_id=D&return=https%3A%2F%2Fevil.invalid',
    '?entry=locator&workspace_id=W&step=1&device_id=private.local', '?entry=locator&workspace_id=W&step=1&device_id=%2Fprivate%2Fpath',
    '?entry=locator&workspace_id=W&step=01&device_id=D', '?entry=locator&workspace_id=W&step=4&group_id=G&editor=1',
  ])('rejects malformed/mixed/unsafe query without fallback: %s', (search) => expect(parseStudioEntityLocator(search)).toEqual({ kind: 'invalid' }));
  it.each(['http://private.local', '/private/file', 'password=secret', '192.168.1.1', '', ' D'])('never emits unsafe identity %s', (device_id) => {
    expect(buildStudioEntityLocatorTarget({ ...empty, device_id })).toBeNull();
  });
  it('does not copy unknown payloads or editor flags into a URL', () => {
    expect(buildStudioEntityLocatorTarget({ ...empty, device_id: 'D', editor: true } as EntityLocator)).toBeNull();
  });
  it('keeps readiness focus intact only when no locator intent is supplied', async () => {
    render(createElement(MemoryRouter, { initialEntries: ['/studio/v2?focus=readiness&step=2&issue=source-required'] }, createElement(AppRoutes)));
    await screen.findByTestId('locator-receiver'); expect(received.props.runtimeReturnFocus).toEqual({ step: 2, issueCode: 'source-required' });
    expect(received.props.locatorIntent).toEqual({ kind: 'absent' });
  });
  it('passes invalid rather than silently invoking readiness or ordinary setup for a mixed intent', async () => {
    render(createElement(MemoryRouter, { initialEntries: ['/studio/v2?entry=locator&workspace_id=W&step=2&device_id=D&rule_id=R&focus=readiness'] }, createElement(AppRoutes)));
    await screen.findByTestId('locator-receiver'); expect(received.props.locatorIntent).toEqual({ kind: 'invalid' }); expect(received.props.runtimeReturnFocus).toBeNull();
  });
});

it('rejects accessor and inherited locators without evaluating them', () => {
  const get = vi.fn(() => 'D'); const value = { ...empty }; Object.defineProperty(value, 'device_id', { get });
  expect(buildStudioEntityLocatorTarget(value)).toBeNull(); expect(get).not.toHaveBeenCalled();
  expect(buildStudioEntityLocatorTarget(Object.create({ ...empty, device_id: 'D' }))).toBeNull();
});
