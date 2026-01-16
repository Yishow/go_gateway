import { useState } from 'react'

interface MonitorControlProps {
  connectionId: string | null
}

export default function MonitorControl({ connectionId }: MonitorControlProps) {
  const [monitoring, setMonitoring] = useState(false)
  const [interval, setInterval] = useState<number>(1000)

  const handleStart = () => {
    if (!connectionId) {
      alert('請先建立連線')
      return
    }
    setMonitoring(true)
    // TODO: 啟動監控模式
  }

  const handleStop = () => {
    setMonitoring(false)
    // TODO: 停止監控模式
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
        <div className="mt-4 p-4 bg-gray-50 rounded-md">
          <p className="text-sm text-gray-600">監控模式運行中...</p>
          {/* TODO: 顯示監控數據圖表 */}
        </div>
      )}
    </div>
  )
}
