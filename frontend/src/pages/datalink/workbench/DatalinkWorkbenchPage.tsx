import { WorkbenchDeviceStep } from './WorkbenchDeviceStep';
import { WorkbenchFrame } from './WorkbenchFrame';
import { WorkbenchProvider, useWorkbench } from './WorkbenchProvider';
import { LocalModbusBoard } from './LocalModbusBoard';
import { SourceCanvasSection } from './SourceCanvasSection';
import { TagBindingStudio } from './TagBindingStudio';

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

function WorkbenchShell() {
  return (
    <WorkbenchFrame>
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
