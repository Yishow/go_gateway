import test from 'node:test';
import assert from 'node:assert/strict';

import { assertFreshTimingWitness } from './fresh-timing.mjs';

test('fresh timing witness reports ordered segments and preserves polls', () => {
  const polls = [
    { t_poll_open: '1000000000', t_poll_complete: '1500000000', row_count: 0 },
    { t_poll_open: '61000000000', t_poll_complete: '61500000000', row_count: 3 },
  ];
  const result = assertFreshTimingWitness({
    t_open: '1000000000',
    t_start: '3000000000',
    t_closed: '63000000000',
    t_sql: '64500000000',
    polls,
  });

  assert.deepEqual(
    {
      setup_ms: result.setup_ms,
      system_wait_ms: result.system_wait_ms,
      delivery_observed_ms: result.delivery_observed_ms,
      end_to_end_ms: result.end_to_end_ms,
    },
    {
      setup_ms: 2000,
      system_wait_ms: 60000,
      delivery_observed_ms: 1500,
      end_to_end_ms: 63500,
    },
  );
  assert.strictEqual(result.polls, polls);
  assert.match(result.polling_limitation, /poll interval/);
});

test('fresh timing witness rejects a timestamp that moves backwards', () => {
  assert.throws(
    () => assertFreshTimingWitness({
      t_open: '1000000000',
      t_start: '3000000000',
      t_closed: '2000000000',
      t_sql: '4000000000',
      polls: [{ row_count: 1 }],
    }),
    /t_open <= t_start <= t_closed <= t_sql/,
  );
});

test('fresh timing witness rejects first SQL after the 300 second bound', () => {
  assert.throws(
    () => assertFreshTimingWitness({
      t_open: '1000000000',
      t_start: '2000000000',
      t_closed: '3000000000',
      t_sql: '301001000000',
      polls: [{ row_count: 1 }],
    }),
    /at most 300 seconds/,
  );
});

test('fresh timing witness requires polling observations', () => {
  assert.throws(
    () => assertFreshTimingWitness({
      t_open: '1000000000',
      t_start: '2000000000',
      t_closed: '3000000000',
      t_sql: '4000000000',
      polls: [],
    }),
    /polling observations are required/,
  );
});
