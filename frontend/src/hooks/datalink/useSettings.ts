import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { dbTargetAPI, settingsAPI } from '../../services/datalink';
import { dbTargetKeys, settingsKeys } from './keys';
import type {
  CreateDatabaseConnectorRequest,
  UpdateDatabaseConnectorRequest,
} from '../../types/datalink';

export function useSettingsItemsQuery() {
  return useQuery({
    queryKey: settingsKeys.items(),
    queryFn: () => settingsAPI.listItems(),
    refetchOnWindowFocus: false,
  });
}

export function useUpdateSettingKeyMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ key, value }: { key: string; value: unknown }) =>
      settingsAPI.updateKey(key, value),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: settingsKeys.all });
    },
  });
}

export function useDbTargetConnectorsQuery() {
  return useQuery({
    queryKey: dbTargetKeys.connectors(),
    queryFn: () => dbTargetAPI.listConnectors(),
    refetchOnWindowFocus: false,
  });
}

export function useCreateDbTargetConnectorMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateDatabaseConnectorRequest) =>
      dbTargetAPI.createConnector(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: dbTargetKeys.connectors() });
    },
  });
}

export function useUpdateDbTargetConnectorMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateDatabaseConnectorRequest }) =>
      dbTargetAPI.updateConnector(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: dbTargetKeys.connectors() });
    },
  });
}

export function useDeleteDbTargetConnectorMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => dbTargetAPI.deleteConnector(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: dbTargetKeys.connectors() });
    },
  });
}

export function useTestDbTargetConnectorMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => dbTargetAPI.testConnector(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: dbTargetKeys.connectors() });
    },
  });
}
