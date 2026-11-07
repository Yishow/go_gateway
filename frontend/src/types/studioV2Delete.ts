/** Server receipt for a permanently removed recording lineage. */
export interface WriteGroupDeleteReceipt {
  deleted: true;
  group_id: string;
  lineage_id: string;
  deleted_generations: number;
}

export interface DeleteBacklogDetails {
  samples_unconsumed: number;
  outbox: Record<string, number>;
}

export interface DeviceDeleteReferences {
  recordings: Array<{ id: string; name: string }>;
  total: number;
}

export interface SuccessorDeleteReferences {
  referenced_by: Array<{ id: string; name: string }>;
  total: number;
}
