import { describe, expect, it } from 'vitest';
import { workbenchV2Reducer, INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import {
    hydrateStudioV2DatabaseConnector,
    toStudioV2DatabaseConfigRequest,
} from '../../../src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave';

describe('studioV2DatabaseAutosave helpers', () => {
    it('hydrates persisted connector without retaining stale password in reducer state', () => {
        const record = {
            id: 'db-1',
            workspace_id: 'workspace-1',
            kind: 'postgres',
            name: 'Line A PG',
            host: '127.0.0.1',
            port: 5432,
            database: 'gateway',
            username: 'gw_writer',
            schema: 'public',
            table: 'sensor_values',
            write_mode: 'insert',
            write_interval_seconds: 5,
            timestamp_column: 'ts',
            status: 'ready',
            save_state: 'saved',
            runtime_apply_status: 'not_running',
            created_at: '2026-06-07T00:00:00Z',
            updated_at: '2026-06-07T00:00:00Z',
        } as const;

        const hydrated = hydrateStudioV2DatabaseConnector(record);
        const stateWithPassword = {
            ...INITIAL_STATE,
            db: {
                ...INITIAL_STATE.db,
                connector: {
                    ...INITIAL_STATE.db.connector,
                    password: 'stale-secret',
                },
            },
        };

        const nextState = workbenchV2Reducer(stateWithPassword, {
            type: 'updateDbConnector',
            patch: hydrated,
        });

        expect(hydrated.password).toBeUndefined();
        expect(nextState.db.connector.password).toBeUndefined();
    });

    it('omits password when sqlite connector still has stale local password', () => {
        const request = toStudioV2DatabaseConfigRequest({
            ...INITIAL_STATE.db.connector,
            kind: 'sqlite',
            password: 'stale-secret',
        });

        expect(request).not.toHaveProperty('password');
    });
});