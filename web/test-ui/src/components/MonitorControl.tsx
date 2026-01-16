import { useState, useEffect, useRef } from 'react'
import { useTestAPI } from '../services/api'

interface MonitorControlProps {
  connectionId: string | null
}

export default function MonitorControl({ connectionId }: MonitorControlProps) {
  const [monitoring, setMonitoring] = useState(false)
  const [interval, setInterval] = useState<number>(1000)
  const [monitorData, setMonitorData] = useState<any[]>([])
  const { startMonitor, stopMonitor } = useTestAPI()
  const socketRef = useRef<WebSocket | null>(null)

  useEffect(() => {
    // 建立 WebSocket 連線用於接收監控數據
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = window.location.host
    const wsUrl = `${protocol}//${host}/api/v1/test/monitor/stream`
    
    const socket = new WebSocket(wsUrl)
    socketRef.current = socket

    socket.onmessage = (event) => {
      const msg = JSON.parse(event.data)
      if (msg.type === 'monitor_update' && msg.connection_id === connectionId) {
        setMonitorData((prev) => [msg, ...prev].slice(0, 50)) // 保留最近 50 筆
      }
    }

    socket.onclose = () => {
      console.log('Monitor WebSocket closed')
    }

    return () => {
      socket.close()
    }
  }, [connectionId])

  const handleStart = async () => {
    if (!connectionId) {
      alert('請先建立連線')
      return
    }
    
    try {
      // 這裡簡單示範監控第一個地址，實際應用中應允許用戶配置多個項目
      await startMonitor({
        connection_id: connectionId,
        interval: interval,
        items: [
          {
            connection_id: connectionId,
            operation: 'read_holding_registers', // 預設操作，視協議而定
            address: 0,
            count: 10
          }
        ]
      })
      setMonitoring(true)
      setMonitorData([])
    } catch (error: any) {
      alert('啟動監控失敗: ' + error.message)
    }
  }

  const handleStop = async () => {
    if (!connectionId) return
    try {
      await stopMonitor(connectionId)
      setMonitoring(false)
    } catch (error: any) {
      alert('停止監控失敗: ' + error.message)
    }
  }

  return (
    <div className="space-y-4">
      <h3 className="text-lg font-semibold">監控模式</h3>
      <div className="flex items-center space-x-4">
        <div className="flex-1">
          <label className="block text-sm font-medium text-gray-700 mb-1">
            讀取間隔 (毫秒)
          </label>
          <input
            type="number"
            value={interval}
            onChange={(e) => setInterval(parseInt(e.target.value) || 1000)}
            disabled={monitoring}
            className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100"
          />
        </div>
        <div className="flex items-end">
          {!monitoring ? (
            <button
              onClick={handleStart}
              disabled={!connectionId}
              className="bg-green-600 text-white px-6 py-2 rounded-md hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              啟動監控
            </button>
          ) : (
            <button
              onClick={handleStop}
              className="bg-red-600 text-white px-6 py-2 rounded-md hover:bg-red-700"
            >
              停止監控
            </button>
          )}
        </div>
      </div>

      {monitoring && (
        <div className="mt-4 p-4 bg-gray-50 rounded-md border border-gray-200">
          <div className="flex justify-between items-center mb-2">
            <p className="text-sm font-medium text-blue-600 animate-pulse">● 監控模式運行中</p>
            <span className="text-xs text-gray-500">保留最近 50 筆數據</span>
          </div>
          
          <div className="max-h-64 overflow-y-auto space-y-2">
            {monitorData.map((entry, idx) => (
              <div key={idx} className="text-xs font-mono bg-white p-2 rounded shadow-sm border-l-4 border-blue-400">
                <div className="text-gray-400 mb-1">{new Date(entry.timestamp).toLocaleTimeString()}</div>
                <div className="grid grid-cols-1 gap-1">
                  {Object.entries(entry.data).map(([key, val]: [string, any]) => (
                    <div key={key} className="flex">
                      <span className="font-bold mr-2">{key}:</span>
                      <span className="text-green-700">
                        {val.error ? `Error: ${val.error}` : JSON.stringify(val)}
                      </span>
                    </div>
                  ))}
                </div>
              </div>
            ))}
            {monitorData.length === 0 && (
              <p className="text-center text-gray-400 py-4 italic">等待數據中...</p>
            )}
          </div>
        </div>
      )}
    </div>
  )
}
