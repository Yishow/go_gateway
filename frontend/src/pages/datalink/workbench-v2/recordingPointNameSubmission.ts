import type { SourceIdentity } from './rangeSourceSubmission';
export interface PointNameSession { id: string; identity: SourceIdentity; displayName: string }
export interface PointNameOutcome { session: PointNameSession; status: 'ready' | 'pending' | 'confirmed' | 'unknown' | 'blocked'; currentName?: string; submittedName?: string; reason?: 'source-review'; canDiscardRejected?: boolean }
/** Same mapping queue, not an independent mutation or operation ledger. */
export interface RecordingPointNamePort {
  prepare(identity: SourceIdentity): Promise<PointNameOutcome | undefined>;
  submit(id: string, name: string): Promise<PointNameOutcome | undefined>;
  inspect(id: string): Promise<PointNameOutcome | undefined>;
  release(id: string): void;
  discardRejected?(id: string): Promise<boolean>;
}
