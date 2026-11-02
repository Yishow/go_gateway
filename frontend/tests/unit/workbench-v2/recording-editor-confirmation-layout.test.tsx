import { useEffect } from 'react';
import { cleanup, fireEvent, render, screen, within } from '@testing-library/react';
import { afterEach, expect, it, vi } from 'vitest';
import { RecordingEditorLayout } from '@/features/datalink/workbench-v2/recording-editor/RecordingEditorLayout';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
afterEach(cleanup);
function Owner({ mounted, unmounted, invoke }: { mounted: () => void; unmounted: () => void; invoke: () => void }) {
  useEffect(() => { mounted(); return unmounted; }, [mounted, unmounted]);
  return <button type="button" onClick={invoke}>Original operation lookup</button>;
}
it.each(['simple', 'guided'] as const)('keeps confirmation outside disabled %s fields without remounting its owner', presentation => {
  const mounted = vi.fn(), unmounted = vi.fn(), invoke = vi.fn();
  const props = { presentation, memberCount: 0, perPoint: false, identity: <input aria-label="Name" />,
    sources: <input aria-label="Source" />, method: <input aria-label="Method" />, destination: <input aria-label="Destination" />,
    policy: <input aria-label="Policy" />, confirmation: <Owner mounted={mounted} unmounted={unmounted} invoke={invoke} />, fieldsDisabled: true };
  const view = render(<RecordingEditorLayout {...props} />);
  for (const name of ['Name', 'Source', 'Method', 'Destination', 'Policy']) expect(screen.getByRole('textbox', { name })).toBeDisabled();
  const command = screen.getByRole('button', { name: 'Original operation lookup' });
  expect(command).toBeEnabled(); expect(command.closest('fieldset')).toBeNull();
  if (presentation === 'simple') expect(within(screen.getByTestId('recording-editor-destination')).getByRole('button', { name: 'Original operation lookup' })).toBe(command);
  fireEvent.click(command); expect(invoke).toHaveBeenCalledOnce();
  view.rerender(<RecordingEditorLayout {...props} fieldsDisabled={false} />);
  for (const name of ['Name', 'Source', 'Method', 'Destination', 'Policy']) expect(screen.getByRole('textbox', { name })).toBeEnabled();
  expect(mounted).toHaveBeenCalledOnce(); expect(unmounted).not.toHaveBeenCalled();
  expect(screen.getByRole('button', { name: 'Original operation lookup' })).toBe(command);
});

it('keeps the simple method controls and destination visible without repeated Next or long section hints', () => {
  const methodChange = vi.fn();
  render(<RecordingEditorLayout presentation="simple" memberCount={2} perPoint identity={<input aria-label="Name" />} sources={<p>Sources</p>}
    method={<select aria-label="Real method" onChange={methodChange}><option value="latest">Latest</option><option value="mean">Mean</option></select>}
    destination={<input aria-label="Destination" />} policy={<input aria-label="Policy" />} />);
  expect(screen.queryByRole('button', { name: /simple.nextMethod|simple.nextDestination/ })).not.toBeInTheDocument();
  expect(screen.queryByText('simple.methodHint')).not.toBeInTheDocument(); expect(screen.queryByText('simple.destinationHint')).not.toBeInTheDocument();
  expect(screen.getByRole('combobox', { name: 'Real method' })).toBeVisible();
  fireEvent.change(screen.getByRole('combobox', { name: 'Real method' }), { target: { value: 'mean' } }); expect(methodChange).toHaveBeenCalledOnce();
  expect(screen.getByRole('textbox', { name: 'Destination' })).toBeVisible(); expect(screen.getByRole('textbox', { name: 'Policy' })).toBeVisible();
});
