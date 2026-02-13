import { useState, useEffect, useRef } from 'react'
import { useTestAPI } from '../services/api'
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts'
import type { ReadRequest } from '../types/api'
import { useProfiles } from '../hooks/useProfiles'
import type { MonitorItem } from '../types/profile'
import { useToast } from '../contexts/ToastContext'

interface MonitorControlProps {
  connectionId: string | null
  protocol: string
}

interface MonitorStreamEntry {
  timestamp: string
  chartTime: string
  data: Record<string, unknown>
}

const toNumericValue = (value: unknown): number | null => {
  if (typeof value === 'number') {
    return Number.isFinite(value) ? value : null
  }
  if (typeof value === 'string') {
    const trimmed = value.trim()
    if (trimmed.length === 0) return null
    const parsed = Number(trimmed)
    return Number.isFinite(parsed) ? parsed : null
  }
  return null
}

export default function MonitorControl({ connectionId, protocol }: MonitorControlProps) {
  const [monitoring, setMonitoring] = useState(false)
  const [interval, setInterval] = useState<number>(1000)
  const [monitorItems, setMonitorItems] = useState<MonitorItem[]>([])
  const [monitorData, setMonitorData] = useState<MonitorStreamEntry[]>([])
  const { startMonitor, stopMonitor } = useTestAPI()
  const { currentProfile, updateProfile } = useProfiles()
  const { showError, showWarning } = useToast()
  const eventSourceRef = useRef<EventSource | null>(null)
  const saveTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const reconnectTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const reconnectAttemptsRef = useRef(0)

  const isModbus = protocol.includes('modbus')
  const isFatek = protocol.includes('fatek')
  const isMCProtocol = protocol.includes('mcprotocol')

  /**
   * 獲取可用的操作選項
   */
  const getOperations = () => {
    if (isModbus) {
      return [
        { value: 'read_coils', label: '讀取線圈 (Coils)' },
        { value: 'read_discrete_inputs', label: '讀取離散輸入 (Discrete Inputs)' },
        { value: 'read_holding_registers', label: '讀取保持暫存器 (Holding Registers)' },
        { value: 'read_input_registers', label: '讀取輸入暫存器 (Input Registers)' },
      ]
    }
    if (isFatek) {
      return [
        { value: 'read_registers', label: '讀取暫存器' },
        { value: 'read_status', label: '讀取狀態' },
      ]
    }
    if (isMCProtocol) {
      return [
        { value: 'batch_read_word', label: '批量讀取字組' },
        { value: 'batch_read_bit', label: '批量讀取位元' },
      ]
    }
    return []
  }

  /**
   * 從 Profile 載入監控配置
   */
  useEffect(() => {
    if (currentProfile?.monitorConfig) {
      setMonitorItems(currentProfile.monitorConfig.items || [])
      setInterval(currentProfile.monitorConfig.interval || 1000)
    }
  }, [currentProfile]) // 當 profile 變更時載入

  /**
   * 保存監控配置到 Profile（使用防抖）
   */
  useEffect(() => {
    if (currentProfile && monitorItems.length > 0) {
      // 清除之前的定時器
      if (saveTimerRef.current) {
        clearTimeout(saveTimerRef.current)
      }
      // 設置新的定時器（防抖：1秒後保存）
      saveTimerRef.current = setTimeout(() => {
        updateProfile(currentProfile.id, {
          monitorConfig: {
            items: monitorItems,
            interval: interval,
          },
        })
      }, 1000)
    }
    
    return () => {
      if (saveTimerRef.current) {
        clearTimeout(saveTimerRef.current)
      }
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [monitorItems, interval, currentProfile, updateProfile])

  /**
   * SSE (Server-Sent Events) 連接處理
   */
  useEffect(() => {
    if (!connectionId) {
      // 斷線時清空數據並關閉連接
      if (eventSourceRef.current) {
        eventSourceRef.current.close()
        eventSourceRef.current = null
      }
      setMonitorData([])
      return
    }

    const host = window.location.host
    const sseUrl = `${window.location.protocol}//${host}/api/v1/test/monitor/stream?connection_id=${connectionId}`

    const eventSource = new EventSource(sseUrl)
    eventSourceRef.current = eventSource

    eventSource.onopen = () => {}

    // 監聽連接確認消息
    eventSource.addEventListener('connected', () => {})

    // 監聽心跳消息
    eventSource.addEventListener('ping', (_event: any) => {
      // 心跳消息，不需要處理
    })

    // 監聽默認消息（監控數據）
    eventSource.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data)

        // 檢查消息類型
        if (msg.type === 'monitor_update') {
          if (msg.connection_id === connectionId) {
            setMonitorData((prev) => {
              const newEntry = {
                chartTime: new Date(msg.timestamp).toLocaleTimeString(),
                timestamp: msg.timestamp,
                data: msg.data || {}, // 保存原始數據
              }
              return [newEntry, ...prev].slice(0, 100) // 保留最近 100 條記錄
            })
          } else {
            console.warn('✗ connection_id 不匹配，忽略消息')
          }
        }
      } catch (error) {
        console.error('解析監控數據失敗:', error, event.data)
      }
    }

    eventSource.onerror = (error) => {
      console.error('SSE error:', error, 'readyState:', eventSource.readyState)
      
      // EventSource 狀態：
      // 0 = CONNECTING
      // 1 = OPEN
      // 2 = CLOSED
      
      if (eventSource.readyState === EventSource.CLOSED) {
        reconnectAttemptsRef.current++
        const maxAttempts = 5
        const delay = Math.min(1000 * Math.pow(2, reconnectAttemptsRef.current - 1), 10000) // 指數退避，最多10秒
        
        if (reconnectAttemptsRef.current <= maxAttempts && connectionId) {
          console.warn(`SSE 連接已關閉，${delay/1000}秒後嘗試重連 (${reconnectAttemptsRef.current}/${maxAttempts})...`)
          
          // 清除之前的重連定時器
          if (reconnectTimerRef.current) {
            clearTimeout(reconnectTimerRef.current)
          }
          
          reconnectTimerRef.current = setTimeout(() => {
            if (connectionId && !eventSourceRef.current) {
              // 重新創建連接（通過重新執行 useEffect）
              // 這裡我們手動觸發重連
              const newEventSource = new EventSource(sseUrl)
              eventSourceRef.current = newEventSource
              
              newEventSource.onopen = () => {
                reconnectAttemptsRef.current = 0
              }
              
              newEventSource.addEventListener('connected', () => {})
              
              newEventSource.addEventListener('ping', () => {
                // 心跳消息
              })
              
              newEventSource.onmessage = (event) => {
                try {
                  const msg = JSON.parse(event.data)
                  if (msg.type === 'monitor_update' && msg.connection_id === connectionId) {
                    setMonitorData((prev) => {
                      const newEntry = {
                        ...msg,
                        chartTime: new Date(msg.timestamp).toLocaleTimeString(),
                        timestamp: msg.timestamp,
                        data: msg.data || {},
                      }
                      return [newEntry, ...prev].slice(0, 100)
                    })
                  }
                } catch (error) {
                  console.error('解析監控數據失敗:', error)
                }
              }
              
              newEventSource.onerror = (err) => {
                console.error('SSE 重連後錯誤:', err)
              }
            }
          }, delay)
        } else {
          console.error('SSE 重連失敗，已達到最大重試次數')
        }
      }
    }

    return () => {
      if (eventSourceRef.current) {
        eventSourceRef.current.close()
        eventSourceRef.current = null
      }
      if (reconnectTimerRef.current) {
        clearTimeout(reconnectTimerRef.current)
        reconnectTimerRef.current = null
      }
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [connectionId])

  /**
   * 添加監控項目
   */
  const handleAddItem = () => {
    const operations = getOperations()
    if (operations.length === 0) return

    const newItem: MonitorItem = {
      id: Date.now().toString(),
      operation: operations[0].value,
      address: 0,
      count: 1,
      label: `監控項目 ${monitorItems.length + 1}`,
      ...(isFatek && { symbol: 'D' }),
      ...(isMCProtocol && { device: 'D' }),
    }
    setMonitorItems([...monitorItems, newItem])
  }

  /**
   * 刪除監控項目
   */
  const handleRemoveItem = (id: string) => {
    setMonitorItems(monitorItems.filter(item => item.id !== id))
  }

  /**
   * 更新監控項目
   */
  const handleUpdateItem = (id: string, updates: Partial<MonitorItem>) => {
    setMonitorItems(monitorItems.map(item => 
      item.id === id ? { ...item, ...updates } : item
    ))
  }

  /**
   * 啟動監控
   */
  const handleStart = async () => {
    if (!connectionId) {
      showError('請先建立連線')
      return
    }

    if (monitorItems.length === 0) {
      showWarning('請至少添加一個監控項目')
      return
    }

    // 確保 SSE 連接已建立
    if (!eventSourceRef.current || eventSourceRef.current.readyState !== EventSource.OPEN) {
      console.warn('SSE 未連接，等待連接建立...')
      // 等待一下讓 SSE 連接建立
      await new Promise(resolve => setTimeout(resolve, 500))
      
      if (!eventSourceRef.current || eventSourceRef.current.readyState !== EventSource.OPEN) {
        console.warn('SSE 連接狀態:', eventSourceRef.current?.readyState)
        // SSE 會自動重連，繼續執行
      }
    }
    
    try {
      // 轉換為 ReadRequest 格式
      const items: ReadRequest[] = monitorItems.map(item => ({
        connection_id: connectionId,
        operation: item.operation,
        address: item.address,
        count: item.count,
        ...(item.symbol && { symbol: item.symbol }),
        ...(item.device && { device: item.device }),
      }))

      await startMonitor({
        connection_id: connectionId,
        interval: interval,
        items: items,
      })

      setMonitoring(true)
      setMonitorData([]) // 清空舊數據
    } catch (error: any) {
      console.error('啟動監控失敗:', error)
      showError('啟動監控失敗: ' + (error.message || error))
    }
  }

  /**
   * 停止監控
   */
  const handleStop = async () => {
    if (!connectionId) return
    try {
      await stopMonitor(connectionId)
      setMonitoring(false)
    } catch (error: any) {
      showError('停止監控失敗: ' + (error.message || error))
    }
  }

  /**
   * 準備圖表數據
   * 後端返回的數據結構：{ "item_0": { values: [...], count: ... }, "item_1": {...} }
   */
  const prepareChartData = () => {
    if (monitorData.length === 0 || monitorItems.length === 0) return []

    // 反轉數據（從舊到新）
    const reversed = [...monitorData].reverse()
    
    // 為每個監控項目創建數據點
    const chartData = reversed.map(entry => {
      const dataPoint: any = {
        time: entry.chartTime || new Date(entry.timestamp).toLocaleTimeString(),
      }
      
      // 為每個監控項目添加數據
      monitorItems.forEach((item, index) => {
        // 後端使用 item_0, item_1 等作為鍵名
        const backendKey = `item_${index}`
        const itemData = entry.data?.[backendKey]
        
        let value: number | null = null
        
        if (itemData) {
          // 處理錯誤情況
          if (typeof itemData === 'object' && itemData !== null && 'error' in itemData) {
            value = null // 錯誤時不顯示數據
          } else if (
            typeof itemData === 'object' &&
            itemData !== null &&
            'values' in itemData &&
            Array.isArray((itemData as { values?: unknown[] }).values)
          ) {
            // 如果有 values 數組，取第一個值（或平均值）
            const values = (itemData as { values: unknown[] }).values
            if (values.length > 0) {
              value = toNumericValue(values[0])
            }
          } else if (typeof itemData === 'number') {
            // 直接是數字
            value = toNumericValue(itemData)
          } else if (typeof itemData === 'string') {
            // 字符串，嘗試轉換
            value = toNumericValue(itemData)
          } else if (Array.isArray(itemData)) {
            // 直接是數組
            value = itemData.length > 0 ? toNumericValue(itemData[0]) : null
          }
        }
        
        // 使用 item.id 作為數據鍵名（用於圖表）
        dataPoint[`item_${item.id}`] = value
      })
      
      return dataPoint
    })

    return chartData
  }

  const chartData = prepareChartData()
  const operations = getOperations()

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between border-b dark:border-gray-700 pb-4">
        <h2 className="text-lg font-bold text-gray-800 dark:text-white flex items-center gap-2">
          <span className="w-1.5 h-6 bg-purple-500 rounded-full"></span>
          即時監控
        </h2>
        <div className="flex items-center gap-3">
          <div className="relative">
            <span className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 dark:text-gray-500 text-xs">間隔(ms)</span>
            <input
              type="number"
              value={interval}
              onChange={(e) => setInterval(Math.max(100, parseInt(e.target.value) || 1000))}
              disabled={monitoring}
              min={100}
              step={100}
              className="pl-16 pr-3 py-2 w-32 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-lg text-sm hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500 focus:bg-white dark:focus:bg-gray-800 text-gray-900 dark:text-white transition-all disabled:bg-gray-100 dark:disabled:bg-gray-700 disabled:cursor-not-allowed"
            />
          </div>
          
          {!monitoring ? (
            <button
              onClick={handleStart}
              disabled={!connectionId || monitorItems.length === 0}
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

      {/* 監控項目配置區域 */}
      {!monitoring && (
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-sm font-semibold text-gray-700 dark:text-gray-300">監控項目配置</h3>
            <button
              onClick={handleAddItem}
              disabled={!connectionId}
              className="text-sm text-purple-600 hover:text-purple-700 dark:text-purple-400 dark:hover:text-purple-300 font-medium flex items-center gap-1 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
              </svg>
              添加項目
            </button>
          </div>

          {monitorItems.length === 0 ? (
            <div className="text-center py-8 bg-gray-50 dark:bg-gray-900/50 rounded-lg border border-dashed border-gray-300 dark:border-gray-700 text-gray-400 dark:text-gray-500 text-sm">
              尚未添加監控項目，點擊「添加項目」開始配置
            </div>
          ) : (
            <div className="space-y-3">
              {monitorItems.map((item, index) => (
                <div key={item.id} className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-4 space-y-3 transition-colors">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <span className="text-xs font-semibold text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded">
                        項目 {index + 1}
                      </span>
                      {item.label && (
                        <span className="text-sm font-medium text-gray-700 dark:text-gray-300">{item.label}</span>
                      )}
                    </div>
                    <button
                      onClick={() => handleRemoveItem(item.id)}
                      className="text-red-500 hover:text-red-700 dark:text-red-400 dark:hover:text-red-300 p-1 transition-colors"
                    >
                      <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                      </svg>
                    </button>
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-3">
                    {/* 操作類型 */}
                    <div>
                      <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1">操作類型</label>
                      <select
                        value={item.operation}
                        onChange={(e) => handleUpdateItem(item.id, { operation: e.target.value })}
                        className="w-full px-3 py-2 text-sm bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white rounded-lg hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500 focus:bg-white dark:focus:bg-gray-700 transition-all"
                      >
                        {operations.map(op => (
                          <option key={op.value} value={op.value}>{op.label}</option>
                        ))}
                      </select>
                    </div>

                    {/* 地址 */}
                    <div>
                      <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1">地址</label>
                      <input
                        type="number"
                        value={item.address}
                        onChange={(e) => handleUpdateItem(item.id, { address: parseInt(e.target.value) || 0 })}
                        min={0}
                        className="w-full px-3 py-2 text-sm bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white rounded-lg hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500 focus:bg-white dark:focus:bg-gray-700 transition-all"
                      />
                    </div>

                    {/* 數量 */}
                    <div>
                      <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1">數量</label>
                      <input
                        type="number"
                        value={item.count}
                        onChange={(e) => handleUpdateItem(item.id, { count: Math.max(1, parseInt(e.target.value) || 1) })}
                        min={1}
                        className="w-full px-3 py-2 text-sm bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white rounded-lg hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500 focus:bg-white dark:focus:bg-gray-700 transition-all"
                      />
                    </div>

                    {/* Fatek Symbol 或 MC Device */}
                    {isFatek && (
                      <div>
                        <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1">符號 (Symbol)</label>
                        <input
                          type="text"
                          value={item.symbol || 'D'}
                          onChange={(e) => handleUpdateItem(item.id, { symbol: e.target.value })}
                          className="w-full px-3 py-2 text-sm bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white rounded-lg hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500 focus:bg-white dark:focus:bg-gray-700 transition-all"
                        />
                      </div>
                    )}

                    {isMCProtocol && (
                      <div>
                        <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1">設備 (Device)</label>
                        <input
                          type="text"
                          value={item.device || 'D'}
                          onChange={(e) => handleUpdateItem(item.id, { device: e.target.value })}
                          className="w-full px-3 py-2 text-sm bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white rounded-lg hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500 focus:bg-white dark:focus:bg-gray-700 transition-all"
                        />
                      </div>
                    )}

                    {/* 標籤（如果有額外空間） */}
                    {(!isFatek && !isMCProtocol) && (
                      <div>
                        <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1">標籤</label>
                        <input
                          type="text"
                          value={item.label || ''}
                          onChange={(e) => handleUpdateItem(item.id, { label: e.target.value })}
                          placeholder={`項目 ${index + 1}`}
                          className="w-full px-3 py-2 text-sm bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white rounded-lg hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500 focus:bg-white dark:focus:bg-gray-700 transition-all"
                        />
                      </div>
                    )}
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      {/* 監控運行時的顯示區域 */}
      {monitoring && (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 animate-fade-in">
          {/* 圖表區域 */}
          <div className="lg:col-span-2 bg-gray-50 dark:bg-gray-900/50 rounded-xl p-4 border border-gray-200 dark:border-gray-700 shadow-inner h-[400px]">
            {chartData.length === 0 ? (
              <div className="flex items-center justify-center h-full text-gray-400 dark:text-gray-500">
                <div className="text-center">
                  <svg className="w-12 h-12 mx-auto mb-2 opacity-50" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M7 12l3-3 3 3 4-4M8 21l4-4 4 4M3 4h18M4 4h16v12a1 1 0 01-1 1H5a1 1 0 01-1-1V4z" />
                  </svg>
                  <p className="text-sm">等待數據...</p>
                  <p className="text-xs mt-2 text-gray-300 dark:text-gray-600">
                    已接收 {monitorData.length} 條數據，配置 {monitorItems.length} 個監控項目
                  </p>
                </div>
              </div>
            ) : (
              <ResponsiveContainer width="100%" height="100%">
                <LineChart data={chartData}>
                  <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" className="dark:stroke-gray-700" />
                  <XAxis 
                    dataKey="time" 
                    tick={{fontSize: 10, fill: '#9ca3af'}} 
                    interval="preserveStartEnd"
                    stroke="#9ca3af"
                  />
                  <YAxis stroke="#9ca3af" tick={{fontSize: 10, fill: '#9ca3af'}} />
                  <Tooltip 
                    contentStyle={{ backgroundColor: '#1f2937', border: 'none', borderRadius: '8px', color: '#fff' }}
                  />
                  <Legend />
                  {monitorItems.map((item, index) => {
                    const colors = ['#8b5cf6', '#ec4899', '#10b981', '#f59e0b', '#3b82f6', '#ef4444']
                    const color = colors[index % colors.length]
                    return (
                      <Line 
                        key={item.id}
                        type="monotone" 
                        dataKey={`item_${item.id}`}
                        name={item.label || `項目 ${index + 1}`}
                        stroke={color}
                        strokeWidth={2} 
                        dot={false} 
                        activeDot={{ r: 6 }}
                        animationDuration={300}
                        connectNulls={false}
                      />
                    )
                  })}
                </LineChart>
              </ResponsiveContainer>
            )}
          </div>

          {/* 數據日誌區域 */}
          <div className="lg:col-span-1 bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 flex flex-col h-[400px] overflow-hidden transition-colors">
            <div className="bg-gray-50 dark:bg-gray-700 px-4 py-2 border-b border-gray-200 dark:border-gray-600 text-xs font-bold text-gray-500 dark:text-gray-300 uppercase flex items-center justify-between">
              <span>數據日誌 ({monitorData.length})</span>
              <button
                onClick={() => setMonitorData([])}
                className="text-xs text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 transition-colors"
              >
                清除
              </button>
            </div>
            <div className="flex-1 overflow-y-auto p-2 space-y-2 scrollbar-thin">
              {monitorData.length === 0 ? (
                <div className="text-center text-gray-400 dark:text-gray-500 text-xs py-8">
                  等待數據...
                </div>
              ) : (
                monitorData.map((entry, idx) => (
                  <div key={idx} className="text-xs p-2 rounded border border-gray-100 dark:border-gray-700 shadow-sm hover:shadow-md transition-shadow bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300">
                    <div className="flex justify-between text-gray-400 dark:text-gray-500 mb-1">
                      <span>#{idx + 1}</span>
                      <span>{new Date(entry.timestamp).toLocaleTimeString()}</span>
                    </div>
                    <div className="font-mono break-all space-y-1">
                      {Object.entries(entry.data || {}).map(([key, val]: [string, any]) => {
                        // 格式化顯示數據
                        let displayValue = ''
                        if (val && typeof val === 'object') {
                          if (val.error) {
                            displayValue = `錯誤: ${val.error}`
                          } else if (val.values && Array.isArray(val.values)) {
                            displayValue = `[${val.values.join(', ')}] (count: ${val.count || val.values.length})`
                          } else {
                            displayValue = JSON.stringify(val)
                          }
                        } else {
                          displayValue = String(val)
                        }
                        
                        // 找到對應的監控項目標籤
                        const itemIndex = parseInt(key.replace('item_', ''))
                        const item = monitorItems[itemIndex]
                        const label = item ? (item.label || `項目 ${itemIndex + 1}`) : key
                        
                        return (
                          <div key={key} className="flex items-start gap-2">
                            <span className="text-purple-600 dark:text-purple-400 font-semibold flex-shrink-0">{label}:</span>
                            <span className="break-all">{displayValue}</span>
                          </div>
                        )
                      })}
                    </div>
                  </div>
                ))
              )}
            </div>
          </div>
        </div>
      )}
      
      {/* 未啟動監控時的提示 */}
      {!monitoring && monitorItems.length > 0 && (
        <div className="flex flex-col items-center justify-center py-12 bg-gray-50 dark:bg-gray-900/50 rounded-xl border border-dashed border-gray-300 dark:border-gray-700 text-gray-400 dark:text-gray-500">
          <svg className="w-12 h-12 mb-2 opacity-50" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M7 12l3-3 3 3 4-4M8 21l4-4 4 4M3 4h18M4 4h16v12a1 1 0 01-1 1H5a1 1 0 01-1-1V4z" />
          </svg>
          <p className="text-sm">已配置 {monitorItems.length} 個監控項目，點擊「啟動監控」開始</p>
        </div>
      )}

      {!monitoring && monitorItems.length === 0 && (
        <div className="flex flex-col items-center justify-center py-12 bg-gray-50 dark:bg-gray-900/50 rounded-xl border border-dashed border-gray-300 dark:border-gray-700 text-gray-400 dark:text-gray-500">
          <svg className="w-12 h-12 mb-2 opacity-50" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M7 12l3-3 3 3 4-4M8 21l4-4 4 4M3 4h18M4 4h16v12a1 1 0 01-1 1H5a1 1 0 01-1-1V4z" />
          </svg>
          <p className="text-sm">準備就緒，請先添加監控項目</p>
        </div>
      )}
    </div>
  )
}
