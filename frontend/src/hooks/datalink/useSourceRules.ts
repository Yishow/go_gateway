import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { sourceRuleAPI } from '../../services/datalink';
import { mappingKeys, pointKeys, sourceRuleKeys, tagKeys } from './keys';
import type {
  CreateSourceRuleRequest,
  UpdateSourceRuleRequest,
} from '../../types/datalink';

export interface SourceRuleListFilters {
  device_id?: string;
  enabled?: boolean;
}

export function useSourceRulesQuery(filters?: SourceRuleListFilters) {
  return useQuery({
    queryKey: sourceRuleKeys.list(filters),
    queryFn: () => sourceRuleAPI.list(filters),
    enabled: Boolean(filters?.device_id),
  });
}

export function useCreateSourceRuleMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateSourceRuleRequest) => sourceRuleAPI.create(data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: sourceRuleKeys.lists() });
      queryClient.invalidateQueries({ queryKey: pointKeys.lists() });
      queryClient.invalidateQueries({ queryKey: tagKeys.lists() });
      queryClient.invalidateQueries({ queryKey: mappingKeys.lists() });
      if (variables.device_id) {
        queryClient.invalidateQueries({
          queryKey: sourceRuleKeys.list({ device_id: variables.device_id }),
        });
      }
    },
  });
}

export function useUpdateSourceRuleMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateSourceRuleRequest }) =>
      sourceRuleAPI.update(id, data),
    onSuccess: (rule, { id }) => {
      queryClient.invalidateQueries({ queryKey: sourceRuleKeys.lists() });
      queryClient.invalidateQueries({ queryKey: sourceRuleKeys.detail(id) });
      queryClient.invalidateQueries({ queryKey: pointKeys.lists() });
      queryClient.invalidateQueries({ queryKey: tagKeys.lists() });
      queryClient.invalidateQueries({ queryKey: mappingKeys.lists() });
      queryClient.invalidateQueries({
        queryKey: sourceRuleKeys.list({ device_id: rule.device_id }),
      });
    },
  });
}

export function useDeleteSourceRuleMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => sourceRuleAPI.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: sourceRuleKeys.lists() });
      queryClient.invalidateQueries({ queryKey: pointKeys.lists() });
      queryClient.invalidateQueries({ queryKey: tagKeys.lists() });
      queryClient.invalidateQueries({ queryKey: mappingKeys.lists() });
    },
  });
}

export function useEnableSourceRuleMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => sourceRuleAPI.enable(id),
    onSuccess: (rule, id) => {
      queryClient.invalidateQueries({ queryKey: sourceRuleKeys.lists() });
      queryClient.invalidateQueries({ queryKey: sourceRuleKeys.detail(id) });
      queryClient.invalidateQueries({ queryKey: pointKeys.lists() });
      queryClient.invalidateQueries({ queryKey: tagKeys.lists() });
      queryClient.invalidateQueries({ queryKey: mappingKeys.lists() });
      queryClient.invalidateQueries({
        queryKey: sourceRuleKeys.list({ device_id: rule.device_id }),
      });
    },
  });
}

export function useDisableSourceRuleMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => sourceRuleAPI.disable(id),
    onSuccess: (rule, id) => {
      queryClient.invalidateQueries({ queryKey: sourceRuleKeys.lists() });
      queryClient.invalidateQueries({ queryKey: sourceRuleKeys.detail(id) });
      queryClient.invalidateQueries({ queryKey: pointKeys.lists() });
      queryClient.invalidateQueries({ queryKey: tagKeys.lists() });
      queryClient.invalidateQueries({ queryKey: mappingKeys.lists() });
      queryClient.invalidateQueries({
        queryKey: sourceRuleKeys.list({ device_id: rule.device_id }),
      });
    },
  });
}
