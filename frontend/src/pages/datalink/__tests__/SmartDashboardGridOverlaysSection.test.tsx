import { createRef } from 'react';
import { QueryClient } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import type { TFunction } from 'i18next';
import { afterEach, describe, expect, it, vi } from 'vitest';
import SmartDashboardGridOverlaysSection from '../smart-dashboard/SmartDashboardGridOverlaysSection';
import type { SmartDashboardTagPanelProps } from '../smart-dashboard/SmartDashboardTagPanel';
import type { GridContextMenuState } from '../smart-dashboard/useSmartDashboardGridOverlays';
import type { Point } from '../../../types/datalink';

vi.mock('../smart-dashboard/SmartDashboardTagPanel', () => ({
  default: () => <div data-testid="tag-panel-mock" />,
}));

vi.mock('../../../components/datalink/PointDetailPanel', () => ({
  PointDetailPanel: () => <div data-testid="point-detail-panel-mock" />,
}));

const t = ((key: string) => key) as unknown as TFunction;

function createPoint(overrides: Partial<Point> = {}): Point {
  return {
    id: 'point-1',
    device_id: 'device-1',
    name: 'Point 1',
    description: '',
    data_type: 'int16',
    address: '40001',
    enabled: true,
    polling_group_id: 'pg-1',
    last_value: 10,
    last_read_at: '',
    last_error: '',
    error_count: 0,
    created_at: '',
    updated_at: '',
    ...overrides,
  };
}

function createTagPanelProps(): SmartDashboardTagPanelProps {
  return {
    selectedSourceAddress: '',
    activePointForLink: null,
    linkedTag: null,
    tagLinkMode: 'existing',
    setTagLinkMode: vi.fn(),
    selectedTagIdForLink: '',
    setSelectedTagIdForLink: vi.fn(),
    tags: [],
    handleLinkTagToSelectedAddress: vi.fn(),
    updateMappingPending: false,
    createMappingPending: false,
    newTagKey: '',
    setNewTagKey: vi.fn(),
    newTagDisplayName: '',
    setNewTagDisplayName: vi.fn(),
    handleCreateTagAndLink: vi.fn(),
    createTagPending: false,
    tagLinkActionMessage: '',
    linkedTagAffectedMappingsCount: 0,
    tagEditDisplayName: '',
    setTagEditDisplayName: vi.fn(),
    tagEditUnit: '',
    setTagEditUnit: vi.fn(),
    tagEditDescription: '',
    setTagEditDescription: vi.fn(),
    handleSaveLinkedTagEdit: vi.fn(),
    updateTagPending: false,
    pendingTagEdit: null,
    handleConfirmTagEdit: vi.fn(),
    clearPendingTagEdit: vi.fn(),
    tagEditMessage: '',
  };
}

describe('SmartDashboardGridOverlaysSection', () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('renders popover tag panel and supports close action', async () => {
    const setGridPopoverOpen = vi.fn();
    const props = {
      gridPopoverOpen: true,
      setGridPopoverOpen,
      gridPopoverPosition: { top: 120, left: 80 },
      gridPopoverAddress: '40001',
      gridPopoverPoint: null,
      gridPopoverContentRef: createRef<HTMLDivElement>(),
      tagPanelProps: createTagPanelProps(),
      queryClient: new QueryClient(),
      deletePoint: vi.fn().mockResolvedValue(undefined),
      selectedPoint: null,
      setSelectedPoint: vi.fn(),
      showSuccess: vi.fn(),
      showError: vi.fn(),
      gridContextMenu: null as GridContextMenuState | null,
      gridContextMenuRef: createRef<HTMLDivElement>(),
      setGridContextMenu: vi.fn(),
      t,
    };

    render(<SmartDashboardGridOverlaysSection {...props} />);

    expect(await screen.findByRole('dialog', { name: 'smartDashboard.gridCellPopover.ariaLabel' })).toBeInTheDocument();
    expect(screen.getByTestId('tag-panel-mock')).toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'common.close' }));
    expect(setGridPopoverOpen).toHaveBeenCalledWith(false);
  });

  it('deletes single point from context menu and clears menu state', async () => {
    const point = createPoint();
    const deletePoint = vi.fn().mockResolvedValue(undefined);
    const showSuccess = vi.fn();
    const setGridContextMenu = vi.fn();
    vi.spyOn(window, 'confirm').mockReturnValue(true);

    render(
      <SmartDashboardGridOverlaysSection
        gridPopoverOpen={false}
        setGridPopoverOpen={vi.fn()}
        gridPopoverPosition={null}
        gridPopoverAddress=""
        gridPopoverPoint={null}
        gridPopoverContentRef={createRef<HTMLDivElement>()}
        tagPanelProps={createTagPanelProps()}
        queryClient={new QueryClient()}
        deletePoint={deletePoint}
        selectedPoint={point}
        setSelectedPoint={vi.fn()}
        showSuccess={showSuccess}
        showError={vi.fn()}
        gridContextMenu={{
          x: 16,
          y: 20,
          point,
          shiftKey: false,
          pointsToDelete: [point],
        }}
        gridContextMenuRef={createRef<HTMLDivElement>()}
        setGridContextMenu={setGridContextMenu}
        t={t}
      />
    );

    fireEvent.click(screen.getByRole('menuitem', { name: 'smartDashboard.gridContextMenu.deletePoint' }));
    await waitFor(() => {
      expect(deletePoint).toHaveBeenCalledWith(point.id);
    });
    expect(showSuccess).toHaveBeenCalledWith('smartDashboard.pointDeleted');
    expect(setGridContextMenu).toHaveBeenCalledWith(null);
  });
});
