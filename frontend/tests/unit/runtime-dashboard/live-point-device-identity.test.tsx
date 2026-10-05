import { render, screen, within } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { LivePointsTable } from '../../../src/features/datalink/runtime-dashboard/LivePointsTable';
import type { RuntimeValueEvent } from '../../../src/types/datalink';
import type { StudioV2RuntimeSetupContext } from '../../../src/types/studioV2RuntimeContext';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string, fallback?: string) => fallback ?? key }) }));
const setup: StudioV2RuntimeSetupContext = { source_rules: [], database_targets: [], mappings: [
  { rule_id: 'rA', device_id: 'A', point_id: 'pA', address: 'D0', display_name: 'A temperature', target_type: 'int16', unit: '°C', enabled: true },
  { rule_id: 'rB', device_id: 'B', point_id: 'pB', address: 'D0', display_name: 'B pressure', target_type: 'uint32', unit: 'kPa', enabled: true },
] };
const event = (device_id: string, point_id: string, value: number | string): RuntimeValueEvent => ({ device_id, point_id, address: 'D0', raw_value: value, transformed_value: value, quality: 'good', stale: false, timestamp: new Date().toISOString() });

describe('device scoped live points', () => {
  it('switches A215/B187 without borrowing same-address values or metadata', () => {
    const a = event('A', 'pA', 215);
    const { rerender } = render(<LivePointsTable selectedDeviceId="A" liveValues={{ pA: a }} setupContext={setup} />);
    expect(screen.getByText('A temperature')).toBeInTheDocument();
    expect(screen.queryByText('B pressure')).not.toBeInTheDocument();
    rerender(<LivePointsTable selectedDeviceId="B" liveValues={{ pA: a }} setupContext={setup} />);
    expect(screen.queryAllByText('215')).toHaveLength(0);
    expect(screen.queryByText('A temperature')).not.toBeInTheDocument();
    rerender(<LivePointsTable selectedDeviceId="B" liveValues={{ pA: a, pB: event('B', 'pB', 187) }} setupContext={setup} />);
    const row = screen.getByText('B pressure').closest('tr')!;
    expect(within(row).getAllByText('187')).toHaveLength(2);
    expect(within(row).getByText('kPa')).toBeInTheDocument();
    expect(within(row).getByText('uint32')).toBeInTheDocument();
    expect(screen.queryAllByText('215')).toHaveLength(0);
  });
  it('displays exact decimal strings without parsing them through a number', () => {
    render(<LivePointsTable selectedDeviceId="B" liveValues={{ pB: event('B', 'pB', '9007199254740993') }} setupContext={setup} />);
    expect(screen.getAllByText('9007199254740993')).toHaveLength(2);
    expect(screen.queryByText('9007199254740992')).not.toBeInTheDocument();
  });
  it('rejects foreign identity even with a matching point key', () => {
    render(<LivePointsTable selectedDeviceId="B" liveValues={{ pB: event('A', 'pB', 215) }} setupContext={setup} />);
    expect(screen.queryAllByText('215')).toHaveLength(0);
  });
});
