import * as React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { describe, expect, it, vi } from 'vitest';
import { fireEvent, render, screen } from '@testing-library/react';
import { ProtocolSelector } from '../../../src/features/datalink/workbench-v2/steps/step1/ProtocolSelector';
import { Step1Device } from '../../../src/features/datalink/workbench-v2/steps/step1/Step1Device';
import { PROTOCOLS } from '../../../src/features/datalink/workbench-v2/state/protocols';
import { getProtocolSupport } from '../../../src/features/datalink/workbench-v2/state/protocolSupport';
import {
  INITIAL_STATE,
  workbenchV2Reducer,
} from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { Device, WorkbenchV2State } from '../../../src/features/datalink/workbench-v2/state/types';
import { AddressParser, getDefaultPlannerStartAddress } from '../../../src/utils/addressParser';
import { inferHydratedProgress } from '../../../src/pages/datalink/workbench-v2/hydratedProgress';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({ t: (key: string) => key }),
}));

const untested: Device = { ...INITIAL_STATE.devices[0], test: null, status: 'draft' };
const tested: Device = {
  ...untested,
  status: 'tested',
  test: { status: 'success', stages: [], total_latency_ms: 5 } as unknown as Device['test'],
};

function Harness({ initial, onContinue }: { initial: WorkbenchV2State; onContinue: () => void }) {
  const [state, dispatch] = React.useReducer(workbenchV2Reducer, initial);
  return (
    <QueryClientProvider client={new QueryClient()}>
      <Step1Device state={state} dispatch={dispatch} onContinue={onContinue} />
    </QueryClientProvider>
  );
}

describe('CapabilityAndReloadReadiness: protocol availability', () => {
  it('derives availability from the address parser, so an option cannot outrun the working path', () => {
    const parser = new AddressParser();
    for (const protocol of PROTOCOLS) {
      const support = getProtocolSupport(protocol.id);
      let parses = true;
      try {
        parser.parse(getDefaultPlannerStartAddress(protocol.id), protocol.id);
      } catch {
        parses = false;
      }
      expect(support.available, protocol.id).toBe(parses);
    }
  });

  it('shows MQTT as unavailable with a reason and never selects it', () => {
    expect(getProtocolSupport('mqtt')).toMatchObject({ available: false, reasonKey: expect.any(String) });
    const onChange = vi.fn();
    render(<ProtocolSelector value="modbus_tcp" onChange={onChange} />);
    const card = screen.getByTestId('protocol-card-mqtt');
    expect(card).toHaveAttribute('aria-disabled', 'true');
    expect(screen.getByTestId('protocol-unavailable-mqtt')).toHaveTextContent(getProtocolSupport('mqtt').reasonKey as string);
    fireEvent.click(card);
    expect(onChange).not.toHaveBeenCalled();
    fireEvent.click(screen.getByTestId('protocol-card-modbus_rtu'));
    expect(onChange).toHaveBeenCalledWith('modbus_rtu');
  });

  it('keeps showing an already saved unavailable protocol instead of silently hiding it', () => {
    render(<ProtocolSelector value="mqtt" onChange={vi.fn()} />);
    expect(screen.getByTestId('protocol-card-mqtt')).toBeInTheDocument();
    expect(screen.getByTestId('protocol-unavailable-mqtt')).toBeInTheDocument();
  });
});

describe('CapabilityAndReloadReadiness: connection identity edits', () => {
  const state: WorkbenchV2State = { ...INITIAL_STATE, devices: [tested] };

  it('invalidates an earlier probe when the connection target changes', () => {
    const next = workbenchV2Reducer(state, { type: 'updateDeviceConfig', deviceId: tested.id, patch: { host: '10.0.0.9' } });
    expect(next.devices[0].test).toBeNull();
    expect(next.devices[0].status).toBe('draft');
  });

  it('keeps the probe when the patch changes nothing', () => {
    const same = workbenchV2Reducer(state, { type: 'updateDeviceConfig', deviceId: tested.id, patch: { host: tested.config.host } });
    expect(same.devices[0].test).not.toBeNull();
  });

  it('keeps the probe when only a request timeout changes', () => {
    const next = workbenchV2Reducer(state, { type: 'updateDeviceConfig', deviceId: tested.id, patch: { timeout: 30 } });
    expect(next.devices[0].test).not.toBeNull();
    expect(next.devices[0].config.timeout).toBe(30);
  });

  it('keeps the probe when only the name or description changes', () => {
    const renamed = workbenchV2Reducer(state, { type: 'renameDevice', deviceId: tested.id, name: 'Line 2' });
    expect(renamed.devices[0].test).not.toBeNull();
  });
});

describe('OfflineDraftNavigation: unverified devices', () => {
  it('lets the operator continue with an unprobed device and marks it unverified', () => {
    const onContinue = vi.fn();
    render(<Harness initial={{ ...INITIAL_STATE, devices: [untested] }} onContinue={onContinue} />);
    const next = screen.getByTestId('btn-continue-step1');
    expect(next).not.toBeDisabled();
    expect(screen.getByTestId('warning-untested-chip')).toBeInTheDocument();
    fireEvent.click(next);
    expect(onContinue).toHaveBeenCalled();
  });

  it('only a missing device blocks continuing', () => {
    render(<Harness initial={{ ...INITIAL_STATE, devices: [] }} onContinue={vi.fn()} />);
    expect(screen.getByTestId('btn-continue-step1')).toBeDisabled();
  });

  it('shows no unverified chip once every device was probed', () => {
    render(<Harness initial={{ ...INITIAL_STATE, devices: [tested] }} onContinue={vi.fn()} />);
    expect(screen.queryByTestId('warning-untested-chip')).not.toBeInTheDocument();
  });

  it('reload and the initial flow agree: an unprobed device still completes step 1 as unverified', () => {
    const progress = inferHydratedProgress([untested], [], {});
    expect(progress.completed.has(1)).toBe(true);
    expect(progress.current).toBe(2);
  });
});
