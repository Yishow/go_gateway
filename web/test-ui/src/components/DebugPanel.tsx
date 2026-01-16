import { useEffect, useState, useRef, useMemo } from 'react'
import { useDebugAPI } from '../services/api'

interface DebugPanelProps {
  connectionId: string | null
}

type DisplayMode = 'hex' | 'ascii' | 'parsed'
type FilterDirection = 'all' | 'request' | 'response'

export default function DebugPanel({ connectionId }: DebugPanelProps) {
  const [activeTab, setActiveTab] = useState<'packets' | 'logs'>('packets')
  const [autoScroll, setAutoScroll] = useState(true)
  const [displayMode, setDisplayMode] = useState<DisplayMode>('hex')
  const [filterDirection, setFilterDirection] = useState<FilterDirection>('all')
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedPacket, setSelectedPacket] = useState<string | null>(null)
  const [showFilters, setShowFilters] = useState(false)
  const scrollRef = useRef<HTMLDivElement>(null)
  const { packets, logs, refresh, clear } = useDebugAPI()

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

  // 過濾數據包
  const filteredPackets = useMemo(() => {
    let filtered = packets

    // 按方向過濾
    if (filterDirection !== 'all') {
      filtered = filtered.filter(p => p.direction === filterDirection)
    }

    // 按連接 ID 過濾
    if (connectionId) {
      filtered = filtered.filter(p => p.connection_id === connectionId)
    }

    // 搜索過濾
    if (searchQuery) {
      const query = searchQuery.toLowerCase()
      filtered = filtered.filter(p => 
        p.hex_data.toLowerCase().includes(query) ||
        p.protocol?.toLowerCase().includes(query) ||
        p.id.toLowerCase().includes(query)
      )
    }

    return filtered
  }, [packets, filterDirection, connectionId, searchQuery])

  // 過濾日誌
  const filteredLogs = useMemo(() => {
    let filtered = logs

    // 按連接 ID 過濾
    if (connectionId) {
      filtered = filtered.filter(l => l.details?.connection_id === connectionId)
    }

    // 搜索過濾
    if (searchQuery) {
      const query = searchQuery.toLowerCase()
      filtered = filtered.filter(l => 
        l.message.toLowerCase().includes(query) ||
        l.level.toLowerCase().includes(query)
      )
    }

    return filtered
  }, [logs, connectionId, searchQuery])

  // 統計信息
  const stats = useMemo(() => {
    const connPackets = packets.filter(p => p.connection_id === connectionId)
    const txCount = connPackets.filter(p => p.direction === 'request').length
    const rxCount = connPackets.filter(p => p.direction === 'response').length
    const errorLogs = logs.filter(l => 
      l.level === 'error' && l.details?.connection_id === connectionId
    ).length

    return { txCount, rxCount, errorLogs, total: connPackets.length }
  }, [packets, logs, connectionId])

  // 格式化數據包顯示
  const formatPacketData = (packet: any) => {
    switch (displayMode) {
      case 'ascii':
        return Array.from(packet.raw_data || [])
          .map(b => {
            const char = String.fromCharCode(b)
            return char >= ' ' && char <= '~' ? char : '.'
          })
          .join('')
      case 'parsed':
        return parseProtocol(packet)
      default:
        return packet.hex_data
    }
  }

  // 簡單的協議解析
  const parseProtocol = (packet: any): string => {
    if (!packet.raw_data || packet.raw_data.length === 0) return packet.hex_data
    
    const data = packet.raw_data
    const protocol = packet.protocol || ''

    if (protocol.includes('modbus')) {
      if (data.length >= 8) {
        const transactionId = (data[0] << 8) | data[1]
        const protocolId = (data[2] << 8) | data[3]
        const length = (data[4] << 8) | data[5]
        const unitId = data[6]
        const functionCode = data[7]
        
        if (protocolId === 0 && length > 0) {
          return `Modbus TCP | TID:${transactionId} | Unit:${unitId} | FC:0x${functionCode.toString(16).toUpperCase().padStart(2, '0')} | Len:${length}`
        }
      }
    }

    return packet.hex_data
  }

  // 導出數據
  const exportData = () => {
    const data = activeTab === 'packets' 
      ? filteredPackets.map(p => ({
          timestamp: p.timestamp,
          direction: p.direction,
          protocol: p.protocol,
          hex_data: p.hex_data,
          connection_id: p.connection_id
        }))
      : filteredLogs.map(l => ({
          timestamp: l.timestamp,
          level: l.level,
          message: l.message,
          details: l.details
        }))

    const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `debug_${activeTab}_${new Date().toISOString().slice(0, 19).replace(/:/g, '-')}.json`
    a.click()
    URL.revokeObjectURL(url)
  }

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
            Packet Stream ({filteredPackets.length})
            {stats.total > 0 && (
              <span className="ml-1 text-xs opacity-75">
                ({stats.txCount}↑/{stats.rxCount}↓)
              </span>
            )}
          </button>
          <button
            onClick={() => setActiveTab('logs')}
            className={`px-3 py-1.5 rounded-t-md transition-colors ${
              activeTab === 'logs'
                ? 'bg-gray-700 text-amber-400 font-bold'
                : 'text-gray-400 hover:text-gray-200 hover:bg-gray-700/50'
            }`}
          >
            System Logs ({filteredLogs.length})
            {stats.errorLogs > 0 && (
              <span className="ml-1 text-xs text-red-400">⚠{stats.errorLogs}</span>
            )}
          </button>
        </div>
        
        <div className="flex items-center space-x-2 text-gray-400">
          {/* 過濾按鈕 */}
          <button
            onClick={() => setShowFilters(!showFilters)}
            className={`px-2 py-1 rounded hover:bg-gray-700 transition-colors ${
              showFilters ? 'bg-gray-700 text-blue-400' : ''
            }`}
            title="Toggle Filters"
          >
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z" />
            </svg>
          </button>

          {/* 顯示模式切換（僅數據包） */}
          {activeTab === 'packets' && (
            <div className="flex items-center space-x-1 bg-gray-700 rounded px-1">
              <button
                onClick={() => setDisplayMode('hex')}
                className={`px-2 py-0.5 rounded text-xs transition-colors ${
                  displayMode === 'hex' ? 'bg-blue-600 text-white' : 'hover:bg-gray-600'
                }`}
                title="Hexadecimal"
              >
                HEX
              </button>
              <button
                onClick={() => setDisplayMode('ascii')}
                className={`px-2 py-0.5 rounded text-xs transition-colors ${
                  displayMode === 'ascii' ? 'bg-blue-600 text-white' : 'hover:bg-gray-600'
                }`}
                title="ASCII"
              >
                ASCII
              </button>
              <button
                onClick={() => setDisplayMode('parsed')}
                className={`px-2 py-0.5 rounded text-xs transition-colors ${
                  displayMode === 'parsed' ? 'bg-blue-600 text-white' : 'hover:bg-gray-600'
                }`}
                title="Parsed"
              >
                PARSED
              </button>
            </div>
          )}

          {/* 導出按鈕 */}
          <button
            onClick={exportData}
            className="px-2 py-1 rounded hover:bg-gray-700 transition-colors"
            title="Export Data"
          >
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </button>

          {/* 清空按鈕 */}
          <button
            onClick={() => {
              if (confirm('確定要清空當前數據嗎？')) {
                clear(connectionId)
              }
            }}
            className="px-2 py-1 rounded hover:bg-gray-700 hover:text-red-400 transition-colors"
            title="Clear Data"
          >
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
          </button>

          {/* Auto-scroll */}
          <label className="flex items-center space-x-1 cursor-pointer hover:text-gray-200">
            <input 
              type="checkbox" 
              checked={autoScroll} 
              onChange={(e) => setAutoScroll(e.target.checked)}
              className="rounded border-gray-600 bg-gray-700 text-blue-500 focus:ring-0 w-3 h-3"
            />
            <span className="text-xs">Auto</span>
          </label>
        </div>
      </div>

      {/* 過濾器面板 */}
      {showFilters && (
        <div className="px-2 py-1.5 bg-gray-800/50 border-b border-gray-700 flex items-center gap-2 flex-wrap">
          {/* 方向過濾 */}
          {activeTab === 'packets' && (
            <div className="flex items-center space-x-1">
              <span className="text-gray-500 text-xs">方向:</span>
              <select
                value={filterDirection}
                onChange={(e) => setFilterDirection(e.target.value as FilterDirection)}
                className="bg-gray-700 text-gray-300 text-xs px-2 py-0.5 rounded border border-gray-600 focus:outline-none focus:border-blue-500"
              >
                <option value="all">全部</option>
                <option value="request">TX (發送)</option>
                <option value="response">RX (接收)</option>
              </select>
            </div>
          )}

          {/* 搜索框 */}
          <div className="flex items-center space-x-1 flex-1 min-w-[200px]">
            <svg className="w-3 h-3 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="搜索..."
              className="bg-gray-700 text-gray-300 text-xs px-2 py-0.5 rounded border border-gray-600 focus:outline-none focus:border-blue-500 flex-1"
            />
            {searchQuery && (
              <button
                onClick={() => setSearchQuery('')}
                className="text-gray-500 hover:text-gray-300"
              >
                <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            )}
          </div>
        </div>
      )}

      {/* Content Area */}
      <div 
        ref={scrollRef}
        className="flex-1 overflow-auto p-2 space-y-1 bg-gray-900 scrollbar-thin scrollbar-thumb-gray-700 scrollbar-track-transparent"
      >
        {activeTab === 'packets' && (
          <>
            {filteredPackets.length === 0 ? (
              <div className="flex flex-col items-center justify-center h-full text-gray-600">
                <span className="text-4xl opacity-20 mb-2">⚡</span>
                <p>{searchQuery || filterDirection !== 'all' ? '沒有匹配的數據包' : 'Waiting for data stream...'}</p>
              </div>
            ) : (
              filteredPackets.map((packet) => (
                <div
                  key={packet.id}
                  onClick={() => setSelectedPacket(selectedPacket === packet.id ? null : packet.id)}
                  className={`group flex gap-3 p-1.5 hover:bg-gray-800 rounded border transition-all cursor-pointer ${
                    selectedPacket === packet.id ? 'bg-gray-800 border-blue-500' : 'border-transparent hover:border-gray-700'
                  }`}
                >
                  <span className="text-gray-500 w-20 shrink-0">
                    {new Date(packet.timestamp).toLocaleTimeString([], { 
                      hour12: false, 
                      hour: '2-digit', 
                      minute: '2-digit', 
                      second: '2-digit', 
                      fractionalSecondDigits: 3 
                    })}
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
                  <div className="flex-1 min-w-0">
                    <div className="text-gray-300 break-all font-mono opacity-90 group-hover:opacity-100">
                      {formatPacketData(packet)}
                    </div>
                    {selectedPacket === packet.id && (
                      <div className="mt-1 pt-1 border-t border-gray-700 text-xs text-gray-400">
                        <div>Protocol: {packet.protocol || 'N/A'}</div>
                        <div>Length: {packet.raw_data?.length || 0} bytes</div>
                        <div>Connection: {packet.connection_id}</div>
                        <button
                          onClick={(e) => {
                            e.stopPropagation()
                            navigator.clipboard.writeText(packet.hex_data)
                          }}
                          className="mt-1 text-blue-400 hover:text-blue-300"
                        >
                          複製 HEX
                        </button>
                      </div>
                    )}
                  </div>
                </div>
              ))
            )}
          </>
        )}

        {activeTab === 'logs' && (
          <>
            {filteredLogs.length === 0 ? (
              <div className="flex flex-col items-center justify-center h-full text-gray-600">
                <span className="text-4xl opacity-20 mb-2">📝</span>
                <p>{searchQuery ? '沒有匹配的日誌' : 'No logs available'}</p>
              </div>
            ) : (
              filteredLogs.map((log) => (
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
                    {new Date(log.timestamp).toLocaleTimeString([], { 
                      hour12: false, 
                      hour: '2-digit', 
                      minute: '2-digit', 
                      second: '2-digit', 
                      fractionalSecondDigits: 3 
                    })}
                  </span>
                  <span className={`uppercase font-bold shrink-0 w-12 ${
                    log.level === 'error' ? 'text-red-400' : log.level === 'warn' ? 'text-yellow-400' : 'text-blue-400'
                  }`}>
                    {log.level}
                  </span>
                  <span className="break-all flex-1">{log.message}</span>
                </div>
              ))
            )}
          </>
        )}
      </div>
    </div>
  )
}
