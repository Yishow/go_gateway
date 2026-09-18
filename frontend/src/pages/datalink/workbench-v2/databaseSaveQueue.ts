export type DatabaseSaveJob<Snapshot> =
  | { kind: 'connector'; snapshot: Snapshot }
  | { kind: 'target'; pointId: string };

export interface DatabaseSaveQueue<Snapshot> {
  enqueueConnector: (snapshot: Snapshot) => void;
  enqueueTarget: (pointId: string) => void;
}

/**
 * 依序執行第 4 步資料庫儲存：每筆要求都帶上一筆回覆的設定版本，避免多筆要求
 * 同時使用同一個版本而互相覆蓋。連線設定優先於已排隊的欄位；連線儲存失敗後，
 * 排隊的欄位等待下一次連線儲存，不以舊設定送出。
 */
export function createDatabaseSaveQueue<Snapshot>(
  run: (job: DatabaseSaveJob<Snapshot>) => Promise<boolean>,
): DatabaseSaveQueue<Snapshot> {
  let running = false;
  let connector: { snapshot: Snapshot } | null = null;
  let targetsBlocked = false;
  const targets: string[] = [];

  const nextJob = (): DatabaseSaveJob<Snapshot> | null => {
    if (connector) {
      const job: DatabaseSaveJob<Snapshot> = { kind: 'connector', snapshot: connector.snapshot };
      connector = null;
      return job;
    }
    if (targetsBlocked) return null;
    const pointId = targets.shift();
    return pointId === undefined ? null : { kind: 'target', pointId };
  };

  const drain = async () => {
    if (running) return;
    running = true;
    try {
      for (let job = nextJob(); job; job = nextJob()) {
        const succeeded = await run(job).catch(() => false);
        if (job.kind === 'connector') targetsBlocked = !succeeded;
      }
    } finally {
      running = false;
    }
  };

  return {
    enqueueConnector(snapshot) {
      connector = { snapshot };
      targetsBlocked = false;
      void drain();
    },
    enqueueTarget(pointId) {
      if (!targets.includes(pointId)) targets.push(pointId);
      void drain();
    },
  };
}
