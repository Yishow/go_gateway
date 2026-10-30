import { afterEach, expect, it } from 'vitest';
import { captureReturnFocus, restoreReturnFocus, markReturnRegion, eligibleFocusTarget } from '@/utils/returnFocus';
afterEach(() => { document.body.replaceChildren(); });
function setup() {
  const region = document.createElement('section'); region.dataset.focusScope = 'W/G/R1';
  const target = document.createElement('button'); const fallback = document.createElement('h2'); fallback.tabIndex = -1;
  region.append(target, fallback); document.body.append(region); return { region, target, fallback };
}
it('returns to an eligible same-owner trigger without a second presentation store', () => {
  const { target, fallback } = setup(); const owner = {}; const ticket = captureReturnFocus(owner, target);
  expect(restoreReturnFocus(ticket, owner, fallback)).toBe('trigger'); expect(target).toHaveFocus();
});
it.each(['removed', 'disabled', 'hidden', 'inert', 'aria-hidden', 'css-hidden', 'revision'])('%s trigger falls back without reviving its authority', kind => {
  const { region, target, fallback } = setup(); const owner = {}; const ticket = captureReturnFocus(owner, target);
  if (kind === 'removed') target.remove(); if (kind === 'disabled') target.disabled = true;
  if (kind === 'hidden') target.hidden = true; if (kind === 'inert') target.setAttribute('inert', '');
  if (kind === 'aria-hidden') target.setAttribute('aria-hidden', 'true'); if (kind === 'css-hidden') target.style.display = 'none';
  if (kind === 'revision') region.dataset.focusScope = 'W/G/R2';
  expect(restoreReturnFocus(ticket, owner, fallback)).toBe('fallback'); expect(fallback).toHaveFocus();
});
it('foreign owner and hidden fallback never steal focus across contexts', () => {
  const { target, fallback } = setup(); const owner = {}; const ticket = captureReturnFocus(owner, target);
  expect(restoreReturnFocus(ticket, {}, fallback)).toBe('none'); expect(target).not.toHaveFocus();
  target.remove(); fallback.hidden = true; expect(restoreReturnFocus(ticket, owner, fallback)).toBe('none');
});
it('a removed/revised entity can fall back to its original still-current workspace heading, never a foreign workspace heading', () => {
  const { region, target } = setup(); const page = document.createElement('main'); page.dataset.focusReturnRegion = 'W/view';
  const heading = document.createElement('h2'); heading.tabIndex = -1; heading.dataset.focusReturnHeading = '';
  document.body.append(page); page.append(heading, region); markReturnRegion(page, {}); const ticket = captureReturnFocus({}, target);
  target.remove(); expect(restoreReturnFocus(ticket, {})).toBe('fallback'); expect(heading).toHaveFocus();
  heading.blur(); page.dataset.focusReturnRegion = 'FOREIGN/view'; expect(restoreReturnFocus(ticket, {})).toBe('none'); expect(heading).not.toHaveFocus();
});

it('workspace A→B→A cannot revive a captured region owner merely because its label matches again', () => {
  const { region, target, fallback } = setup(); region.dataset.focusReturnRegion = 'A'; fallback.dataset.focusReturnHeading = '';
  const ownerA = {}; markReturnRegion(region, ownerA); const ticket = captureReturnFocus(ownerA, target);
  markReturnRegion(region, {}); region.dataset.focusReturnRegion = 'B';
  markReturnRegion(region, {}); region.dataset.focusReturnRegion = 'A'; target.remove();
  expect(restoreReturnFocus(ticket, {})).toBe('none'); expect(fallback).not.toHaveFocus();
});

it('closed native details rejects content and local fallback but preserves its first direct summary and descendants', () => {
  const page = document.createElement('main'); page.dataset.focusReturnRegion = 'W/view';
  const heading = document.createElement('h2'); heading.tabIndex = -1; heading.dataset.focusReturnHeading = '';
  const details = document.createElement('details'); details.open = true;
  const summary = document.createElement('summary'); const label = document.createElement('span'); summary.append(label);
  const second = document.createElement('summary'); const trigger = document.createElement('button'); const local = document.createElement('h3'); local.tabIndex = -1;
  details.append(summary, second, trigger, local); page.append(heading, details); document.body.append(page);
  const owner = {}; markReturnRegion(page, owner); const ticket = captureReturnFocus(owner, trigger); details.open = false;
  expect(eligibleFocusTarget(summary)).toBe(true); expect(eligibleFocusTarget(label)).toBe(true); expect(eligibleFocusTarget(second)).toBe(false);
  expect(restoreReturnFocus(ticket, owner, local)).toBe('fallback'); expect(heading).toHaveFocus();
});
it('every closed details ancestor is considered, including nested summaries hidden by an outer disclosure', () => {
  const outer = document.createElement('details'); const outerSummary = document.createElement('summary');
  const inner = document.createElement('details'); const innerSummary = document.createElement('summary'); const content = document.createElement('button');
  inner.append(innerSummary, content); outer.append(outerSummary, inner); document.body.append(outer);
  expect(eligibleFocusTarget(innerSummary)).toBe(false); outer.open = true; expect(eligibleFocusTarget(innerSummary)).toBe(true);
  expect(eligibleFocusTarget(content)).toBe(false); inner.open = true; expect(eligibleFocusTarget(content)).toBe(true);
  outer.open = false; expect(eligibleFocusTarget(content)).toBe(false);
});
