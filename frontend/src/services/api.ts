import axios from 'axios'
import { useState } from 'react'
import { logger } from '../utils/logger';
import { resolveGatewayUiVersion } from '../features/gateway/uiVersion';
import type { 
  ConnectResponse, 
  ReadResponse, 
  ConnectionState
} from '../types/api'
import { VITE_API_BASE_URL } from '../env'

const apiBaseURL = VITE_API_BASE_URL

const api = axios.create({
  baseURL: apiBaseURL,
  headers: {
    'Content-Type': 'application/json',
  },
})

api.interceptors.request.use((config) => {
  const uiVersion = resolveGatewayUiVersion();
  if (uiVersion) {
    config.headers = config.headers ?? {};
    if (!config.headers['X-UI-Version']) {
      config.headers['X-UI-Version'] = uiVersion;
    }
  }
  return config;
})

interface DebugPacket {
  id: string
  connection_id?: string
  direction: 'request' | 'response' | string
  protocol?: string
  hex_data: string
  raw_data?: number[]
  timestamp: string
  [key: string]: unknown
}

interface DebugLog {
  id: string
  timestamp: string
  level: string
  message: string
  details?: {
    connection_id?: string
    [key: string]: unknown
  }
  [key: string]: unknown
}

// 測試 API
export function useTestAPI() {
  const connect = async (protocol: string, config: Record<string, unknown>): Promise<ConnectResponse> => {
    const response = await api.post<ConnectResponse>('/test/connect', {
      protocol,
      config,
    })
    return response.data
  }

  const disconnect = async (connectionId: string): Promise<void> => {
    await api.post('/test/disconnect', null, {
      params: { connection_id: connectionId },
    })
  }

  const getStatus = async (connectionId: string): Promise<ConnectionState> => {
    const response = await api.get<ConnectionState>('/test/status', {
      params: { connection_id: connectionId },
    })
    return response.data
  }

  const read = async (
    connectionId: string,
    params: {
      operation: string
      address: number
      count: number
      symbol?: string
      device?: string
      unit_id?: number
      station?: number
    }
  ): Promise<ReadResponse> => {
    const response = await api.post<ReadResponse>('/test/read', {
      connection_id: connectionId,
      ...params,
    })
    return response.data
  }

  const write = async (
    connectionId: string,
    params: {
      operation: string
      address: number
      values: Array<string | number | boolean>
      symbol?: string
      unit_id?: number
      station?: number
    }
  ): Promise<void> => {
    await api.post('/test/write', {
      connection_id: connectionId,
      ...params,
    })
  }

  const batch = async (req: import('../types/api').BatchRequest): Promise<import('../types/api').BatchResponse> => {
    const response = await api.post<import('../types/api').BatchResponse>('/test/batch', req)
    return response.data
  }

  const startMonitor = async (req: import('../types/api').MonitorRequest): Promise<void> => {
    await api.post('/test/monitor/start', req)
  }

  const stopMonitor = async (connectionId: string): Promise<void> => {
    await api.post('/test/monitor/stop', { connection_id: connectionId })
  }

  return {
    connect,
    disconnect,
    getStatus,
    read,
    write,
    batch,
    startMonitor,
    stopMonitor,
  }
}

// Debug API Hook
export function useDebugAPI() {
  const [packets, setPackets] = useState<DebugPacket[]>([])
  const [logs, setLogs] = useState<DebugLog[]>([])

  const refresh = async () => {
    try {
      const [packetsRes, logsRes] = await Promise.all([
        api.get<{ packets?: DebugPacket[] }>('/debug/packets'),
        api.get<{ logs?: DebugLog[] }>('/debug/logs'),
      ])
      setPackets(packetsRes.data.packets || [])
      setLogs(logsRes.data.logs || [])
    } catch (error) {
      logger.error('Failed to refresh debug data:', error)
    }
  }

  const clear = async (connectionId: string | null) => {
    try {
      const params = connectionId ? { connection_id: connectionId } : {}
      const beforePacketsCount = packets.length
      const beforeLogsCount = logs.length
      
      const response = await api.delete('/debug/clear', { params })
      logger.log('Clear response:', response.data)
      
      // 立即更新本地狀態，不調用 refresh()，讓自動刷新機制自然處理
      // 這樣可以避免在清空後立即重新載入數據
      if (connectionId) {
        setPackets(prev => {
          const filtered = prev.filter(p => p.connection_id !== connectionId)
          logger.log(`Cleared packets for connection ${connectionId}: ${prev.length} -> ${filtered.length}`)
          return filtered
        })
        setLogs(prev => {
          const filtered = prev.filter(l => {
            const details = l.details as { connection_id?: string } | undefined
            return details?.connection_id !== connectionId
          })
          logger.log(`Cleared logs for connection ${connectionId}: ${prev.length} -> ${filtered.length}`)
          return filtered
        })
      } else {
        logger.log(`Cleared all: packets ${beforePacketsCount} -> 0, logs ${beforeLogsCount} -> 0`)
        setPackets([])
        setLogs([])
      }
    } catch (error) {
      logger.error('Failed to clear debug data:', error)
      // 如果後端不支持，則前端清空
      if (connectionId) {
        setPackets(prev => prev.filter(p => p.connection_id !== connectionId))
        setLogs(prev => prev.filter(l => {
          const details = l.details as { connection_id?: string } | undefined
          return details?.connection_id !== connectionId
        }))
      } else {
        setPackets([])
        setLogs([])
      }
    }
  }

  return {
    packets,
    logs,
    refresh,
    clear,
  }
}
