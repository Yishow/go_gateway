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
