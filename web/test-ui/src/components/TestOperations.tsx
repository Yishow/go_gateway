import { useState, useEffect, useRef, useCallback } from 'react'
import { useTestAPI } from '../services/api'
import type { BatchOperation } from '../types/api'

interface TestOperationsProps {
  protocol: string
  connectionId: string | null
}

/**
 * 輪詢結果記錄介面
 */
interface PollingRecord {
  timestamp: string
  data: any
  error?: string
}

export default function TestOperations({
  protocol,
  connectionId,
}: TestOperationsProps) {
  const [activeTab, setActiveTab] = useState<'single' | 'batch' | 'polling'>('single')
  const [operation, setOperation] = useState<string>('read_holding_registers')
  const [address, setAddress] = useState<number>(0)
  const [count, setCount] = useState<number>(10)
  const [symbol, setSymbol] = useState<string>('D')
  const [device, setDevice] = useState<string>('D')
  const [values, setValues] = useState<string>('')
  const [result, setResult] = useState<any>(null)
  const [loading, setLoading] = useState(false)
  
  // Batch state
  const [batchQueue, setBatchQueue] = useState<BatchOperation[]>([])
  
  // Polling state
  const [polling, setPolling] = useState<boolean>(false)
  const [pollingInterval, setPollingInterval] = useState<number>(1000)
  const [pollingRecords, setPollingRecords] = useState<PollingRecord[]>([])
  const pollingIntervalRef = useRef<number | null>(null)
  
  const { read, write, batch } = useTestAPI()

  const isModbus = protocol.includes('modbus')
  const isFatek = protocol.includes('fatek')
  const isMCProtocol = protocol.includes('mcprotocol')

  const isReadOp = operation.startsWith('read') || operation.includes('read')

  const getOperations = () => {
    if (isModbus) {
      return [
        { value: 'read_coils', label: '讀取線圈 (Coils)', type: 'read' },
        { value: 'read_discrete_inputs', label: '讀取離散輸入 (Discrete Inputs)', type: 'read' },
        { value: 'read_holding_registers', label: '讀取保持暫存器 (Holding Registers)', type: 'read' },
        { value: 'read_input_registers', label: '讀取輸入暫存器 (Input Registers)', type: 'read' },
        { value: 'write_single_coil', label: '寫入單個線圈', type: 'write' },
        { value: 'write_single_register', label: '寫入單個暫存器', type: 'write' },
        { value: 'write_multiple_coils', label: '寫入多個線圈', type: 'write' },
        { value: 'write_multiple_registers', label: '寫入多個暫存器', type: 'write' },
      ]
    }
    if (isFatek) {
      return [
        { value: 'read_registers', label: '讀取暫存器', type: 'read' },
        { value: 'read_status', label: '讀取狀態', type: 'read' },
        { value: 'write_registers', label: '寫入暫存器', type: 'write' },
        { value: 'write_status', label: '寫入狀態', type: 'write' },
      ]
    }
    if (isMCProtocol) {
      return [
        { value: 'batch_read_word', label: '批量讀取字組', type: 'read' },
        { value: 'batch_read_bit', label: '批量讀取位元', type: 'read' },
        { value: 'batch_write_word', label: '批量寫入字組', type: 'write' },
        { value: 'batch_write_bit', label: '批量寫入位元', type: 'write' },
      ]
    }
    return []
  }

  const handleExecute = async () => {
    if (!connectionId) {
      alert('請先建立連線')
      return
    }

    setLoading(true)
    try {
      if (isReadOp) {
        const data = await read(connectionId, {
          operation,
          address,
          count,
          symbol: isFatek ? symbol : undefined,
        })
        setResult(data)
      } else {
        const valuesArray = values.split(',').map((v) => {
          const trimmed = v.trim()
          if (trimmed === 'true') return true
          if (trimmed === 'false') return false
          const num = Number(trimmed)
          return isNaN(num) ? trimmed : num
        })
        await write(connectionId, {
          operation,
          address,
          values: valuesArray,
          symbol: isFatek ? symbol : undefined,
        })
        setResult({ status: 'success', message: '寫入成功' })
      }
    } catch (error: any) {
      setResult({ error: error.message })
    } finally {
      setLoading(false)
    }
  }

  const addToBatch = () => {
    const newOp: BatchOperation = {
      type: isReadOp ? 'read' : 'write',
      read_request: isReadOp ? {
        connection_id: connectionId!,
        operation,
        address,
        count,
        symbol: isFatek ? symbol : undefined,
        device: isMCProtocol ? device : undefined,
      } : undefined,
      write_request: !isReadOp ? {
        connection_id: connectionId!,
        operation,
        address,
        values: values.split(',').map(v => v.trim()),
        symbol: isFatek ? symbol : undefined,
        device: isMCProtocol ? device : undefined,
      } : undefined
    }
    setBatchQueue([...batchQueue, newOp])
  }

  const handleExecuteBatch = async () => {
    if (!connectionId || batchQueue.length === 0) return
    setLoading(true)
    try {
      const res = await batch({
        connection_id: connectionId,
        operations: batchQueue
      })
      setResult(res)
    } catch (error: any) {
      setResult({ error: error.message })
    } finally {
      setLoading(false)
    }
  }

  /**
   * 執行單次輪詢讀取
   * 使用 useCallback 確保總是使用最新的參數值
   */
  const executePollingRead = useCallback(async () => {
    if (!connectionId) return

    try {
      const data = await read(connectionId, {
        operation,
        address,
        count,
        symbol: isFatek ? symbol : undefined,
        device: isMCProtocol ? device : undefined,
      })
      
      const newRecord: PollingRecord = {
        timestamp: new Date().toLocaleTimeString('zh-TW', { 
          hour12: false,
          hour: '2-digit',
          minute: '2-digit',
          second: '2-digit',
          fractionalSecondDigits: 3
        }),
        data,
      }
      
      setPollingRecords(prev => {
        const updated = [newRecord, ...prev]
        // 最多保留 100 筆記錄
        return updated.slice(0, 100)
      })
    } catch (error: any) {
      const newRecord: PollingRecord = {
        timestamp: new Date().toLocaleTimeString('zh-TW', { 
          hour12: false,
          hour: '2-digit',
          minute: '2-digit',
          second: '2-digit',
          fractionalSecondDigits: 3
        }),
        data: null,
        error: error.message,
      }
      
      setPollingRecords(prev => {
        const updated = [newRecord, ...prev]
        return updated.slice(0, 100)
      })
    }
  }, [connectionId, operation, address, count, symbol, device, isFatek, isMCProtocol, read])

  /**
   * 啟動輪詢
   */
  const handleStartPolling = () => {
    if (!connectionId) {
      alert('請先建立連線')
      return
    }

    if (!isReadOp) {
      alert('輪詢功能僅支援讀取操作')
      return
    }

    setPolling(true)
    // 立即執行一次
    executePollingRead()
  }

  /**
   * 停止輪詢
   */
  const handleStopPolling = () => {
    setPolling(false)
  }

  /**
   * 清除輪詢記錄
   */
  const handleClearPollingRecords = () => {
    setPollingRecords([])
  }

  /**
   * 輪詢效果處理
   */
  useEffect(() => {
    if (polling && connectionId) {
      pollingIntervalRef.current = setInterval(() => {
        executePollingRead()
      }, pollingInterval)
    } else {
      if (pollingIntervalRef.current) {
        clearInterval(pollingIntervalRef.current)
        pollingIntervalRef.current = null
      }
    }

    return () => {
      if (pollingIntervalRef.current) {
        clearInterval(pollingIntervalRef.current)
      }
    }
  }, [polling, pollingInterval, executePollingRead])

  /**
   * 切換 tab 時停止輪詢
   */
  useEffect(() => {
    if (activeTab !== 'polling' && polling) {
      setPolling(false)
    }
  }, [activeTab, polling])

  /**
   * 協議切換時重置操作類型為該協議的第一個讀取操作
   */
  useEffect(() => {
    const operations = getOperations()
    if (operations.length === 0) return
    
    const firstReadOp = operations.find(op => op.type === 'read')
    if (firstReadOp) {
      // 檢查當前操作是否在新協議中有效
      const isValid = operations.some(op => op.value === operation)
      if (!isValid) {
        setOperation(firstReadOp.value)
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [protocol])

  const InputField = ({ label, children }: { label: string, children: React.ReactNode }) => (
    <div className="space-y-1">
      <label className="text-sm font-medium text-gray-700">{label}</label>
      {children}
    </div>
  )

  const StyledInput = (props: React.InputHTMLAttributes<HTMLInputElement>) => (
    <input
      {...props}
      className="w-full px-3 py-2 bg-white border border-gray-300 rounded-lg shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 focus:bg-white transition-all duration-200 disabled:bg-gray-100 disabled:cursor-not-allowed"
    />
  )

  return (
    <div className="space-y-6">
      {/* Tabs */}
      <div className="flex bg-gray-100 p-1 rounded-lg">
        <button
          className={`flex-1 py-2 px-4 rounded-md text-sm font-medium transition-all duration-200 ${activeTab === 'single' ? 'bg-white text-blue-700 shadow-sm' : 'text-gray-500 hover:text-gray-700'}`}
          onClick={() => setActiveTab('single')}
        >
          單次操作
        </button>
        <button
          className={`flex-1 py-2 px-4 rounded-md text-sm font-medium transition-all duration-200 ${activeTab === 'batch' ? 'bg-white text-blue-700 shadow-sm' : 'text-gray-500 hover:text-gray-700'}`}
          onClick={() => setActiveTab('batch')}
        >
          批量操作 {batchQueue.length > 0 && <span className="ml-1 px-1.5 py-0.5 bg-blue-100 text-blue-600 rounded-full text-xs">{batchQueue.length}</span>}
        </button>
        <button
          className={`flex-1 py-2 px-4 rounded-md text-sm font-medium transition-all duration-200 ${activeTab === 'polling' ? 'bg-white text-blue-700 shadow-sm' : 'text-gray-500 hover:text-gray-700'}`}
          onClick={() => setActiveTab('polling')}
        >
          輪詢模式 {polling && <span className="ml-1 px-1.5 py-0.5 bg-red-100 text-red-600 rounded-full text-xs animate-pulse">運行中</span>}
        </button>
      </div>

      <div className="space-y-4 animate-fade-in">
        <InputField label="操作類型">
          <div className="relative">
             <select
              value={operation}
              onChange={(e) => setOperation(e.target.value)}
              className="w-full px-4 py-2 appearance-none border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-white"
            >
              {getOperations().map((op) => (
                <option key={op.value} value={op.value}>
                  {op.label}
                </option>
              ))}
            </select>
             <div className="absolute inset-y-0 right-0 flex items-center px-2 pointer-events-none text-gray-500">
              <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
              </svg>
            </div>
          </div>
        </InputField>

        <div className="grid grid-cols-2 gap-4">
          {isFatek && (
            <InputField label="組件符號">
              <StyledInput
                type="text"
                value={symbol}
                onChange={(e) => setSymbol(e.target.value)}
                placeholder="D, M, X, Y"
              />
            </InputField>
          )}

          {isMCProtocol && (
            <InputField label="設備代號">
              <StyledInput
                type="text"
                value={device}
                onChange={(e) => setDevice(e.target.value)}
                placeholder="D, M, X, Y"
              />
            </InputField>
          )}

          <InputField label="起始位址">
            <StyledInput
              type="number"
              value={address}
              onChange={(e) => setAddress(parseInt(e.target.value) || 0)}
            />
          </InputField>

          {isReadOp && (
            <InputField label="讀取數量">
              <StyledInput
                type="number"
                value={count}
                onChange={(e) => setCount(parseInt(e.target.value) || 1)}
              />
            </InputField>
          )}
        </div>

        {!isReadOp && (
          <InputField label="寫入數值 (逗號分隔)">
            <StyledInput
              type="text"
              value={values}
              onChange={(e) => setValues(e.target.value)}
              placeholder="1, 2, 3 或 true, false"
            />
          </InputField>
        )}

        {activeTab === 'single' ? (
          <button
            onClick={handleExecute}
            disabled={loading || !connectionId}
            className={`w-full py-3 rounded-xl font-semibold shadow-lg transition-all active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed
              ${isReadOp 
                ? 'bg-blue-600 text-white shadow-blue-500/30 hover:bg-blue-700 hover:shadow-blue-500/40' 
                : 'bg-amber-500 text-white shadow-amber-500/30 hover:bg-amber-600 hover:shadow-amber-500/40'
              }
            `}
          >
            {loading ? '執行中...' : `執行${isReadOp ? '讀取' : '寫入'}`}
          </button>
        ) : activeTab === 'batch' ? (
          <div className="flex gap-2">
            <button
              onClick={addToBatch}
              disabled={!connectionId}
              className="flex-1 bg-white border border-gray-300 text-gray-700 px-4 py-3 rounded-xl font-medium hover:bg-gray-50 active:scale-[0.98] transition-all disabled:opacity-50"
            >
              加入清單
            </button>
            <button
              onClick={handleExecuteBatch}
              disabled={loading || !connectionId || batchQueue.length === 0}
              className="flex-1 bg-green-600 text-white px-4 py-3 rounded-xl font-semibold shadow-lg shadow-green-500/30 hover:bg-green-700 hover:shadow-green-500/40 active:scale-[0.98] transition-all disabled:opacity-50"
            >
              {loading ? '批次執行中...' : '執行批量'}
            </button>
            <button
              onClick={() => setBatchQueue([])}
              className="px-4 py-3 text-red-600 hover:bg-red-50 rounded-xl transition-all"
              title="清除清單"
            >
              <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </button>
          </div>
        ) : (
          <div className="space-y-4">
            <InputField label="輪詢間隔 (毫秒)">
              <StyledInput
                type="number"
                value={pollingInterval}
                onChange={(e) => setPollingInterval(Math.max(100, parseInt(e.target.value) || 1000))}
                disabled={polling}
                min={100}
                step={100}
              />
            </InputField>
            <div className="flex gap-2">
              {!polling ? (
                <button
                  onClick={handleStartPolling}
                  disabled={!connectionId || !isReadOp}
                  className="flex-1 bg-purple-600 text-white px-4 py-3 rounded-xl font-semibold shadow-lg shadow-purple-500/30 hover:bg-purple-700 hover:shadow-purple-500/40 active:scale-[0.98] transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  啟動輪詢
                </button>
              ) : (
                <button
                  onClick={handleStopPolling}
                  className="flex-1 bg-red-600 text-white px-4 py-3 rounded-xl font-semibold shadow-lg shadow-red-500/30 hover:bg-red-700 hover:shadow-red-500/40 active:scale-[0.98] transition-all flex items-center justify-center gap-2"
                >
                  <span className="animate-pulse w-2 h-2 rounded-full bg-white"></span>
                  停止輪詢
                </button>
              )}
              <button
                onClick={handleClearPollingRecords}
                disabled={pollingRecords.length === 0}
                className="px-4 py-3 text-gray-600 hover:bg-gray-50 rounded-xl transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                title="清除記錄"
              >
                <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>
            {!isReadOp && (
              <div className="bg-amber-50 border border-amber-200 rounded-lg p-3 text-sm text-amber-800">
                <span className="font-medium">提示：</span>輪詢功能僅支援讀取操作，請選擇讀取類型的操作。
              </div>
            )}
          </div>
        )}
      </div>

      {activeTab === 'batch' && batchQueue.length > 0 && (
        <div className="bg-gray-50 p-4 rounded-xl border border-dashed border-gray-300 animate-fade-in">
          <h4 className="text-xs font-bold text-gray-500 uppercase mb-3">待執行項目</h4>
          <ul className="text-sm space-y-2 max-h-40 overflow-y-auto scrollbar-thin">
            {batchQueue.map((op, idx) => (
              <li key={idx} className="flex items-center gap-2 text-gray-600 bg-white p-2 rounded border border-gray-100 shadow-sm">
                <span className={`w-2 h-2 rounded-full ${op.type === 'read' ? 'bg-blue-500' : 'bg-amber-500'}`}></span>
                <span className="font-mono text-xs text-gray-400">{idx + 1}.</span>
                <span className="font-medium">{op.type === 'read' ? op.read_request?.operation : op.write_request?.operation}</span>
                <span className="text-gray-400 text-xs ml-auto">@ {op.type === 'read' ? op.read_request?.address : op.write_request?.address}</span>
              </li>
            ))}
          </ul>
        </div>
      )}

      {activeTab === 'polling' && pollingRecords.length > 0 && (
        <div className="bg-gray-50 p-4 rounded-xl border border-dashed border-gray-300 animate-fade-in">
          <div className="flex justify-between items-center mb-3">
            <h4 className="text-xs font-bold text-gray-500 uppercase">輪詢記錄 ({pollingRecords.length})</h4>
            <button
              onClick={handleClearPollingRecords}
              className="text-xs text-gray-500 hover:text-gray-700 font-medium"
            >
              清除
            </button>
          </div>
          <div className="max-h-64 overflow-y-auto scrollbar-thin space-y-2">
            {pollingRecords.map((record, idx) => (
              <div
                key={idx}
                className={`bg-white p-3 rounded-lg border shadow-sm ${
                  record.error ? 'border-red-200 bg-red-50' : 'border-gray-100'
                }`}
              >
                <div className="flex items-center justify-between mb-2">
                  <span className="text-xs font-mono text-gray-500">{record.timestamp}</span>
                  {record.error ? (
                    <span className="text-xs text-red-600 font-medium">錯誤</span>
                  ) : (
                    <span className="text-xs text-green-600 font-medium">成功</span>
                  )}
                </div>
                {record.error ? (
                  <div className="text-sm text-red-700 font-mono">{record.error}</div>
                ) : (
                  <div className="text-xs font-mono text-gray-700 overflow-x-auto">
                    <div className="space-y-1">
                      {Array.isArray(record.data?.values) ? (
                        <div className="flex flex-wrap gap-1">
                          {record.data.values.map((val: any, i: number) => (
                            <span
                              key={i}
                              className="px-2 py-1 bg-blue-100 text-blue-800 rounded text-xs"
                            >
                              [{i}]: {String(val)}
                            </span>
                          ))}
                        </div>
                      ) : (
                        <pre className="text-xs">{JSON.stringify(record.data, null, 2)}</pre>
                      )}
                    </div>
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>
      )}

      {result && (
        <div className="mt-6 animate-slide-up">
          <div className="flex justify-between items-center mb-2">
            <label className="text-xs font-bold text-gray-500 uppercase tracking-wider">執行結果</label>
            <button onClick={() => setResult(null)} className="text-xs text-blue-500 hover:text-blue-700 font-medium">清除</button>
          </div>
          <div className="relative group">
            <div className="absolute inset-0 bg-gray-900 rounded-lg shadow-inner"></div>
            <div className="relative p-4 rounded-lg overflow-hidden">
               <pre className="text-xs font-mono text-green-400 overflow-auto max-h-64 scrollbar-thin">
                {JSON.stringify(result, null, 2)}
              </pre>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
