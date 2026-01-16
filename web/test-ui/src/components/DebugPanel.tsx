import { useEffect, useState } from 'react'
import { useDebugAPI } from '../services/api'

interface DebugPanelProps {
  connectionId: string | null
}

export default function DebugPanel({ connectionId }: DebugPanelProps) {
  const [activeTab, setActiveTab] = useState<'packets' | 'logs'>('packets')
  const { packets, logs, refresh } = useDebugAPI()

  useEffect(() => {
    if (connectionId) {
      const interval = setInterval(() => {
        refresh()
      }, 1000)
      return () => clearInterval(interval)
    }
  }, [connectionId, refresh])

  return (
    <div className="space-y-4">
      {/* 標籤切換 */}
      <div className="flex border-b">
        <button
          onClick={() => setActiveTab('packets')}
          className={`px-4 py-2 text-sm font-medium ${
            activeTab === 'packets'
              ? 'border-b-2 border-blue-500 text-blue-600'
              : 'text-gray-500 hover:text-gray-700'
          }`}
        >
          數據包
        </button>
        <button
          onClick={() => setActiveTab('logs')}
          className={`px-4 py-2 text-sm font-medium ${
            activeTab === 'logs'
              ? 'border-b-2 border-blue-500 text-blue-600'
              : 'text-gray-500 hover:text-gray-700'
          }`}
        >
          日誌
        </button>
      </div>

      {/* 內容區域 */}
      <div className="h-64 overflow-auto">
        {activeTab === 'packets' && (
          <div className="space-y-2">
            {packets.length === 0 ? (
              <p className="text-gray-500 text-sm text-center py-8">
                尚無數據包記錄
              </p>
            ) : (
              packets.map((packet) => (
                <div
                  key={packet.id}
                  className="p-2 bg-gray-50 rounded text-xs font-mono"
                >
                  <div className="flex justify-between mb-1">
                    <span className="font-semibold">
                      {packet.direction === 'request' ? '→' : '←'}
                    </span>
                    <span className="text-gray-500">
                      {new Date(packet.timestamp).toLocaleTimeString()}
                    </span>
                  </div>
                  <div className="text-gray-700">{packet.hex_data}</div>
                </div>
              ))
            )}
          </div>
        )}

        {activeTab === 'logs' && (
          <div className="space-y-2">
            {logs.length === 0 ? (
              <p className="text-gray-500 text-sm text-center py-8">
                尚無日誌記錄
              </p>
            ) : (
              logs.map((log) => (
                <div
                  key={log.id}
                  className={`p-2 rounded text-xs ${
                    log.level === 'error'
                      ? 'bg-red-50 text-red-800'
                      : log.level === 'warn'
                      ? 'bg-yellow-50 text-yellow-800'
                      : 'bg-gray-50 text-gray-800'
                  }`}
                >
                  <div className="flex justify-between mb-1">
                    <span className="font-semibold uppercase">{log.level}</span>
                    <span className="text-gray-500">
                      {new Date(log.timestamp).toLocaleTimeString()}
                    </span>
                  </div>
                  <div>{log.message}</div>
                </div>
              ))
            )}
          </div>
        )}
      </div>
    </div>
  )
}
