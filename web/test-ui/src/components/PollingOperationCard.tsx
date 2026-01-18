import { useState, useEffect, useRef, useCallback } from 'react'
import { useTestAPI } from '../services/api'
import { useToast } from '../contexts/ToastContext'

/**
 * 單個 Polling 操作的配置
 */
export interface PollingOperation {
  id: string
  operation: string
  address: number
  count: number
  values?: string // 寫入操作的數值（逗號分隔）
  symbol?: string
  device?: string
  unitId?: string
  station?: string
  interval: number
  baudRate?: number // 波特率（每個操作獨立）
  enabled: boolean
}

/**
 * Polling 記錄
 */
export interface PollingRecord {
  timestamp: string
  data: any
  error?: string
}

interface PollingOperationCardProps {
  operation: PollingOperation
  onUpdate: (operation: PollingOperation) => void
  onDelete: () => void
  connectionId: string | null
  protocol: string
  isModbus: boolean
  isFatek: boolean
  isMCProtocol: boolean
  getOperations: () => Array<{ value: string; label: string; type: string }>
  isReadOperation: (operation: string) => boolean
}

/**
 * 單個 Polling 操作卡片組件
 */
export default function PollingOperationCard({
  operation,
  onUpdate,
  onDelete,
  connectionId,
  protocol: _protocol,
  isModbus,
  isFatek,
  isMCProtocol,
  getOperations,
  isReadOperation,
}: PollingOperationCardProps) {
  const [records, setRecords] = useState<PollingRecord[]>([])
  const intervalRef = useRef<number | null>(null)
  const operationRef = useRef(operation)
  
  const { read, write } = useTestAPI()
  const { showError } = useToast()

  // 更新 ref
  useEffect(() => {
    operationRef.current = operation
  }, [operation])

  /**
   * 執行單次 polling 讀取
   */
  const executePollingRead = useCallback(async () => {
    if (!connectionId) return

    try {
      const readParams: any = {
        operation: operationRef.current.operation,
        address: operationRef.current.address,
        count: operationRef.current.count,
        symbol: isFatek ? operationRef.current.symbol : undefined,
        device: isMCProtocol ? operationRef.current.device : undefined,
      }
      
      // 添加站號參數
      if (isModbus && operationRef.current.unitId) {
        const unitIdNum = parseInt(operationRef.current.unitId)
        if (!isNaN(unitIdNum)) {
          readParams.unit_id = unitIdNum
        }
      }
      if (isFatek && operationRef.current.station) {
        const stationNum = parseInt(operationRef.current.station)
        if (!isNaN(stationNum)) {
          readParams.station = stationNum
        }
      }
      
      const data = await read(connectionId, readParams)
      
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
      
      setRecords(prev => {
        const updated = [newRecord, ...prev]
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
      
      setRecords(prev => {
        const updated = [newRecord, ...prev]
        return updated.slice(0, 100)
      })
    }
  }, [connectionId, isFatek, isMCProtocol, isModbus, read])

  /**
   * 執行單次 polling 寫入
   */
  const executePollingWrite = useCallback(async () => {
    if (!connectionId) return

    try {
      if (!operationRef.current.values || operationRef.current.values.trim() === '') {
        throw new Error('請輸入寫入數值')
      }

      const valuesArray = operationRef.current.values.split(',').map((v) => {
        const trimmed = v.trim()
        if (trimmed === 'true') return true
        if (trimmed === 'false') return false
        const num = Number(trimmed)
        return isNaN(num) ? trimmed : num
      })

      const writeParams: any = {
        operation: operationRef.current.operation,
        address: operationRef.current.address,
        values: valuesArray,
        symbol: isFatek ? operationRef.current.symbol : undefined,
        device: isMCProtocol ? operationRef.current.device : undefined,
      }
      
      // 添加站號參數
      if (isModbus && operationRef.current.unitId) {
        const unitIdNum = parseInt(operationRef.current.unitId)
        if (!isNaN(unitIdNum)) {
          writeParams.unit_id = unitIdNum
        }
      }
      if (isFatek && operationRef.current.station) {
        const stationNum = parseInt(operationRef.current.station)
        if (!isNaN(stationNum)) {
          writeParams.station = stationNum
        }
      }
      
      await write(connectionId, writeParams)
      
      const newRecord: PollingRecord = {
        timestamp: new Date().toLocaleTimeString('zh-TW', { 
          hour12: false,
          hour: '2-digit',
          minute: '2-digit',
          second: '2-digit',
          fractionalSecondDigits: 3
        }),
        data: { status: 'success', message: '寫入成功', values: valuesArray },
      }
      
      setRecords(prev => {
        const updated = [newRecord, ...prev]
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
      
      setRecords(prev => {
        const updated = [newRecord, ...prev]
        return updated.slice(0, 100)
      })
    }
  }, [connectionId, isFatek, isMCProtocol, isModbus, write])

  /**
   * 執行單次 polling 操作（讀取或寫入）
   */
  const executePollingOperation = useCallback(async () => {
    if (isReadOperation(operationRef.current.operation)) {
      await executePollingRead()
    } else {
      await executePollingWrite()
    }
  }, [executePollingRead, executePollingWrite, isReadOperation])

  /**
   * 啟動 Polling
   */
  const handleStart = () => {
    if (!connectionId) {
      showError('請先建立連線')
      return
    }
    
    onUpdate({ ...operation, enabled: true })
    // 立即執行一次
    executePollingOperation()
  }

  /**
   * 停止 Polling
   */
  const handleStop = () => {
    onUpdate({ ...operation, enabled: false })
  }

  /**
   * 清除記錄
   */
  const handleClearRecords = () => {
    setRecords([])
  }

  /**
   * Polling 效果處理
   */
  useEffect(() => {
    if (operation.enabled && connectionId) {
      intervalRef.current = setInterval(() => {
        executePollingOperation()
      }, operation.interval)
    } else {
      if (intervalRef.current) {
        clearInterval(intervalRef.current)
        intervalRef.current = null
      }
    }

    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current)
      }
    }
  }, [operation.enabled, operation.interval, connectionId, executePollingOperation])

  /**
   * 當連線斷開時，停止 polling
   */
  useEffect(() => {
    if (!connectionId && operation.enabled) {
      onUpdate({ ...operation, enabled: false })
    }
  }, [connectionId, operation.enabled, onUpdate]) // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <div className="border border-gray-200 dark:border-gray-700 rounded-lg p-2.5 space-y-2">
      <div className="flex items-center justify-between">
        <h3 className="text-xs font-semibold text-gray-700 dark:text-gray-300">
          Polling 操作 {operation.id}
        </h3>
        <div className="flex items-center gap-2">
          {operation.enabled && (
            <span className="px-1.5 py-0.5 bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400 rounded text-xs font-medium animate-pulse">
              運行中
            </span>
          )}
          <button
            onClick={onDelete}
            className="px-1.5 py-1 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded transition-all"
            title="刪除操作"
          >
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>

      {/* 操作類型、起始位址、讀取數量/寫入數值、間隔、波特率 - 同一排 */}
      <div className="flex items-end gap-2 flex-wrap">
        <div className="flex-shrink-0">
          <label className="text-xs font-medium text-gray-600 dark:text-gray-400 mb-0.5 block">操作類型</label>
          <div className="relative">
            <select
              value={operation.operation}
              onChange={(e) => onUpdate({ ...operation, operation: e.target.value })}
              disabled={operation.enabled}
              className="px-2 py-1 text-xs appearance-none border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-800 text-gray-900 dark:text-white disabled:opacity-50 disabled:cursor-not-allowed w-40"
            >
              {getOperations().map((op) => (
                <option key={op.value} value={op.value}>
                  {op.label}
                </option>
              ))}
            </select>
          </div>
        </div>

        <div className="flex-shrink-0">
          <label className="text-xs font-medium text-gray-600 dark:text-gray-400 mb-0.5 block">波特率</label>
          <div className="relative">
            <select
              value={operation.baudRate || 9600}
              onChange={(e) => onUpdate({ ...operation, baudRate: parseInt(e.target.value) })}
              disabled={operation.enabled}
              className="px-2 py-1 text-xs appearance-none border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-800 text-gray-900 dark:text-white disabled:opacity-50 disabled:cursor-not-allowed w-20"
            >
              <option value={9600}>9600</option>
              <option value={19200}>19200</option>
              <option value={38400}>38400</option>
              <option value={57600}>57600</option>
              <option value={115200}>115200</option>
            </select>
          </div>
        </div>

        <div className="flex-shrink-0">
          <label className="text-xs font-medium text-gray-600 dark:text-gray-400 mb-0.5 block">起始位址</label>
          <input
            type="number"
            value={operation.address}
            onChange={(e) => onUpdate({ ...operation, address: parseInt(e.target.value) || 0 })}
            disabled={operation.enabled}
            className="w-20 px-2 py-1 text-xs border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-800 text-gray-900 dark:text-white disabled:opacity-50 disabled:cursor-not-allowed"
          />
        </div>

        {isReadOperation(operation.operation) ? (
          <div className="flex-shrink-0">
            <label className="text-xs font-medium text-gray-600 dark:text-gray-400 mb-0.5 block">讀取數量</label>
            <input
              type="number"
              value={operation.count}
              onChange={(e) => onUpdate({ ...operation, count: parseInt(e.target.value) || 1 })}
              disabled={operation.enabled}
              className="w-20 px-2 py-1 text-xs border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-800 text-gray-900 dark:text-white disabled:opacity-50 disabled:cursor-not-allowed"
            />
          </div>
        ) : (
          <div className="flex-shrink-0">
            <label className="text-xs font-medium text-gray-600 dark:text-gray-400 mb-0.5 block">寫入數值</label>
            <input
              type="text"
              value={operation.values || ''}
              onChange={(e) => onUpdate({ ...operation, values: e.target.value })}
              disabled={operation.enabled}
              placeholder="1,2,3 或 true,false"
              className="w-32 px-2 py-1 text-xs border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-800 text-gray-900 dark:text-white disabled:opacity-50 disabled:cursor-not-allowed"
            />
          </div>
        )}

        <div className="flex-shrink-0">
          <label className="text-xs font-medium text-gray-600 dark:text-gray-400 mb-0.5 block">間隔 (ms)</label>
          <input
            type="number"
            value={operation.interval}
            onChange={(e) => onUpdate({ ...operation, interval: Math.max(100, parseInt(e.target.value) || 1000) })}
            disabled={operation.enabled}
            min={100}
            step={100}
            className="w-20 px-2 py-1 text-xs border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-800 text-gray-900 dark:text-white disabled:opacity-50 disabled:cursor-not-allowed"
          />
        </div>

        {(isFatek || isMCProtocol) && (
          <div className="flex-shrink-0">
            <label className="text-xs font-medium text-gray-600 dark:text-gray-400 mb-0.5 block">
              {isFatek ? '組件符號' : '設備代號'}
            </label>
            <input
              type="text"
              value={operation.symbol || operation.device || ''}
              onChange={(e) => {
                if (isFatek) {
                  onUpdate({ ...operation, symbol: e.target.value })
                } else {
                  onUpdate({ ...operation, device: e.target.value })
                }
              }}
              disabled={operation.enabled}
              placeholder={isFatek ? 'D' : 'D'}
              className="w-16 px-2 py-1 text-xs border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-800 text-gray-900 dark:text-white disabled:opacity-50 disabled:cursor-not-allowed"
            />
          </div>
        )}

        {isModbus && (
          <div className="flex-shrink-0">
            <label className="text-xs font-medium text-gray-600 dark:text-gray-400 mb-0.5 block">站號</label>
            <input
              type="number"
              value={operation.unitId || ''}
              onChange={(e) => onUpdate({ ...operation, unitId: e.target.value })}
              disabled={operation.enabled}
              placeholder="留空"
              min={1}
              max={255}
              className="w-16 px-2 py-1 text-xs border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-800 text-gray-900 dark:text-white disabled:opacity-50 disabled:cursor-not-allowed"
            />
          </div>
        )}

        {isFatek && (
          <div className="flex-shrink-0">
            <label className="text-xs font-medium text-gray-600 dark:text-gray-400 mb-0.5 block">站號</label>
            <input
              type="number"
              value={operation.station || ''}
              onChange={(e) => onUpdate({ ...operation, station: e.target.value })}
              disabled={operation.enabled}
              placeholder="留空"
              min={1}
              max={255}
              className="w-16 px-2 py-1 text-xs border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-800 text-gray-900 dark:text-white disabled:opacity-50 disabled:cursor-not-allowed"
            />
          </div>
        )}
      </div>

      {/* 最新回傳值顯示 */}
      {records.length > 0 && (
        <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-md p-1.5">
          <div className="flex items-center justify-between mb-1">
            <span className="text-xs font-semibold text-blue-700 dark:text-blue-400">最新回傳值</span>
            <span className="text-xs font-mono text-blue-600 dark:text-blue-500">{records[0].timestamp}</span>
          </div>
          {records[0].error ? (
            <div className="text-xs text-red-700 dark:text-red-300 font-mono">{records[0].error}</div>
          ) : (
            <div className="text-xs font-mono text-gray-700 dark:text-gray-300">
              {records[0].data?.status === 'success' ? (
                <div className="text-green-700 dark:text-green-300">
                  {records[0].data.message}
                  {Array.isArray(records[0].data.values) && (
                    <div className="flex flex-wrap gap-1 mt-1">
                      {records[0].data.values.map((val: any, i: number) => (
                        <span
                          key={i}
                          className="px-1.5 py-0.5 bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-300 rounded text-xs"
                        >
                          [{i}]: {String(val)}
                        </span>
                      ))}
                    </div>
                  )}
                </div>
              ) : Array.isArray(records[0].data?.values) ? (
                <div className="flex flex-wrap gap-1">
                  {records[0].data.values.map((val: any, i: number) => (
                    <span
                      key={i}
                      className="px-1.5 py-0.5 bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-300 rounded text-xs"
                    >
                      [{i}]: {String(val)}
                    </span>
                  ))}
                </div>
              ) : (
                <pre className="text-xs whitespace-pre-wrap break-words">{JSON.stringify(records[0].data, null, 2)}</pre>
              )}
            </div>
          )}
        </div>
      )}

      {/* 開始/停止按鈕和清除按鈕 - 同一排 */}
      <div className="flex items-center gap-2">
        {!operation.enabled ? (
          <button
            onClick={handleStart}
            disabled={!connectionId}
            className="bg-purple-600 text-white px-3 py-1 rounded-md font-semibold shadow-md shadow-purple-500/30 hover:bg-purple-700 hover:shadow-purple-500/40 active:scale-[0.98] transition-all disabled:opacity-50 disabled:cursor-not-allowed text-xs"
          >
            開始
          </button>
        ) : (
          <button
            onClick={handleStop}
            className="bg-red-600 text-white px-3 py-1 rounded-md font-semibold shadow-md shadow-red-500/30 hover:bg-red-700 hover:shadow-red-500/40 active:scale-[0.98] transition-all flex items-center gap-1.5 text-xs"
          >
            <span className="animate-pulse w-1.5 h-1.5 rounded-full bg-white"></span>
            停止
          </button>
        )}
        <button
          onClick={handleClearRecords}
          disabled={records.length === 0}
          className="px-2 py-1 text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-700 rounded-md transition-all disabled:opacity-50 disabled:cursor-not-allowed"
          title="清除記錄"
        >
          <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
        </button>
      </div>

      {/* 記錄顯示 */}
      {records.length > 1 && (
        <div className="bg-gray-50 dark:bg-gray-900/50 p-2 rounded-md border border-dashed border-gray-300 dark:border-gray-700">
          <div className="flex justify-between items-center mb-1">
            <h4 className="text-xs font-bold text-gray-500 dark:text-gray-400 uppercase">
              歷史記錄 ({records.length - 1})
            </h4>
          </div>
          <div className="max-h-24 overflow-y-auto scrollbar-thin space-y-0.5">
            {records.slice(1, 4).map((record, idx) => (
              <div
                key={idx}
                className={`text-xs p-1.5 rounded border ${
                  record.error 
                    ? 'border-red-200 dark:border-red-800 bg-red-50 dark:bg-red-900/20' 
                    : 'border-gray-100 dark:border-gray-700 bg-white dark:bg-gray-800'
                }`}
              >
                <div className="flex items-center justify-between">
                  <span className="font-mono text-gray-500 dark:text-gray-400 text-xs">{record.timestamp}</span>
                  {record.error ? (
                    <span className="text-red-600 dark:text-red-400 font-medium text-xs">錯誤</span>
                  ) : (
                    <span className="text-green-600 dark:text-green-400 font-medium text-xs">成功</span>
                  )}
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}
