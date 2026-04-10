import { useEffect, useState } from 'react';
import { useInRouterContext, useSearchParams } from 'react-router-dom';
import Box from '@mui/material/Box';
import { MuiDeviceStep } from './MuiDeviceStep';
import { MuiOutputCommandDeck } from './MuiOutputCommandDeck';
import { MuiSourceCommandDeck } from './MuiSourceCommandDeck';
import { MuiTagCommandDeck } from './MuiTagCommandDeck';
import { MuiWorkbenchFrame } from './MuiWorkbenchFrame';
import { MuiWorkbenchShell } from './MuiWorkbenchShell';
import { WorkbenchProvider, useWorkbench } from './WorkbenchProvider';
import { WORKBENCH_STEPS, type OutputTarget, type WorkbenchStep } from './workbenchTypes';

const WORKBENCH_OUTPUT_TARGETS = ['modbus', 'database'] as const;

function isWorkbenchStepParam(value: string | null): value is WorkbenchStep {
  return value !== null && WORKBENCH_STEPS.includes(value as WorkbenchStep);
}

function isWorkbenchOutputTargetParam(value: string | null): value is OutputTarget {
  return value !== null && WORKBENCH_OUTPUT_TARGETS.includes(value as OutputTarget);
}

function StepContent() {
  const { activeStep } = useWorkbench();

  switch (activeStep) {
    case 'device':
      return <MuiDeviceStep />;
    case 'source':
      return <MuiSourceCommandDeck />;
    case 'tag':
      return <MuiTagCommandDeck />;
    case 'output':
      return <MuiOutputCommandDeck />;
    default:
      return null;
  }
}

/**
 * Two-phase route ↔ provider sync:
 *
 * Phase 1 (mount): deep-link query params hydrate provider state.
 * Phase 2 (post-hydration): provider state is the source of truth and
 *   the URL is updated to stay coherent with it.  This prevents stale
 *   query params from overriding user-driven step navigation.
 */
function WorkbenchRouteStateSyncWithRouter() {
  const [searchParams, setSearchParams] = useSearchParams();
  const { activeStep, setActiveStep, activeOutputTarget, setActiveOutputTarget } = useWorkbench();
  const [hydrated, setHydrated] = useState(false);

  // Phase 1: apply deep-link params → provider state (once on mount).
  useEffect(() => {
    const step = searchParams.get('step');
    const target = searchParams.get('target');

    if (isWorkbenchStepParam(step) && step !== activeStep) {
      setActiveStep(step);
    }
    if (isWorkbenchOutputTargetParam(target) && target !== activeOutputTarget) {
      setActiveOutputTarget(target);
    }

    setHydrated(true);
    // eslint-disable-next-line react-hooks/exhaustive-deps -- intentionally run once
  }, []);

  // Phase 2: provider state → URL (after hydration).
  useEffect(() => {
    if (!hydrated) return;

    const next = new URLSearchParams(searchParams);

    if (activeStep === 'device') {
      next.delete('step');
    } else {
      next.set('step', activeStep);
    }

    if (activeOutputTarget === 'modbus') {
      next.delete('target');
    } else {
      next.set('target', activeOutputTarget);
    }

    if (next.toString() !== searchParams.toString()) {
      setSearchParams(next, { replace: true });
    }
  }, [hydrated, activeStep, activeOutputTarget, searchParams, setSearchParams]);

  return null;
}

function WorkbenchRouteStateSync() {
  if (!useInRouterContext()) {
    return null;
  }

  return <WorkbenchRouteStateSyncWithRouter />;
}

/**
 * 將步驟內容包在可伸縮欄位中，配合 {@link WorkbenchFrame} 主欄 `overflow-hidden`，
 * 由各步驟自行決定內部捲動（例如來源步驟僅記憶體格區捲動）。
 */
function WorkbenchShell() {
  return (
    <MuiWorkbenchFrame>
      <WorkbenchRouteStateSync />
      <Box sx={{ display: 'flex', minHeight: 0, minWidth: 0, flex: 1, flexDirection: 'column', overflow: 'hidden' }}>
        <StepContent />
      </Box>
    </MuiWorkbenchFrame>
  );
}

export default function DatalinkWorkbenchPage() {
  return (
    <MuiWorkbenchShell>
      <WorkbenchProvider>
        <WorkbenchShell />
      </WorkbenchProvider>
    </MuiWorkbenchShell>
  );
}
