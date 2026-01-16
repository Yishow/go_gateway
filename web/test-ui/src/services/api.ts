import axios from 'axios'
import { useState } from 'react'
import type { 
  ConnectResponse, 
  ReadResponse, 
  ConnectionState
} from '../types/api'

const api = axios.create({
  baseURL: '/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
})

// 測試 API
export function useTestAPI() {
  const connect = async (protocol: string, config: Record<string, any>): Promise<ConnectResponse> => {
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
      values: any[]
      symbol?: string
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
  const [packets, setPackets] = useState<any[]>([])
  const [logs, setLogs] = useState<any[]>([])

  const refresh = async () => {
    try {
      const [packetsRes, logsRes] = await Promise.all([
        api.get('/debug/packets'),
        api.get('/debug/logs'),
      ])
      setPackets(packetsRes.data.packets || [])
      setLogs(logsRes.data.logs || [])
    } catch (error) {
      console.error('Failed to refresh debug data:', error)
    }
  }

  return {
    packets,
    logs,
    refresh,
  }
}
