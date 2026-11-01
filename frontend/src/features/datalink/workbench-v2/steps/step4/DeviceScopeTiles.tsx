import * as React from 'react';

export interface DeviceScopeTile {
  id: string;
  name?: string;
  address?: string;
}

export interface DeviceScopeTilesProps {
  devices: DeviceScopeTile[];
  selectedId: string;
  onSelect: (deviceId: string) => void;
  disabled?: boolean;
  label: string;
}

/**
 * The single persisted-device selection of Step 4. Each tile is a native button
 * bound to the device ID, never a dropdown; a shared display name or address
 * never merges two devices. Selecting a tile only changes the viewed scope.
 */
export const DeviceScopeTiles: React.FC<DeviceScopeTilesProps> = ({ devices, selectedId, onSelect, disabled, label }) => (
  <div role="group" aria-label={label} data-testid="device-scope-tiles" className="flex flex-wrap gap-2">
    {devices.map((device) => {
      const selected = device.id === selectedId;
      return (
        <button
          key={device.id}
          type="button"
          data-testid="basic-recording-device-tile"
          data-device-id={device.id}
          aria-pressed={selected}
          disabled={disabled}
          onClick={() => onSelect(device.id)}
          className={`rounded border px-3 py-2 text-left text-xs ${selected ? 'border-[var(--gw-focus)] bg-[var(--gw-active)] text-[var(--gw-accent)]' : 'border-[var(--gw-control)] bg-[var(--gw-surface)] text-[var(--gw-text)]'} disabled:cursor-not-allowed disabled:text-[var(--gw-muted)] disabled:bg-[var(--gw-ground)]`}
        >
          <span className="block font-semibold">{device.name || device.id}</span>
          {device.address && <span className="block text-xs text-[var(--gw-muted)]">{device.address}</span>}
        </button>
      );
    })}
  </div>
);
