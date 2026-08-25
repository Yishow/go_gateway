import { expectTypeOf } from 'vitest';
import { modbusShareAPI } from './modbusShare';
import { modbusShareAPI as legacyModbusShareAPI } from './datalink';

expectTypeOf(legacyModbusShareAPI).toEqualTypeOf<typeof modbusShareAPI>();
expectTypeOf(legacyModbusShareAPI.reconcile).toEqualTypeOf(modbusShareAPI.reconcile);
