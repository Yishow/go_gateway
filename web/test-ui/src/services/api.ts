import axios from 'axios'
import { useState } from 'react'

const api = axios.create({
  baseURL: '/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
})

// 測試 API
export function useTestAPI() {
  const connect = async (protocol: string, config: Record<string, any>) => {
    const response = await api.post('/test/connect', {
      protocol,
      config,
    })
    return response.data
  }

  const disconnect = async (connectionId: string) => {
    await api.post('/test/disconnect', null, {
      params: { connection_id: connectionId },
    })
  }

  const getStatus = async (connectionId: string) => {
    const response = await api.get('/test/status', {
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
  ) => {
    const response = await api.post('/test/read', {
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
  ) => {
    const response = await api.post('/test/write', {
      connection_id: connectionId,
      ...params,
    })
    return response.data
  }

  return {
    connect,
    disconnect,
    getStatus,
    read,
    write,
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
