import { useEffect, useState, useRef } from 'react'
import { useDebugAPI } from '../services/api'

interface DebugPanelProps {
  connectionId: string | null
}

export default function DebugPanel({ connectionId }: DebugPanelProps) {
  const [activeTab, setActiveTab] = useState<'packets' | 'logs'>('packets')
  const [autoScroll, setAutoScroll] = useState(true)
  const scrollRef = useRef<HTMLDivElement>(null)
  const { packets, logs, refresh } = useDebugAPI()

  useEffect(() => {
    if (connectionId) {
      const interval = setInterval(() => {
        refresh()
      }, 1000)
      return () => clearInterval(interval)
    }
  }, [connectionId, refresh])

  useEffect(() => {
    if (autoScroll && scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight
    }
  }, [packets, logs, autoScroll])

  return (
    <div className="flex flex-col h-[500px] bg-gray-900 rounded-xl shadow-xl border border-gray-700 overflow-hidden font-mono text-xs">
      {/* Header / Tabs */}
      <div className="flex items-center justify-between px-2 py-1 bg-gray-800 border-b border-gray-700">
        <div className="flex space-x-1">
          <button
            onClick={() => setActiveTab('packets')}
            className={`px-3 py-1.5 rounded-t-md transition-colors ${
              activeTab === 'packets'
                ? 'bg-gray-700 text-blue-400 font-bold'
                : 'text-gray-400 hover:text-gray-200 hover:bg-gray-700/50'
            }`}
          >
            Packet Stream ({packets.length})
          </button>
          <button
            onClick={() => setActiveTab('logs')}
            className={`px-3 py-1.5 rounded-t-md transition-colors ${
              activeTab === 'logs'
                ? 'bg-gray-700 text-amber-400 font-bold'
                : 'text-gray-400 hover:text-gray-200 hover:bg-gray-700/50'
            }`}
          >
            System Logs ({logs.length})
          </button>
        </div>
        
        <div className="flex items-center space-x-2 text-gray-400 px-2">
           <label className="flex items-center space-x-1 cursor-pointer hover:text-gray-200">
            <input 
              type="checkbox" 
              checked={autoScroll} 
              onChange={(e) => setAutoScroll(e.target.checked)}
              className="rounded border-gray-600 bg-gray-700 text-blue-500 focus:ring-0 w-3 h-3"
            />
            <span>Auto-scroll</span>
          </label>
          <button 
            onClick={() => {/* Implement clear logic if API supports it, or local clear */}} 
            className="hover:text-red-400"
            title="Clear Console"
          >
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
          </button>
        </div>
      </div>

      {/* Content Area */}
      <div 
        ref={scrollRef}
        className="flex-1 overflow-auto p-2 space-y-1 bg-gray-900 scrollbar-thin scrollbar-thumb-gray-700 scrollbar-track-transparent"
      >
        {activeTab === 'packets' && (
          <>
            {packets.length === 0 ? (
              <div className="flex flex-col items-center justify-center h-full text-gray-600">
                <span className="text-4xl opacity-20 mb-2">⚡</span>
                <p>Waiting for data stream...</p>
              </div>
            ) : (
              packets.map((packet) => (
                <div
                  key={packet.id}
                  className="group flex gap-3 p-1.5 hover:bg-gray-800 rounded border border-transparent hover:border-gray-700 transition-all"
                >
                  <span className="text-gray-500 w-20 shrink-0">
                    {new Date(packet.timestamp).toLocaleTimeString([], { hour12: false, hour: '2-digit', minute: '2-digit', second: '2-digit', fractionalSecondDigits: 3 })}
                  </span>
                  <div className="flex items-center gap-2 shrink-0 w-16">
                     {packet.direction === 'request' ? (
                       <span className="text-blue-400 font-bold flex items-center gap-1">
                         <span className="text-xs">TX</span> →
                       </span>
                     ) : (
                       <span className="text-green-400 font-bold flex items-center gap-1">
                         ← <span className="text-xs">RX</span>
                       </span>
                     )}
                  </div>
                  <div className="text-gray-300 break-all font-mono opacity-90 group-hover:opacity-100">
                    {packet.hex_data}
                  </div>
                </div>
              ))
            )}
          </>
        )}

        {activeTab === 'logs' && (
          <>
            {logs.length === 0 ? (
              <div className="flex flex-col items-center justify-center h-full text-gray-600">
                <span className="text-4xl opacity-20 mb-2">📝</span>
                <p>No logs available</p>
              </div>
            ) : (
              logs.map((log) => (
                <div
                  key={log.id}
                  className={`flex gap-3 p-1.5 rounded border-l-2 ${
                    log.level === 'error'
                      ? 'border-red-500 bg-red-900/10 text-red-200'
                      : log.level === 'warn'
                      ? 'border-yellow-500 bg-yellow-900/10 text-yellow-200'
                      : 'border-blue-500 bg-gray-800/30 text-gray-300'
                  }`}
                >
                   <span className="text-gray-500 w-20 shrink-0">
                    {new Date(log.timestamp).toLocaleTimeString([], { hour12: false, hour: '2-digit', minute: '2-digit', second: '2-digit', fractionalSecondDigits: 3 })}
                  </span>
                  <span className={`uppercase font-bold shrink-0 w-12 ${
                     log.level === 'error' ? 'text-red-400' : log.level === 'warn' ? 'text-yellow-400' : 'text-blue-400'
                  }`}>
                    {log.level}
                  </span>
                  <span className="break-all">{log.message}</span>
                </div>
              ))
            )}
          </>
        )}
      </div>
    </div>
  )
}
