import { act, fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import i18n from '../../../src/i18n/config';
import { MappingRow } from '../../../src/features/datalink/workbench-v2/steps/step3/MappingRow';
import { hydrateStudioV2Mapping } from '../../../src/features/datalink/workbench-v2/state/studioV2MappingAutosave';
import type { Mapping } from '../../../src/features/datalink/workbench-v2/state/types';
import { mappingFixture, mappingPoints } from './mapping-autosave-page.testHarness';

vi.mock('../../../src/features/datalink/workbench-v2/steps/step3/MappingPreviewCells', () => ({ MappingPreviewCells: () => <td /> }));
vi.mock('../../../src/features/datalink/workbench-v2/steps/step3/MappingPayloadDialog', () => ({ MappingPayloadDialog: () => <td /> }));

async function show(state: Mapping['save_state'], advanced = false) {
  await act(async () => { await i18n.changeLanguage('zh-TW'); });
  const point = mappingPoints[0], dispatch = vi.fn(), select = vi.fn();
  const mapping = { ...hydrateStudioV2Mapping(point, mappingFixture()), save_state: state,
    save_error: 'SQL password private', save_error_detail: { code: 'workspace_mapping_save_failed', status: 500, requestId: 'req-safe' } };
  const view = render(<><button>outside</button><table><tbody><MappingRow point={point} mapping={mapping}
    advanced={advanced} devices={[]} isSelected={false} dispatch={dispatch} onSelect={select} /></tbody></table></>);
  return { ...view, point, dispatch, select };
}

describe.each([false, true])('inline mapping save info (advanced=%s)', advanced => {
  it('keeps saved copy hidden until hover, with a green info button beside the name', async () => {
    const { point } = await show('saved', advanced);
    expect(screen.queryByText('已儲存')).not.toBeInTheDocument();
    const info = screen.getByTestId(`mapping-save-info-${point.id}`);
    expect(info).toHaveClass('text-emerald-400');
    expect(info.parentElement).toContainElement(screen.getByTestId(`input-display-name-${point.id}`));
    fireEvent.mouseEnter(info);
    expect(screen.getByRole('dialog')).toHaveTextContent('已儲存');
    fireEvent.keyDown(info, { key: 'Escape' });
    expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
  });

  it('opens safe errors on click, keeps retry functional, and closes outside', async () => {
    const { point, dispatch, select } = await show('save-error', advanced);
    expect(screen.queryByTestId(`mapping-save-retry-${point.id}`)).not.toBeInTheDocument();
    const info = screen.getByTestId(`mapping-save-info-${point.id}`);
    expect(info).toHaveClass('text-rose-400');
    fireEvent.click(info);
    const details = screen.getByRole('dialog');
    expect(details).toHaveTextContent('req-safe');
    expect(details).not.toHaveTextContent('SQL password private');
    fireEvent.click(screen.getByTestId(`mapping-save-retry-${point.id}`));
    expect(dispatch).toHaveBeenCalledExactlyOnceWith({ type: 'retryMappingSave', pointId: point.id });
    expect(select).not.toHaveBeenCalled();
    fireEvent.pointerDown(screen.getByText('outside'));
    expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
  });
});

it.each([['draft-invalid', 'text-rose-400'], ['saving', 'text-amber-400']] as const)('uses a status-specific icon for %s', async (state, color) => {
  const { point } = await show(state);
  const info = screen.getByTestId(`mapping-save-info-${point.id}`);
  expect(info).toHaveClass(color);
  fireEvent.focus(info);
  expect(screen.getByRole('dialog')).toBeInTheDocument();
});

it.each([undefined, 'idle'] as const)('shows no status for %s', async state => {
  await show(state);
  expect(screen.queryByTestId(/^mapping-save-info-/)).not.toBeInTheDocument();
});
