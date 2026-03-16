import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { WorkbenchProvider, useWorkbench } from '../WorkbenchProvider';

function WorkbenchProbe() {
  const {
    activeStep,
    selectedDeviceId,
    steps,
    setActiveStep,
    setSelectedDeviceId,
  } = useWorkbench();

  return (
    <div>
      <p data-testid="active-step">{activeStep}</p>
      <p data-testid="selected-device">{selectedDeviceId ?? 'none'}</p>
      <p data-testid="steps">{steps.join(',')}</p>
      <button type="button" onClick={() => setActiveStep('output')}>
        go-output
      </button>
      <button type="button" onClick={() => setSelectedDeviceId('device-42')}>
        select-device
      </button>
    </div>
  );
}

describe('WorkbenchProvider', () => {
  it('provides ordered steps and defaults to the device step', () => {
    render(
      <WorkbenchProvider>
        <WorkbenchProbe />
      </WorkbenchProvider>,
    );

    expect(screen.getByTestId('active-step')).toHaveTextContent('device');
    expect(screen.getByTestId('selected-device')).toHaveTextContent('none');
    expect(screen.getByTestId('steps')).toHaveTextContent('device,source,tag,output');
  });

  it('updates the active step and selected device id', () => {
    render(
      <WorkbenchProvider>
        <WorkbenchProbe />
      </WorkbenchProvider>,
    );

    fireEvent.click(screen.getByRole('button', { name: 'go-output' }));
    fireEvent.click(screen.getByRole('button', { name: 'select-device' }));

    expect(screen.getByTestId('active-step')).toHaveTextContent('output');
    expect(screen.getByTestId('selected-device')).toHaveTextContent('device-42');
  });
});
