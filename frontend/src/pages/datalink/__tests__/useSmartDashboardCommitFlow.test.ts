import { act, renderHook } from '@testing-library/react';
import type { TFunction } from 'i18next';
import { describe, expect, it, vi } from 'vitest';
import type { PlannedAllocation } from '../../../components/datalink/MemoryGrid';
import type { Point, PollingGroup } from '../../../types/datalink';
import { useSmartDashboardCommitFlow } from '../smart-dashboard/useSmartDashboardCommitFlow';

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

const pollingGroups: PollingGroup[] = [
  {
    id: 'pg-1',
    name: 'PG1',
    description: '',
    interval_ms: 1000,
    priority: 1,
    enabled: true,
    created_at: '',
    updated_at: '',
  },
];

const plannedAllocations: PlannedAllocation[] = [
  {
    id: 'plan-1',
    dataType: 'int16',
    addresses: ['40001'],
    label: 'SRC-001',
  },
  {
    id: 'plan-2',
    dataType: 'int16',
    addresses: ['40002'],
    label: 'SRC-002',
  },
];

describe('useSmartDashboardCommitFlow', () => {
  it('marks failed chunks and retries after conflict is resolved', () => {
    const markActive = vi.fn();
    const markError = vi.fn();
    const setGuideStage = vi.fn();
    const scheduleGuideStageReset = vi.fn();

    const { result, rerender } = renderHook(
      ({ allPoints }) =>
        useSmartDashboardCommitFlow({
          allPoints,
          linkedAddresses: [],
          plannedAllocations,
          pendingGlobalTagEdit: false,
          pollingGroups,
          canActivate: true,
          mappingEnabled: true,
          markActive,
          markError,
          setGuideStage,
          scheduleGuideStageReset,
          t,
        }),
      {
        initialProps: {
          allPoints: [createPoint({ id: 'conflict-point', address: '40001' })],
        },
      },
    );

    expect(result.current.commitQueueSummary).toMatchObject({
      total: 2,
      conflict: 1,
      pending: 1,
    });

    act(() => {
      result.current.handleCommitFlow();
    });

    expect(markError).toHaveBeenCalled();
    expect(result.current.hasFailedChunk).toBe(true);
    expect(result.current.commitActionMessage).toContain('Commit 部分成功');

    rerender({
      allPoints: [],
    });

    act(() => {
      result.current.handleRetryFailedCommits();
    });

    expect(result.current.hasFailedChunk).toBe(false);
    expect(result.current.commitActionMessage).toContain('Retry 成功');
    expect(markActive).not.toHaveBeenCalled();
    expect(setGuideStage).not.toHaveBeenCalledWith('commit');
  });

  it('commits successfully and supports rollback', () => {
    const markActive = vi.fn();
    const markError = vi.fn();
    const setGuideStage = vi.fn();
    const scheduleGuideStageReset = vi.fn();

    const { result } = renderHook(() =>
      useSmartDashboardCommitFlow({
        allPoints: [],
        linkedAddresses: [],
        plannedAllocations,
        pendingGlobalTagEdit: false,
        pollingGroups,
        canActivate: true,
        mappingEnabled: true,
        markActive,
        markError,
        setGuideStage,
        scheduleGuideStageReset,
        t,
      }),
    );

    act(() => {
      result.current.handleCommitFlow();
    });

    expect(markError).not.toHaveBeenCalled();
    expect(markActive).toHaveBeenCalledTimes(1);
    expect(setGuideStage).toHaveBeenCalledWith('commit');
    expect(scheduleGuideStageReset).toHaveBeenCalledTimes(1);
    expect(result.current.commitActionMessage).toContain('Commit 成功');
    expect(result.current.canRollback).toBe(true);

    act(() => {
      result.current.handleRollbackCommitRun();
    });

    expect(result.current.canRollback).toBe(false);
    expect(result.current.commitActionMessage).toContain('已回滾');
  });
});
