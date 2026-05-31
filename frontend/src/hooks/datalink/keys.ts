/**
 * Datalink Query Key Factories
 *
 * 結構化的 Query Key 設計，支援精準的 Cache Invalidation。
 * 遵循 TanStack Query 最佳實踐：階層式 Key 結構。
 */

/** Device 相關 Query Keys */
export const deviceKeys = {
  /** 所有 Device 相關快取的根 Key */
  all: ['devices'] as const,

  /** Device 列表查詢 */
  lists: () => [...deviceKeys.all, 'list'] as const,

  /** 帶篩選條件的 Device 列表 */
  list: (filters?: { protocol?: string; status?: string }) =>
    [...deviceKeys.lists(), filters ?? {}] as const,

  /** Device 詳情查詢 */
  details: () => [...deviceKeys.all, 'detail'] as const,

  /** 特定 Device 詳情 */
  detail: (id: string) => [...deviceKeys.details(), id] as const,
};

/** Point 相關 Query Keys */
export const pointKeys = {
  all: ['points'] as const,
  lists: () => [...pointKeys.all, 'list'] as const,
  list: (filters?: { device_id?: string; polling_group_id?: string }) =>
    [...pointKeys.lists(), filters ?? {}] as const,
  details: () => [...pointKeys.all, 'detail'] as const,
  detail: (id: string) => [...pointKeys.details(), id] as const,
};

/** Source Rule 相關 Query Keys */
export const sourceRuleKeys = {
  all: ['source-rules'] as const,
  lists: () => [...sourceRuleKeys.all, 'list'] as const,
  list: (filters?: { device_id?: string; enabled?: boolean }) =>
    [...sourceRuleKeys.lists(), filters ?? {}] as const,
  details: () => [...sourceRuleKeys.all, 'detail'] as const,
  detail: (id: string) => [...sourceRuleKeys.details(), id] as const,
  candidates: (id: string) => [...sourceRuleKeys.all, 'candidates', id] as const,
  reviewDecisions: (id: string) => [...sourceRuleKeys.all, 'review-decisions', id] as const,
};

/** Tag 相關 Query Keys */
export const tagKeys = {
  all: ['tags'] as const,
  lists: () => [...tagKeys.all, 'list'] as const,
  list: (filters?: { status?: string; data_type?: string; key_prefix?: string }) =>
    [...tagKeys.lists(), filters ?? {}] as const,
  details: () => [...tagKeys.all, 'detail'] as const,
  detail: (id: string) => [...tagKeys.details(), id] as const,
};

/** Mapping 相關 Query Keys */
export const mappingKeys = {
  all: ['mappings'] as const,
  lists: () => [...mappingKeys.all, 'list'] as const,
  list: (filters?: { point_id?: string; tag_id?: string; enabled?: boolean }) =>
    [...mappingKeys.lists(), filters ?? {}] as const,
  details: () => [...mappingKeys.all, 'detail'] as const,
  detail: (id: string) => [...mappingKeys.details(), id] as const,
};

/** Polling Group 相關 Query Keys */
export const pollingGroupKeys = {
  all: ['polling-groups'] as const,
  lists: () => [...pollingGroupKeys.all, 'list'] as const,
  list: () => [...pollingGroupKeys.lists()] as const,
};

/** Protocol 相關 Query Keys */
export const protocolKeys = {
  all: ['protocols'] as const,
  list: () => [...protocolKeys.all, 'list'] as const,
};

/** Settings 相關 Query Keys */
export const settingsKeys = {
  all: ['settings'] as const,
  items: () => [...settingsKeys.all, 'items'] as const,
};

/** Database Target 相關 Query Keys */
export const dbTargetKeys = {
  all: ['db-targets'] as const,
  connectors: () => [...dbTargetKeys.all, 'connectors'] as const,
};

/** Studio V2 Workspace 相關 Query Keys */
export const studioV2WorkspaceKeys = {
  all: ['studio-v2-workspace'] as const,
  bootstrap: () => [...studioV2WorkspaceKeys.all, 'bootstrap'] as const,
  runtimeContext: () => [...studioV2WorkspaceKeys.all, 'runtime-context'] as const,
  devices: () => [...studioV2WorkspaceKeys.all, 'devices'] as const,
  sourceRules: () => [...studioV2WorkspaceKeys.all, 'source-rules'] as const,
  mappings: () => [...studioV2WorkspaceKeys.all, 'mappings'] as const,
  databaseConfig: () => [...studioV2WorkspaceKeys.all, 'database-config'] as const,
  databaseTargets: () => [...studioV2WorkspaceKeys.all, 'database-targets'] as const,
};
