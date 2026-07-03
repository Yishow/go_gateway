import { expectTypeOf } from 'vitest';
import type { DbConnector, DbRowGroup, DbTarget, CommitLog, CommitState, WorkbenchV2State } from './types';

/**
 * @file types-step4.test-d.ts
 * @description 驗證 Step 4 Database 相關型別契約，確保符合設計決策與資料模型。
 */

// 1. 斷言 DbConnector 結構與型別
expectTypeOf<DbConnector['kind']>().toEqualTypeOf<'sqlite' | 'postgres' | 'mysql' | 'sqlserver'>();
expectTypeOf<DbConnector['write_mode']>().toEqualTypeOf<'insert' | 'upsert'>();
expectTypeOf<DbConnector['write_interval_seconds']>().toEqualTypeOf<number>();
expectTypeOf<DbConnector['host']>().toEqualTypeOf<string>();

// 2. 斷言 DbTarget 結構與型別
expectTypeOf<DbTarget['tag_id']>().toEqualTypeOf<string>();
expectTypeOf<DbTarget['column_name']>().toEqualTypeOf<string>();
expectTypeOf<DbTarget['enabled']>().toEqualTypeOf<boolean>();
expectTypeOf<DbTarget['row_group_id']>().toEqualTypeOf<string | undefined>();
expectTypeOf<DbRowGroup['member_point_ids']>().toEqualTypeOf<string[]>();

// 3. 斷言 CommitLog 結構與型別
expectTypeOf<CommitLog['label']>().toEqualTypeOf<string>();
expectTypeOf<CommitLog['detail']>().toEqualTypeOf<string>();
expectTypeOf<CommitLog['status']>().toEqualTypeOf<'pending' | 'running' | 'success' | 'failed'>();

// 4. 斷言 CommitState 結構與型別
expectTypeOf<CommitState['status']>().toEqualTypeOf<'idle' | 'committing' | 'success' | 'failed'>();
expectTypeOf<CommitState['logs']>().toEqualTypeOf<CommitLog[]>();

// 5. 斷言 WorkbenchV2State 中包含 db 與 commit 結構
expectTypeOf<WorkbenchV2State['db']>().toEqualTypeOf<{
  connector: DbConnector;
  row_groups?: DbRowGroup[];
  targets: Record<string, DbTarget>;
}>();
expectTypeOf<WorkbenchV2State['committed']>().toEqualTypeOf<boolean>();
