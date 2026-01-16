import { useState, useEffect, useRef } from 'react'
import { useTestAPI } from '../services/api'
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'

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
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = window.location.host
    const wsUrl = `${protocol}//${host}/api/v1/test/monitor/stream`
    
    const socket = new WebSocket(wsUrl)
    socketRef.current = socket

    socket.onmessage = (event) => {
      const msg = JSON.parse(event.data)
      if (msg.type === 'monitor_update' && msg.connection_id === connectionId) {
        setMonitorData((prev) => {
          // Flatten data for chart: { time: '...', value: 123 }
          // Assuming the first key in data is the value we want to plot for now
          // In a real app, we'd map multiple lines
          const dataPoints = msg.data
          const firstKey = Object.keys(dataPoints)[0]
          const numericValue = typeof dataPoints[firstKey] === 'number' ? dataPoints[firstKey] : null
          
          const newEntry = {
            ...msg,
            chartTime: new Date(msg.timestamp).toLocaleTimeString(),
            val: numericValue,
            ...dataPoints // spread other data for flexible access
          }
          return [newEntry, ...prev].slice(0, 50) 
        })
      }
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
      await startMonitor({
        connection_id: connectionId,
        interval: interval,
        items: [
          {
            connection_id: connectionId,
            operation: 'read_holding_registers',
            address: 0,
            count: 1
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

  // Reverse data for chart (oldest to newest)
  const chartData = [...monitorData].reverse()

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between border-b pb-4">
        <h2 className="text-lg font-bold text-gray-800 flex items-center gap-2">
          <span className="w-1.5 h-6 bg-purple-500 rounded-full"></span>
          即時監控
        </h2>
        <div className="flex items-center gap-3">
           <div className="relative">
            <span className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 text-xs">間隔(ms)</span>
             <input
              type="number"
              value={interval}
              onChange={(e) => setInterval(parseInt(e.target.value) || 1000)}
              disabled={monitoring}
              className="pl-16 pr-3 py-2 w-32 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500 transition-all disabled:bg-gray-100"
            />
           </div>
           
           {!monitoring ? (
            <button
              onClick={handleStart}
              disabled={!connectionId}
              className="bg-purple-600 text-white px-6 py-2 rounded-lg font-semibold shadow-lg shadow-purple-500/30 hover:bg-purple-700 hover:shadow-purple-500/40 active:scale-[0.98] transition-all disabled:opacity-50 disabled:cursor-not-allowed"
            >
              啟動監控
            </button>
          ) : (
            <button
              onClick={handleStop}
              className="bg-red-500 text-white px-6 py-2 rounded-lg font-semibold shadow-lg shadow-red-500/30 hover:bg-red-600 hover:shadow-red-500/40 active:scale-[0.98] transition-all flex items-center gap-2"
            >
              <span className="animate-pulse w-2 h-2 rounded-full bg-white"></span>
              停止監控
            </button>
          )}
        </div>
      </div>

      {monitoring && (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 animate-fade-in">
          {/* Chart Area */}
          <div className="lg:col-span-2 bg-gray-50 rounded-xl p-4 border border-gray-200 shadow-inner h-[300px]">
            <ResponsiveContainer width="100%" height="100%">
              <LineChart data={chartData}>
                <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
                <XAxis 
                  dataKey="chartTime" 
                  tick={{fontSize: 10}} 
                  interval="preserveStartEnd"
                  stroke="#9ca3af"
                />
                <YAxis stroke="#9ca3af" tick={{fontSize: 10}} />
                <Tooltip 
                  contentStyle={{ backgroundColor: '#1f2937', border: 'none', borderRadius: '8px', color: '#fff' }}
                />
                <Line 
                  type="monotone" 
                  dataKey="val" 
                  stroke="#8b5cf6" 
                  strokeWidth={2} 
                  dot={false} 
                  activeDot={{ r: 6 }} 
                  animationDuration={300}
                />
              </LineChart>
            </ResponsiveContainer>
          </div>

          {/* Data Log Area */}
          <div className="lg:col-span-1 bg-white rounded-xl border border-gray-200 flex flex-col h-[300px] overflow-hidden">
             <div className="bg-gray-50 px-4 py-2 border-b border-gray-200 text-xs font-bold text-gray-500 uppercase">
               數據日誌 ({monitorData.length})
             </div>
             <div className="flex-1 overflow-y-auto p-2 space-y-2 scrollbar-thin">
                {monitorData.map((entry, idx) => (
                  <div key={idx} className="text-xs p-2 rounded border border-gray-100 shadow-sm hover:shadow-md transition-shadow bg-white">
                    <div className="flex justify-between text-gray-400 mb-1">
                      <span>#{idx + 1}</span>
                      <span>{new Date(entry.timestamp).toLocaleTimeString()}</span>
                    </div>
                    <div className="font-mono text-gray-700 break-all">
                      {Object.entries(entry.data).map(([key, val]: [string, any]) => (
                        <span key={key} className="block">
                          <span className="text-purple-600 font-semibold">{key}:</span> {String(val)}
                        </span>
                      ))}
                    </div>
                  </div>
                ))}
             </div>
          </div>
        </div>
      )}
      
      {!monitoring && (
        <div className="flex flex-col items-center justify-center py-12 bg-gray-50 rounded-xl border border-dashed border-gray-300 text-gray-400">
           <svg className="w-12 h-12 mb-2 opacity-50" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M7 12l3-3 3 3 4-4M8 21l4-4 4 4M3 4h18M4 4h16v12a1 1 0 01-1 1H5a1 1 0 01-1-1V4z" />
          </svg>
          <p>準備就緒，點擊啟動監控以開始</p>
        </div>
      )}
    </div>
  )
}
