import {
  createContext,
  useContext,
  useMemo,
  useState,
  type ReactNode,
} from 'react';
import { WORKBENCH_STEPS, type WorkbenchStep } from './workbenchTypes';

type WorkbenchContextValue = {
  steps: readonly WorkbenchStep[];
  activeStep: WorkbenchStep;
  selectedDeviceId: string | null;
  setActiveStep: (step: WorkbenchStep) => void;
  setSelectedDeviceId: (deviceId: string | null) => void;
};

const WorkbenchContext = createContext<WorkbenchContextValue | null>(null);

export function WorkbenchProvider({ children }: { children: ReactNode }) {
  const [activeStep, setActiveStep] = useState<WorkbenchStep>('device');
  const [selectedDeviceId, setSelectedDeviceId] = useState<string | null>(null);

  const value = useMemo<WorkbenchContextValue>(
    () => ({
      steps: WORKBENCH_STEPS,
      activeStep,
      selectedDeviceId,
      setActiveStep,
      setSelectedDeviceId,
    }),
    [activeStep, selectedDeviceId],
  );

  return <WorkbenchContext.Provider value={value}>{children}</WorkbenchContext.Provider>;
}

export function useWorkbench() {
  const context = useContext(WorkbenchContext);

  if (!context) {
    throw new Error('useWorkbench must be used within WorkbenchProvider');
  }

  return context;
}
