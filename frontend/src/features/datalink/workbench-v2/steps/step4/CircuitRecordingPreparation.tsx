import { CircuitProfileStartStatus } from './CircuitProfileBasicActions';
import { GuidedRecordingPreparation } from './GuidedRecordingPreparation';
import { CircuitProfileSchemaPreparation, type CircuitProfileSchemaPreparationProps } from './CircuitProfileSchemaPreparation';
import type { GuidedRecordingPreparationProps } from './recordingPreparationPresentation';
import type { useCircuitProfileStart } from './useCircuitProfileStart';
import { recordingStartObservation } from './useRecordingPreparationOwner';

export function CircuitRecordingPreparation({ flow, schema, onContinueToRuntime, deviceNames, groupNames, ...props }: Omit<GuidedRecordingPreparationProps,
  'start' | 'renderPreparation' | 'testIdPrefix'> & {
  flow: ReturnType<typeof useCircuitProfileStart>; schema: CircuitProfileSchemaPreparationProps;
  onContinueToRuntime?: (deviceIds: string[]) => void; deviceNames?: Record<string, string>; groupNames?: Record<string, string>;
}) {
  return <GuidedRecordingPreparation {...props} testIdPrefix="circuit-profile"
    completedOperationId={flow.completion.state === 'available' && flow.completion.handoffDeviceId ? flow.completion.operationId : undefined}
    renderStartResult={props.prepareOnly ? undefined : primary => <CircuitProfileStartStatus flow={flow} primary={primary}
      onContinueToRuntime={onContinueToRuntime} deviceNames={deviceNames} groupNames={groupNames} />}
    start={{ observation: recordingStartObservation(props.ownerKey, flow), busy: flow.busy,
      prepare: flow.prepare, start: () => flow.start(false), check: flow.check }}
    renderPreparation={bridge => <CircuitProfileSchemaPreparation {...schema} guided={bridge} />} />;
}
