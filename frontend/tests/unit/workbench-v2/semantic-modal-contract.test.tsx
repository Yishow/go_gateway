import * as React from 'react';
import { readFileSync } from 'node:fs';
import { cleanup, fireEvent, render, screen, within } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { RecordingChangeController } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/RecordingChangeController';
import { RecordingChangeDialog } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/RecordingChangeDialog';
import { SavedRecordingRemoveDialog } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/SavedRecordingRemoveDialog';
import { fixture } from '../../fixtures/recordingChange';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
afterEach(() => { cleanup(); document.documentElement.className = ''; });
function Change() {
  const f = React.useMemo(fixture, []); const controller = React.useMemo(() => new RecordingChangeController(f.workspace, f.deps), [f]);
  const state = React.useSyncExternalStore(controller.subscribe, controller.snapshot);
  return <><button onClick={() => controller.requestAction(f.intent)}>Open change</button><RecordingChangeDialog controller={{ ...state, state, controller, requestAction: controller.requestAction, dismiss: controller.dismiss, unresolved: controller.unresolved }} /></>;
}
describe('actual Radix consumers share the semantic modal contract, not a second primitive', () => {
  it.each(['light', 'dark'])('%s form portal uses explicit form sizing and fixed header/body/footer without operation changes', async theme => {
    document.documentElement.className = theme; const view = render(<Change />); fireEvent.click(screen.getByText('Open change'));
    const dialog = await screen.findByRole('dialog', { name: 'dialog.title' });
    expect(dialog).toHaveClass('gw-modal'); expect(dialog).toHaveAttribute('data-size', 'form');
    expect(view.container).not.toContainElement(dialog); expect(document.body).toContainElement(dialog);
    expect(dialog.querySelector('.gw-modal-header')).toBeInTheDocument(); expect(dialog.querySelector('.gw-modal-body')).toBeInTheDocument(); expect(dialog.querySelector('.gw-modal-footer')).toBeInTheDocument();
    expect(within(dialog).getByTestId('recording-change-review')).toHaveClass('gw-primary');
    expect(dialog.style.backgroundColor).toBe('var(--gw-raised)'); expect(document.documentElement).toHaveClass(theme);
    const name = within(dialog).getByRole('textbox'); fireEvent.change(name, { target: { value: 'Unsaved recording name' } });
    document.documentElement.className = theme === 'light' ? 'dark' : 'light'; view.rerender(<Change />);
    expect(name).toHaveValue('Unsaved recording name'); expect(dialog).toBeInTheDocument();
  });
  it('short destructive confirm stays explicitly compact, with body scroll and original busy guard', () => {
    const close = vi.fn(), confirm = vi.fn();
    render(<SavedRecordingRemoveDialog open busy onClose={close} onConfirm={confirm} />);
    const dialog = screen.getByRole('alertdialog'); expect(dialog).toHaveAttribute('data-size', 'confirm');
    expect(dialog).toHaveClass('gw-modal'); expect(dialog.querySelector('.gw-modal-body')).toBeInTheDocument();
    expect(screen.getByTestId('group-delete-confirm-button')).toBeDisabled();
    fireEvent.keyDown(dialog, { key: 'Escape' }); expect(close).not.toHaveBeenCalled(); expect(confirm).not.toHaveBeenCalled();
  });
  it('shared CSS has no palette literals and only explicit form classification receives75% sizing', () => {
    const css = readFileSync('src/styles/semantic-modal.css', 'utf8');
    expect(css).not.toMatch(/#[0-9a-f]{3,8}\b|rgba?\(/i); expect(css).toContain('.gw-modal[data-size="form"]');
    expect(css).toContain('width: 75vw'); expect(css).toContain('height: 75dvh');
    expect(css).toContain('.gw-modal[data-size="confirm"]'); expect(css).toContain('height: auto');
    expect(css).toContain('(max-height: 479px)'); expect(css).toContain('safe-area-inset'); expect(css).toContain('prefers-reduced-motion');
    expect(css).toContain(':focus-visible'); expect(css).toContain('var(--gw-focus)');
  });
});
it('global and workbench TS adapters share the same CSS-only semantic authority', async () => {
  const global = await import('@/styles/tokens'); const workbench = await import('@/features/datalink/workbench-v2/tokens');
  expect(workbench.semanticTokens).toBe(global.semanticTokens);
  const variables = [...JSON.stringify(global.semanticTokens).matchAll(/var\((--gw-[a-z-]+)\)/g)].map(match => match[1]);
  expect(variables.length).toBeGreaterThan(30);
  const css = readFileSync('src/index.css', 'utf8');
  for (const variable of variables) expect(css).toContain(`${variable}:`);
  expect(JSON.stringify(global.semanticTokens)).not.toMatch(/#[0-9a-f]{3,8}|rgba?\(/i);
});

it('the actual out-of-modal recovery action also receives44px touch height on narrow or short viewports', () => {
  const css = readFileSync('src/styles/semantic-modal.css', 'utf8');
  const narrow = css.slice(css.indexOf('@media (max-width: 639px)'), css.indexOf('@media (prefers-reduced-motion: reduce)'));
  expect(narrow).toMatch(/\.gw-semantic-recovery button\s*\{\s*min-height:\s*var\(--gw-control-mobile\)/);
});
