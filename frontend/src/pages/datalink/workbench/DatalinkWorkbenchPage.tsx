import { useEffect } from 'react';
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

function WorkbenchRouteStateSyncWithRouter() {
  const [searchParams] = useSearchParams();
  const requestedStep = searchParams.get('step');
  const requestedTarget = searchParams.get('target');
  const { activeStep, setActiveStep, activeOutputTarget, setActiveOutputTarget } = useWorkbench();

  useEffect(() => {
    if (isWorkbenchStepParam(requestedStep) && requestedStep !== activeStep) {
      setActiveStep(requestedStep);
    }
  }, [activeStep, requestedStep, setActiveStep]);

  useEffect(() => {
    if (isWorkbenchOutputTargetParam(requestedTarget) && requestedTarget !== activeOutputTarget) {
      setActiveOutputTarget(requestedTarget);
    }
  }, [activeOutputTarget, requestedTarget, setActiveOutputTarget]);

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
