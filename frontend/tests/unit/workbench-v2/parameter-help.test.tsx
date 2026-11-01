import * as Dialog from '@radix-ui/react-dialog';
import { act, cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, expect, it, vi } from 'vitest';
import { ParameterHelp } from '@/features/datalink/workbench-v2/components/ParameterHelp';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string, values?: { title?: string }) => key === 'open' ? `About ${values?.title}` : key }) }));
afterEach(cleanup);
it('opens field-specific help by pointer without submitting or editing the form', () => {
 const submit = vi.fn(), change = vi.fn();
 render(<form onSubmit={submit}><input aria-label="Value" defaultValue="draft" onChange={change} /><ParameterHelp help="device.name" /></form>);
 const button = screen.getByRole('button', { name: 'About device.name.title' }); expect(button).toHaveAttribute('type', 'button');
 fireEvent.click(button); expect(button).toHaveAttribute('aria-expanded', 'true');
 const help = screen.getByRole('note'); expect(help).toHaveTextContent('device.name.what'); expect(help).toHaveTextContent('device.name.how'); expect(help).toHaveTextContent('device.name.example');
 expect(help).toHaveAttribute('popover', 'manual'); expect(submit).not.toHaveBeenCalled(); expect(change).not.toHaveBeenCalled(); expect(screen.getByLabelText('Value')).toHaveValue('draft');
});
it('accepts native keyboard-generated click and Escape returns to the same trigger', () => {
 render(<ParameterHelp help="device.name" />);
 const button = screen.getByRole('button', { name: 'About device.name.title' }); button.focus(); fireEvent.click(button, { detail: 0 });
 expect(screen.getByRole('note')).toBeVisible(); fireEvent.keyDown(window, { key: 'Escape' }); expect(screen.queryByRole('note')).not.toBeInTheDocument(); expect(button).toHaveFocus();
});
it('outside pointer dismissal does not steal the newly chosen focus', () => {
 render(<><ParameterHelp help="device.name" /><input aria-label="Next field" /></>);
 fireEvent.click(screen.getByRole('button', { name: 'About device.name.title' })); const next = screen.getByLabelText('Next field'); next.focus(); fireEvent.pointerDown(next);
 expect(screen.queryByRole('note')).not.toBeInTheDocument(); expect(screen.getByLabelText('Next field')).toHaveFocus();
});
it('Escape closes help before its existing modal, without creating a second dialog', () => {
 const close = vi.fn();
 render(<Dialog.Root open onOpenChange={close}><Dialog.Content aria-describedby={undefined}><Dialog.Title>Native form</Dialog.Title><ParameterHelp help="device.name" /></Dialog.Content></Dialog.Root>);
 const trigger = screen.getByRole('button', { name: 'About device.name.title' }); fireEvent.click(trigger);
 expect(screen.getAllByRole('dialog')).toHaveLength(1); expect(screen.getByRole('dialog')).toContainElement(screen.getByRole('note'));
 fireEvent.keyDown(window, { key: 'Escape' }); expect(close).not.toHaveBeenCalled(); expect(trigger).toHaveFocus(); expect(screen.queryByRole('note')).not.toBeInTheDocument();
});
it('IME Escape is not interpreted as a help close', () => {
 render(<ParameterHelp help="device.name" />); fireEvent.click(screen.getByRole('button', { name: 'About device.name.title' }));
 fireEvent.keyDown(window, { key: 'Escape', isComposing: true }); expect(screen.getByRole('note')).toBeVisible();
});

it('Field keeps help outside its native label and preserves its control name without a mutation', async () => {
 const { Field } = await import('@/features/datalink/workbench-v2/components/Field'); const change = vi.fn();
 render(<Field label="Device name" help="device.name"><input defaultValue="Unsent" onChange={change} /></Field>);
 const input = screen.getByLabelText('Device name'), button = screen.getByRole('button', { name: 'About device.name.title' });
 expect(button.closest('label')).toBeNull(); fireEvent.click(button); expect(input).toHaveValue('Unsent'); expect(change).not.toHaveBeenCalled();
});
it('disabled configuration still has keyboard-reachable help and inline validation', async () => {
 const { Field } = await import('@/features/datalink/workbench-v2/components/Field');
 render(<Field label="Device name" help="device.name" error="Source is unavailable"><input disabled defaultValue="Saved" /></Field>);
 expect(screen.getByLabelText(/Device name/)).toBeDisabled(); expect(screen.getByRole('button', { name: 'About device.name.title' })).toBeEnabled();
 expect(screen.getByText('Source is unavailable')).toBeVisible();
});

it('keeps the circle touch target square despite the existing modal control padding', () => {
 render(<ParameterHelp help="device.name" />);
 expect(screen.getByRole('button', { name: 'About device.name.title' })).toHaveStyle({ width: '44px', height: '44px', minWidth: '44px', padding: '0px', borderRadius: '50%' });
});
it('fits its actual width inside the offset visual viewport and recomputes on resize', () => {
 const viewport = Object.assign(new EventTarget(), { width: 300, height: 400, offsetLeft: 40, offsetTop: 30 });
 vi.stubGlobal('visualViewport', viewport);
 try {
  render(<ParameterHelp help="device.name" />);
  fireEvent.click(screen.getByRole('button', { name: 'About device.name.title' }));
  const note = screen.getByRole('note');
  expect(note).toHaveStyle({ width: '284px', left: '48px' });
  viewport.width = 180; viewport.offsetLeft = 80;
  act(() => { viewport.dispatchEvent(new Event('resize')); });
  expect(note).toHaveStyle({ width: '164px', left: '88px' });
  expect(parseFloat(note.style.left) + parseFloat(note.style.width)).toBeLessThanOrEqual(80 + 180 - 8);
 } finally { vi.unstubAllGlobals(); }
});
