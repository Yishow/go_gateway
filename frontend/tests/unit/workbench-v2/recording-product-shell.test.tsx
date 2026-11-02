import { afterEach, expect, it, vi } from 'vitest';
import { cleanup, fireEvent, render, screen, within } from '@testing-library/react';
import { RecordingHome } from '@/features/datalink/workbench-v2/recording-home/RecordingHome';
import { SimpleRecordingShell } from '@/features/datalink/workbench-v2/shell/SimpleRecordingShell';
import { recordingHomeGuidance } from '@/features/datalink/workbench-v2/recording-home/recordingHomeGuidance';
import { fixture } from './recording-home-guidance-fixtures';
const theme = vi.hoisted(() => ({ preference: 'system', saveState: 'saved', setPreference: vi.fn(), persistCurrent: vi.fn() }));
vi.mock('@/contexts/ThemeContext', () => ({ useOptionalTheme: () => theme }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
afterEach(() => { cleanup(); vi.clearAllMocks(); });
it('gives Home product identity and theme controls through the existing owner without navigation', () => {
 const navigate = vi.fn(); render(<RecordingHome guidance={recordingHomeGuidance(fixture(['D1']))} location={{ search: '?view=recordings', hash: '' }} onRetry={vi.fn()} onNavigate={navigate} />);
 const shell = screen.getByTestId('recording-product-header'); expect(within(shell).getByText('go_gateway')).toBeInTheDocument();
 expect(within(shell).getByText('product.currentWorkspace')).toHaveAttribute('title', 'W1');
 expect(within(shell).queryByText('W1')).not.toBeInTheDocument();
 fireEvent.change(within(shell).getByLabelText('product.appearance'), { target: { value: 'dark' } });
 expect(theme.setPreference).toHaveBeenCalledExactlyOnceWith('dark'); expect(navigate).not.toHaveBeenCalled();
 expect(within(shell).getByText('product.tools').closest('details')).not.toHaveAttribute('open');
});
it('uses the same quiet product header for simple editing without replacing its form owner', () => {
 const back = vi.fn(); render(<SimpleRecordingShell onBack={back} onAdvanced={vi.fn()}><input aria-label="draft" defaultValue="kept" /></SimpleRecordingShell>);
 expect(within(screen.getByTestId('recording-product-header')).getByText('go_gateway')).toBeInTheDocument();
 expect(within(screen.getByTestId('recording-product-header')).getByText('unknown')).not.toHaveAttribute('title');
 expect(screen.getByRole('heading', { level: 1, name: 'simple.create' })).toBeVisible();
 expect(screen.getByLabelText('draft')).toHaveValue('kept'); fireEvent.click(screen.getByText('simple.back')); expect(back).toHaveBeenCalledOnce();
});

it('names located editing and new copies without replacing the mounted draft or navigation owners', () => {
 const back = vi.fn(), advanced = vi.fn();
 const view = render(<SimpleRecordingShell mode="located" onBack={back} onAdvanced={advanced}><input aria-label="kept draft" defaultValue="original" /></SimpleRecordingShell>);
 const draft = screen.getByLabelText('kept draft'); fireEvent.change(draft, { target: { value: 'unsent' } });
 expect(screen.getByRole('heading', { level: 1, name: 'workflow.editTitle' })).toBeVisible();
 view.rerender(<SimpleRecordingShell mode="copy" onBack={back} onAdvanced={advanced}><input aria-label="kept draft" defaultValue="original" /></SimpleRecordingShell>);
 expect(screen.getByRole('heading', { level: 1, name: 'simple.create' })).toBeVisible(); expect(screen.getByLabelText('kept draft')).toBe(draft); expect(draft).toHaveValue('unsent');
 expect(back).not.toHaveBeenCalled(); expect(advanced).not.toHaveBeenCalled();
});
