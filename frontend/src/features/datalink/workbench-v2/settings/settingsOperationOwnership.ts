export interface SettingsOperationOwnership {
  begin: (scope?: string) => number;
  isCurrent: (token: number, scope?: string) => boolean;
}

/** Tracks the latest request independently for each settings operation scope. */
export function createSettingsOperationOwnership(): SettingsOperationOwnership {
  const generations = new Map<string, number>();
  let nextToken = 0;
  return {
    begin: (scope = 'global') => {
      const token = ++nextToken;
      generations.set(scope, token);
      return token;
    },
    isCurrent: (token, scope = 'global') => token === generations.get(scope),
  };
}
