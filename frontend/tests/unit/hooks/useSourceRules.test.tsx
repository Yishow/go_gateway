import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, renderHook } from '@testing-library/react';
import type { ReactNode } from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { mappingKeys, tagKeys } from '@/hooks/datalink/keys';
import {
  useCreateSourceRuleMutation,
  useEnableSourceRuleMutation,
} from '@/hooks/datalink/useSourceRules';
import { sourceRuleAPI } from '@/services/datalink';
import type { SourceRuleRecord } from '@/types/datalink';

vi.mock('@/services/datalink', () => ({
  sourceRuleAPI: {
    create: vi.fn(),
    enable: vi.fn(),
  },
}));

function createRule(overrides: Partial<SourceRuleRecord> = {}): SourceRuleRecord {
  return {
    id: 'rule-1',
    device_id: 'device-1',
    start_address: '40001',
    count: 1,
    data_type: 'int16',
    naming_prefix: 'MIXER',
    enabled: true,
    locked: false,
    origin: 'manual',
    template_name: '',
    skipped_addresses: [],
    created_at: '',
    updated_at: '',
    ...overrides,
  };
}

function createWrapper(queryClient: QueryClient) {
  return function Wrapper({ children }: { children: ReactNode }) {
    return (
      <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    );
  };
}

describe('useSourceRules mutations', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('invalidates tag and mapping lists after create', async () => {
    const queryClient = new QueryClient({
      defaultOptions: {
        queries: { retry: false },
        mutations: { retry: false },
      },
    });
    const invalidateSpy = vi.spyOn(queryClient, 'invalidateQueries');
    vi.mocked(sourceRuleAPI.create).mockResolvedValue(createRule());

    const { result } = renderHook(() => useCreateSourceRuleMutation(), {
      wrapper: createWrapper(queryClient),
    });

    await act(async () => {
      await result.current.mutateAsync({
        device_id: 'device-1',
        start_address: '40001',
        count: 1,
        data_type: 'int16',
        naming_prefix: 'MIXER',
        enabled: true,
      });
    });

    expect(invalidateSpy).toHaveBeenCalledWith({ queryKey: tagKeys.lists() });
    expect(invalidateSpy).toHaveBeenCalledWith({
      queryKey: mappingKeys.lists(),
    });
  });

  it('invalidates tag and mapping lists after enable', async () => {
    const queryClient = new QueryClient({
      defaultOptions: {
        queries: { retry: false },
        mutations: { retry: false },
      },
    });
    const invalidateSpy = vi.spyOn(queryClient, 'invalidateQueries');
    vi.mocked(sourceRuleAPI.enable).mockResolvedValue(createRule());

    const { result } = renderHook(() => useEnableSourceRuleMutation(), {
      wrapper: createWrapper(queryClient),
    });

    await act(async () => {
      await result.current.mutateAsync('rule-1');
    });

    expect(invalidateSpy).toHaveBeenCalledWith({ queryKey: tagKeys.lists() });
    expect(invalidateSpy).toHaveBeenCalledWith({
      queryKey: mappingKeys.lists(),
    });
  });
});
