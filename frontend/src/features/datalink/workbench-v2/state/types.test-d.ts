import { expectTypeOf } from 'vitest';
import type { WorkbenchV2State } from './types';

// 斷言 view 屬性只能為 'flow' 或 'settings'
expectTypeOf<WorkbenchV2State['view']>().toEqualTypeOf<'flow' | 'settings'>();

// 斷言 current 屬性只能為 1, 2, 3, 4
expectTypeOf<WorkbenchV2State['current']>().toEqualTypeOf<1 | 2 | 3 | 4>();
