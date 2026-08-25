import { expectTypeOf } from 'vitest';
import type {
  ModbusShareDesiredMapping as DedicatedDesiredMapping,
  ModbusShareReconcileOutcome as DedicatedReconcileOutcome,
  ModbusShareStatus as DedicatedStatus,
  TypedAPIError as DedicatedTypedAPIError,
} from './modbusShare';
import type {
  ModbusShareDesiredMapping,
  ModbusShareReconcileOutcome,
  ModbusShareStatus,
  TypedAPIError,
} from './datalink';

expectTypeOf<ModbusShareStatus>().toEqualTypeOf<DedicatedStatus>();
expectTypeOf<ModbusShareDesiredMapping>().toEqualTypeOf<DedicatedDesiredMapping>();
expectTypeOf<ModbusShareReconcileOutcome>().toEqualTypeOf<DedicatedReconcileOutcome>();
expectTypeOf<TypedAPIError>().toEqualTypeOf<DedicatedTypedAPIError>();
