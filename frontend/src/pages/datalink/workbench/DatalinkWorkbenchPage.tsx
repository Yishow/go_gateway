import { useEffect, useState } from 'react';
import { useInRouterContext, useSearchParams } from 'react-router-dom';
import { WorkbenchDeviceStep } from './WorkbenchDeviceStep';
import { WorkbenchFrame } from './WorkbenchFrame';
import { WorkbenchProvider, useWorkbench } from './WorkbenchProvider';
import { LocalModbusBoard } from './LocalModbusBoard';
import { SourceCanvasSection } from './SourceCanvasSection';
import { TagBindingStudio } from './TagBindingStudio';
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
      return <WorkbenchDeviceStep />;
    case 'source':
      return <SourceCanvasSection />;
    case 'tag':
      return <TagBindingStudio />;
    case 'output':
      return <LocalModbusBoard />;
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

function WorkbenchShell() {
  return (
    <WorkbenchFrame>
      <WorkbenchRouteStateSync />
      <StepContent />
    </WorkbenchFrame>
  );
}

export default function DatalinkWorkbenchPage() {
  return (
    <WorkbenchProvider>
      <WorkbenchShell />
    </WorkbenchProvider>
  );
}
