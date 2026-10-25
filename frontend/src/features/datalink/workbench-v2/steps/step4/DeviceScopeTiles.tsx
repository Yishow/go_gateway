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
          className={`rounded border px-3 py-2 text-left text-xs ${selected ? 'border-cyan-500 bg-cyan-950/60 text-cyan-100' : 'border-slate-700 bg-slate-950 text-slate-200'} disabled:cursor-not-allowed disabled:opacity-40`}
        >
          <span className="block font-semibold">{device.name || device.id}</span>
          {device.address && <span className="block text-[11px] text-slate-400">{device.address}</span>}
        </button>
      );
    })}
  </div>
);
