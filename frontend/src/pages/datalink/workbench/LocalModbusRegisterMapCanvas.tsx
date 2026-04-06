import type { ModbusShareMapping } from '../../../types/datalink';
import type { AutoMapStrategy, DryRunResult, WorkbenchOutputCandidate } from './workbenchOutputTypes';
import { toModbusDisplayRegister } from './workbenchOutputTypes';

type RegisterMapConflict = {
  register: number;
  mappings: ModbusShareMapping[];
};

type LocalModbusRegisterMapCanvasProps = {
  candidates: WorkbenchOutputCandidate[];
  shareMappings: ModbusShareMapping[];
  conflicts: RegisterMapConflict[];
  selectedTagId: string;
  viewportAnchorRegister: number | null;
  onSlotClick: (slotIndex: number, occupant: WorkbenchOutputCandidate | null) => void;
  onAutoMap: (strategy: AutoMapStrategy) => void;
  onDryRun: () => void;
  dryRunResults: DryRunResult[] | null;
  t: (key: string, params?: Record<string, unknown>) => string;
};

function dataTypeWidth(dataType: string): number {
  switch (dataType) {
    case 'int32':
    case 'uint32':
    case 'float32':
      return 2;
    case 'int64':
    case 'uint64':
    case 'float64':
      return 4;
    default:
      return 1;
  }
}

function getRegisterSlots(register: number, dataType: string): number[] {
  const width = dataTypeWidth(dataType);
  return Array.from({ length: width }, (_, index) => register + index);
}

export function LocalModbusRegisterMapCanvas({
  candidates,
  shareMappings,
  conflicts,
  selectedTagId,
  viewportAnchorRegister,
  onSlotClick,
  onAutoMap,
  onDryRun,
  dryRunResults,
  t,
}: LocalModbusRegisterMapCanvasProps) {
  const maxCanvasSlots = 256;
  const conflictRegisters = new Set(conflicts.map((conflict) => conflict.register));
  const mappedSlots = new Map<number, WorkbenchOutputCandidate>();

  for (const mapping of shareMappings) {
    const candidate = candidates.find((item) => item.tagId === mapping.tag_id);
    if (!candidate) {
      continue;
    }
    getRegisterSlots(mapping.register, mapping.data_type).forEach((slot) => {
      mappedSlots.set(slot, candidate);
    });
  }

  const maxRegister = Math.max(16, ...Array.from(mappedSlots.keys()).map((register) => register + 4));
  const slotCount = Math.min(maxRegister, maxCanvasSlots);
  const viewportStart = slotCount >= maxRegister
    ? 0
    : Math.max(
      0,
      Math.min(
        (viewportAnchorRegister ?? 0) - Math.floor(slotCount / 2),
        maxRegister - slotCount,
      ),
    );

  return (
    <div data-testid="register-map-canvas" className="space-y-4">
      <div className="flex flex-wrap gap-1">
        {Array.from({ length: slotCount }, (_, index) => {
          const slotRegister = viewportStart + index;
          const candidate = mappedSlots.get(slotRegister);
          const isConflict = conflictRegisters.has(slotRegister);
          const isSlotForSelected = candidate?.tagId === selectedTagId && Boolean(selectedTagId);
          const displayRegister = toModbusDisplayRegister(slotRegister);
          return (
            <button
              type="button"
              key={slotRegister}
              data-testid={`register-slot-${displayRegister}`}
              data-conflict={isConflict ? 'true' : undefined}
              onClick={() => onSlotClick(slotRegister, candidate ?? null)}
              className={`flex min-w-[60px] cursor-pointer flex-col items-center rounded-lg border px-2 py-1 text-[10px] transition hover:ring-1 hover:ring-cyan-400/40 ${
                isConflict
                  ? 'border-amber-500/40 bg-amber-500/10 text-amber-200'
                  : isSlotForSelected
                    ? 'border-cyan-400 bg-cyan-500/20 text-cyan-100 ring-1 ring-cyan-400/50'
                    : candidate
                      ? 'border-cyan-500/30 bg-cyan-500/10 text-cyan-200'
                      : selectedTagId
                        ? 'border-slate-700 bg-slate-900/60 text-slate-400 hover:border-cyan-500/40 hover:bg-cyan-500/5'
                        : 'border-slate-800 bg-slate-950/50 text-slate-500'
              }`}
            >
              <span className="font-mono text-[9px] text-slate-500">HR{displayRegister}</span>
              <span className="truncate font-medium">
                {candidate?.tagKey ?? '—'}
              </span>
            </button>
          );
        })}
      </div>

      <div className="flex flex-wrap gap-2">
        <button
          type="button"
          onClick={() => onAutoMap('sequential')}
          className="rounded-lg border border-slate-700 px-3 py-1.5 text-xs text-slate-200 transition hover:border-cyan-500/40"
        >
          {t('workbench.output.modbusStudio.autoMap.sequential')}
        </button>
        <button
          type="button"
          onClick={() => onAutoMap('gapAware')}
          className="rounded-lg border border-slate-700 px-3 py-1.5 text-xs text-slate-200 transition hover:border-cyan-500/40"
        >
          {t('workbench.output.modbusStudio.autoMap.gapAware')}
        </button>
        <button
          type="button"
          onClick={() => onAutoMap('aligned')}
          className="rounded-lg border border-slate-700 px-3 py-1.5 text-xs text-slate-200 transition hover:border-cyan-500/40"
        >
          {t('workbench.output.modbusStudio.autoMap.aligned')}
        </button>
        <button
          type="button"
          onClick={onDryRun}
          className="rounded-lg border border-emerald-500/30 bg-emerald-500/10 px-3 py-1.5 text-xs font-medium text-emerald-200 transition hover:bg-emerald-500/20"
        >
          {t('workbench.output.modbusStudio.actions.dryRun')}
        </button>
      </div>

      {dryRunResults ? (
        <div data-testid="dry-run-results" className="space-y-2 rounded-xl border border-slate-800 bg-slate-950/60 p-3">
          <p className="text-xs font-semibold uppercase tracking-[0.18em] text-slate-400">
            {t('workbench.output.modbusStudio.dryRun.heading')}
          </p>
          <ul className="space-y-1">
            {dryRunResults.map((result) => (
              <li
                key={result.tagId}
                className={`rounded-lg border px-3 py-1.5 text-xs ${
                  result.valid
                    ? 'border-emerald-500/20 bg-emerald-500/5 text-emerald-200'
                    : 'border-rose-500/20 bg-rose-500/5 text-rose-200'
                }`}
              >
                {result.tagKey}: {result.valid
                  ? t('workbench.output.modbusStudio.dryRun.valid', {
                    register: toModbusDisplayRegister(result.register),
                  })
                  : t('workbench.output.modbusStudio.dryRun.invalid', { reason: result.reason ?? 'unknown' })}
              </li>
            ))}
          </ul>
        </div>
      ) : null}
    </div>
  );
}
