import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { WorkbenchProvider, useWorkbench } from '@/pages/datalink/workbench/WorkbenchProvider';

function WorkbenchProbe() {
  const {
    activeStep,
    selectedDeviceId,
    steps,
    inspectorSelection,
    activeOutputTarget,
    outputSelectionState,
    crossStepContext,
    sourcePlanningState,
    setActiveStep,
    setSelectedDeviceId,
    setInspectorSelection,
    clearInspectorSelection,
    setActiveOutputTarget,
    setSelectedOutputTagId,
    setFocusedRuleId,
    setFocusedTagIds,
    setSourcePlannerStartAddress,
    clearOutputSelectionState,
    clearCrossStepContext,
  } = useWorkbench();

  return (
    <div>
      <p data-testid="active-step">{activeStep}</p>
      <p data-testid="selected-device">{selectedDeviceId ?? 'none'}</p>
      <p data-testid="steps">{steps.join(',')}</p>
      <p data-testid="inspector-kind">{inspectorSelection.kind}</p>
      <p data-testid="output-target">{activeOutputTarget}</p>
      <p data-testid="output-tag-modbus">{outputSelectionState.modbus || 'none'}</p>
      <p data-testid="output-tag-database">{outputSelectionState.database || 'none'}</p>
      <p data-testid="focused-rule-id">{crossStepContext.focusedRuleId ?? 'none'}</p>
      <p data-testid="focused-tag-ids">{crossStepContext.focusedTagIds.join(',') || 'none'}</p>
      <p data-testid="device-42-start-address">
        {sourcePlanningState.plannerStartAddressByDeviceId['device-42'] ?? 'none'}
      </p>
      <p data-testid="device-99-start-address">
        {sourcePlanningState.plannerStartAddressByDeviceId['device-99'] ?? 'none'}
      </p>
      <button type="button" onClick={() => setActiveStep('output')}>
        go-output
      </button>
      <button type="button" onClick={() => setActiveStep('source')}>
        go-source
      </button>
      <button type="button" onClick={() => setActiveStep('tag')}>
        go-tag
      </button>
      <button type="button" onClick={() => setSelectedDeviceId('device-42')}>
        select-device
      </button>
      <button type="button" onClick={() => setSelectedDeviceId('device-99')}>
        select-device-99
      </button>
      <button
        type="button"
        onClick={() => setInspectorSelection({ kind: 'device', deviceId: 'dev-1' })}
      >
        select-inspector-device
      </button>
      <button
        type="button"
        onClick={() => setInspectorSelection({ kind: 'tag', tagId: 'tag-1', pointId: 'pt-1' })}
      >
        select-inspector-tag
      </button>
      <button
        type="button"
        onClick={() =>
          setInspectorSelection({ kind: 'outputCandidate', tagId: 'tag-2', target: 'database' })
        }
      >
        select-inspector-output
      </button>
      <button type="button" onClick={clearInspectorSelection}>
        clear-inspector
      </button>
      <button type="button" onClick={() => setActiveOutputTarget('database')}>
        switch-to-database
      </button>
      <button type="button" onClick={() => setSelectedOutputTagId('modbus', 'tag-modbus')}>
        select-modbus-tag
      </button>
      <button type="button" onClick={() => setSelectedOutputTagId('database', 'tag-database')}>
        select-database-tag
      </button>
      <button type="button" onClick={clearOutputSelectionState}>
        clear-output-selection
      </button>
      <button type="button" onClick={() => setFocusedRuleId('rule-abc')}>
        focus-rule
      </button>
      <button type="button" onClick={() => setFocusedRuleId(null)}>
        clear-focused-rule
      </button>
      <button type="button" onClick={() => setFocusedTagIds(['tag-1', 'tag-2'])}>
        focus-tags
      </button>
      <button
        type="button"
        onClick={() => setSourcePlannerStartAddress('device-42', '40011')}
      >
        set-device-42-start
      </button>
      <button
        type="button"
        onClick={() => setSourcePlannerStartAddress('device-99', 'D20')}
      >
        set-device-99-start
      </button>
      <button type="button" onClick={clearCrossStepContext}>
        clear-cross-step
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

  describe('inspector selection', () => {
    it('defaults to no selection', () => {
      render(
        <WorkbenchProvider>
          <WorkbenchProbe />
        </WorkbenchProvider>,
      );

      expect(screen.getByTestId('inspector-kind')).toHaveTextContent('none');
    });

    it('updates inspector selection for device kind', () => {
      render(
        <WorkbenchProvider>
          <WorkbenchProbe />
        </WorkbenchProvider>,
      );

      fireEvent.click(screen.getByRole('button', { name: 'select-inspector-device' }));
      expect(screen.getByTestId('inspector-kind')).toHaveTextContent('device');
    });

    it('updates inspector selection for tag kind with pointId', () => {
      render(
        <WorkbenchProvider>
          <WorkbenchProbe />
        </WorkbenchProvider>,
      );

      fireEvent.click(screen.getByRole('button', { name: 'select-inspector-tag' }));
      expect(screen.getByTestId('inspector-kind')).toHaveTextContent('tag');
    });

    it('updates inspector selection for outputCandidate kind with target', () => {
      render(
        <WorkbenchProvider>
          <WorkbenchProbe />
        </WorkbenchProvider>,
      );

      fireEvent.click(screen.getByRole('button', { name: 'select-inspector-output' }));
      expect(screen.getByTestId('inspector-kind')).toHaveTextContent('outputCandidate');
    });

    it('clears inspector selection', () => {
      render(
        <WorkbenchProvider>
          <WorkbenchProbe />
        </WorkbenchProvider>,
      );

      fireEvent.click(screen.getByRole('button', { name: 'select-inspector-device' }));
      expect(screen.getByTestId('inspector-kind')).toHaveTextContent('device');

      fireEvent.click(screen.getByRole('button', { name: 'clear-inspector' }));
      expect(screen.getByTestId('inspector-kind')).toHaveTextContent('none');
    });

    it('clears inspector selection when switching steps', () => {
      render(
        <WorkbenchProvider>
          <WorkbenchProbe />
        </WorkbenchProvider>,
      );

      fireEvent.click(screen.getByRole('button', { name: 'select-inspector-device' }));
      expect(screen.getByTestId('inspector-kind')).toHaveTextContent('device');

      fireEvent.click(screen.getByRole('button', { name: 'go-source' }));
      expect(screen.getByTestId('inspector-kind')).toHaveTextContent('none');
      expect(screen.getByTestId('active-step')).toHaveTextContent('source');
    });
  });

  describe('output target', () => {
    it('defaults to modbus', () => {
      render(
        <WorkbenchProvider>
          <WorkbenchProbe />
        </WorkbenchProvider>,
      );

      expect(screen.getByTestId('output-target')).toHaveTextContent('modbus');
    });

    it('switches to database target', () => {
      render(
        <WorkbenchProvider>
          <WorkbenchProbe />
        </WorkbenchProvider>,
      );

      fireEvent.click(screen.getByRole('button', { name: 'switch-to-database' }));
      expect(screen.getByTestId('output-target')).toHaveTextContent('database');
    });

    it('stores selected tag separately for each output target', () => {
      render(
        <WorkbenchProvider>
          <WorkbenchProbe />
        </WorkbenchProvider>,
      );

      fireEvent.click(screen.getByRole('button', { name: 'select-modbus-tag' }));
      fireEvent.click(screen.getByRole('button', { name: 'select-database-tag' }));

      expect(screen.getByTestId('output-tag-modbus')).toHaveTextContent('tag-modbus');
      expect(screen.getByTestId('output-tag-database')).toHaveTextContent('tag-database');
    });

    it('clears output selections when the device changes', () => {
      render(
        <WorkbenchProvider>
          <WorkbenchProbe />
        </WorkbenchProvider>,
      );

      fireEvent.click(screen.getByRole('button', { name: 'select-modbus-tag' }));
      fireEvent.click(screen.getByRole('button', { name: 'select-database-tag' }));
      expect(screen.getByTestId('output-tag-modbus')).toHaveTextContent('tag-modbus');
      expect(screen.getByTestId('output-tag-database')).toHaveTextContent('tag-database');

      fireEvent.click(screen.getByRole('button', { name: 'select-device' }));

      expect(screen.getByTestId('output-tag-modbus')).toHaveTextContent('none');
      expect(screen.getByTestId('output-tag-database')).toHaveTextContent('none');
    });
  });

  describe('cross-step traceability context', () => {
    it('defaults to empty cross-step context', () => {
      render(
        <WorkbenchProvider>
          <WorkbenchProbe />
        </WorkbenchProvider>,
      );

      expect(screen.getByTestId('focused-rule-id')).toHaveTextContent('none');
      expect(screen.getByTestId('focused-tag-ids')).toHaveTextContent('none');
    });

    it('updates focusedRuleId via setFocusedRuleId', () => {
      render(
        <WorkbenchProvider>
          <WorkbenchProbe />
        </WorkbenchProvider>,
      );

      fireEvent.click(screen.getByRole('button', { name: 'focus-rule' }));
      expect(screen.getByTestId('focused-rule-id')).toHaveTextContent('rule-abc');
    });

    it('updates focusedTagIds via setFocusedTagIds', () => {
      render(
        <WorkbenchProvider>
          <WorkbenchProbe />
        </WorkbenchProvider>,
      );

      fireEvent.click(screen.getByRole('button', { name: 'focus-tags' }));
      expect(screen.getByTestId('focused-tag-ids')).toHaveTextContent('tag-1,tag-2');
    });

    it('clears focusedRuleId via setFocusedRuleId(null)', () => {
      render(
        <WorkbenchProvider>
          <WorkbenchProbe />
        </WorkbenchProvider>,
      );

      fireEvent.click(screen.getByRole('button', { name: 'focus-rule' }));
      expect(screen.getByTestId('focused-rule-id')).toHaveTextContent('rule-abc');

      fireEvent.click(screen.getByRole('button', { name: 'clear-focused-rule' }));
      expect(screen.getByTestId('focused-rule-id')).toHaveTextContent('none');
    });

    it('clears entire cross-step context via clearCrossStepContext', () => {
      render(
        <WorkbenchProvider>
          <WorkbenchProbe />
        </WorkbenchProvider>,
      );

      fireEvent.click(screen.getByRole('button', { name: 'focus-rule' }));
      fireEvent.click(screen.getByRole('button', { name: 'focus-tags' }));
      expect(screen.getByTestId('focused-rule-id')).toHaveTextContent('rule-abc');
      expect(screen.getByTestId('focused-tag-ids')).toHaveTextContent('tag-1,tag-2');

      fireEvent.click(screen.getByRole('button', { name: 'clear-cross-step' }));
      expect(screen.getByTestId('focused-rule-id')).toHaveTextContent('none');
      expect(screen.getByTestId('focused-tag-ids')).toHaveTextContent('none');
    });

    it('PERSISTS cross-step context when switching steps (unlike inspectorSelection)', () => {
      render(
        <WorkbenchProvider>
          <WorkbenchProbe />
        </WorkbenchProvider>,
      );

      // Set a rule context in step 2 (source)
      fireEvent.click(screen.getByRole('button', { name: 'go-source' }));
      fireEvent.click(screen.getByRole('button', { name: 'focus-rule' }));
      expect(screen.getByTestId('focused-rule-id')).toHaveTextContent('rule-abc');
      expect(screen.getByTestId('inspector-kind')).toHaveTextContent('none'); // was cleared by step switch

      // Navigate to step 3 (tag) — cross-step context should remain
      fireEvent.click(screen.getByRole('button', { name: 'go-tag' }));
      expect(screen.getByTestId('active-step')).toHaveTextContent('tag');
      expect(screen.getByTestId('focused-rule-id')).toHaveTextContent('rule-abc');
      expect(screen.getByTestId('inspector-kind')).toHaveTextContent('none');
    });

    it('clears cross-step context when a new device is selected', () => {
      render(
        <WorkbenchProvider>
          <WorkbenchProbe />
        </WorkbenchProvider>,
      );

      fireEvent.click(screen.getByRole('button', { name: 'select-device' }));
      fireEvent.click(screen.getByRole('button', { name: 'focus-rule' }));
      fireEvent.click(screen.getByRole('button', { name: 'focus-tags' }));
      expect(screen.getByTestId('focused-rule-id')).toHaveTextContent('rule-abc');
      expect(screen.getByTestId('focused-tag-ids')).toHaveTextContent('tag-1,tag-2');

      // Selecting a different device resets the source context
      fireEvent.click(screen.getByRole('button', { name: 'select-device-99' }));
      expect(screen.getByTestId('selected-device')).toHaveTextContent('device-99');
      expect(screen.getByTestId('focused-rule-id')).toHaveTextContent('none');
      expect(screen.getByTestId('focused-tag-ids')).toHaveTextContent('none');
    });

    it('preserves remembered planner start addresses when switching devices', () => {
      render(
        <WorkbenchProvider>
          <WorkbenchProbe />
        </WorkbenchProvider>,
      );

      fireEvent.click(screen.getByRole('button', { name: 'set-device-42-start' }));
      fireEvent.click(screen.getByRole('button', { name: 'set-device-99-start' }));
      fireEvent.click(screen.getByRole('button', { name: 'select-device' }));
      fireEvent.click(screen.getByRole('button', { name: 'select-device-99' }));

      expect(screen.getByTestId('device-42-start-address')).toHaveTextContent('40011');
      expect(screen.getByTestId('device-99-start-address')).toHaveTextContent('D20');
    });
  });
});
