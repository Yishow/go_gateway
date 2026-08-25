import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render } from '@testing-library/react';
import { MemoryRouter, Route, Routes, useLocation } from 'react-router-dom';
import { vi } from 'vitest';
import { RuntimeDashboardRoute } from '../../../src/features/datalink/runtime-dashboard/RuntimeDashboardRoute';
import type { Point, RuntimeStatus } from '../../../src/types/datalink';
import type { StudioV2RuntimeSetupContext } from '../../../src/types/studioV2RuntimeContext';
import { createRuntimeSetupFixture } from './runtimeSetupFixture';

const fixture = vi.hoisted(() => ({
  mockRuntimeContext: {
    workspace_id: 'workspace-1',
    devices: [] as Array<{ device_id: string; name: string; protocol: string; running: boolean; availability_status: 'available' | 'unavailable'; availability_reason: string | null }>,
    default_device_id: null as string | null,
    setup: null as StudioV2RuntimeSetupContext | null,
  },
  mockPoints: [] as Point[],
  mockRuntimeStatus: vi.fn<() => Promise<RuntimeStatus>>(),
  mockEventSources: [] as MockEventSource[],
}));
export const { mockRuntimeContext, mockPoints, mockRuntimeStatus, mockEventSources } = fixture;

type EventHandler = (event: MessageEvent<string>) => void;
export class MockEventSource {
  url: string;
  onopen: ((event: Event) => void) | null = null;
  onerror: ((event: Event) => void) | null = null;
  private listeners = new Map<string, Set<EventHandler>>();
  constructor(url: string) { this.url = url; fixture.mockEventSources.push(this); }
  addEventListener(type: string, listener: EventHandler) {
    const existing = this.listeners.get(type) ?? new Set<EventHandler>();
    existing.add(listener); this.listeners.set(type, existing);
  }
  removeEventListener(type: string, listener: EventHandler) { this.listeners.get(type)?.delete(listener); }
  close() {}
  emit(type: string, data: unknown) {
    const payload = { data: JSON.stringify(data) } as MessageEvent<string>;
    this.listeners.get(type)?.forEach((listener) => listener(payload));
  }
  fail() { this.onerror?.(new Event('error')); }
  open() {
    this.onopen?.(new Event('open'));
    this.emit('stream_state', { stream_state: { state: 'ready', empty: false, degraded: false, unavailable: false, stale: false } });
  }
}

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string, fallback?: unknown) => typeof fallback === 'string' ? fallback : key }) }));
vi.mock('../../../src/hooks/datalink/useStudioV2RuntimeContext', () => ({ useStudioV2RuntimeContextQuery: () => ({ data: fixture.mockRuntimeContext, isLoading: false, isError: false, error: null }) }));
vi.mock('../../../src/hooks/datalink/usePoints', () => ({ usePointsQuery: (filters?: { device_id?: string }) => ({ data: filters?.device_id ? fixture.mockPoints.filter((point) => point.device_id === filters.device_id) : [], isLoading: false }) }));
vi.mock('../../../src/services/datalink', () => ({ runtimeAPI: {
  getStatus: fixture.mockRuntimeStatus,
  getStreamUrl: vi.fn((deviceId: string, pointIds?: string[]) => {
    const params = new URLSearchParams({ device_id: deviceId });
    if (pointIds && pointIds.length > 0) params.set('point_ids', pointIds.join(','));
    return `/api/v1/datalink/runtime/stream?${params.toString()}`;
  }),
} }));

function LocationProbe() { return <div data-testid="runtime-dashboard-location">{useLocation().search}</div>; }

export function renderRoute(initialEntry: string) {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(<QueryClientProvider client={queryClient}><MemoryRouter initialEntries={[initialEntry]}><Routes><Route path="/studio/runtime" element={<><RuntimeDashboardRoute /><LocationProbe /></>} /></Routes></MemoryRouter></QueryClientProvider>);
}

export function resetRuntimeDashboardFixture() {
  vi.clearAllMocks(); mockEventSources.splice(0, mockEventSources.length); vi.stubGlobal('EventSource', MockEventSource);
  mockRuntimeContext.workspace_id = 'workspace-1';
  mockRuntimeContext.devices.splice(0, mockRuntimeContext.devices.length, {
    device_id: 'device-B', name: 'Filler PLC', protocol: 'modbus_tcp', running: false, availability_status: 'unavailable', availability_reason: 'invalid Step 1 configuration',
  }, {
    device_id: 'device-A', name: 'Mixer PLC', protocol: 'modbus_tcp', running: true, availability_status: 'available', availability_reason: null,
  });
  mockRuntimeContext.default_device_id = 'device-A'; mockRuntimeContext.setup = createRuntimeSetupFixture();
  mockPoints.splice(0, mockPoints.length, {
    id: 'point-1', device_id: 'device-A', name: 'Flow', description: '', address: '40001', data_type: 'float32', enabled: true, polling_group_id: 'group-A', created_at: '2026-05-29T00:00:00Z', updated_at: '2026-05-29T00:00:00Z', last_read_at: '', last_value: null, last_error: '', error_count: 0,
  } as Point, {
    id: 'point-2', device_id: 'device-B', name: 'Level', description: '', address: '40002', data_type: 'float32', enabled: true, polling_group_id: 'group-B', created_at: '2026-05-29T00:00:00Z', updated_at: '2026-05-29T00:00:00Z', last_read_at: '', last_value: null, last_error: '', error_count: 0,
  } as Point);
  mockRuntimeStatus.mockImplementation(async (deviceId?: string) => ({ running: deviceId !== 'device-B', uptime_seconds: 12, metrics: { collected_total: 10, write_success_total: 8, write_error_total: 0, mapping_error_total: 0, point_state_error_total: 0 }, collectors: [{ device_id: deviceId ?? 'device-A', device_name: deviceId === 'device-B' ? 'Filler PLC' : 'Mixer PLC', protocol: 'modbus_tcp', status: deviceId === 'device-B' ? 'idle' : 'running', availability_status: deviceId === 'device-B' ? 'unavailable' : 'available', availability_reason: deviceId === 'device-B' ? 'invalid Step 1 configuration' : null, running: deviceId !== 'device-B', points_total: 1, points_healthy: deviceId === 'device-B' ? 0 : 1, points_stale: 0, points_error: 0, last_read_at: '2026-05-29T00:00:00Z', last_error: null, breaker_state: 'closed' }] }));
}
