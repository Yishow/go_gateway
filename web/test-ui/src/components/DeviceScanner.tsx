import { useState, useCallback, useRef, useEffect } from 'react'
import { useTestAPI } from '../services/api'

interface DeviceScannerProps {
  protocol: string
  connectionId: string | null
}

/**
 * 掃描結果介面
 */
interface ScanResult {
  address: number
  success: boolean
  data?: any
  error?: string
}

/**
 * 設備掃描器組件
 * 支援輸入位置，自動修改相關參數，掃描連線底下的設備數量
 * 適配所有協議：Modbus、FATEK、MC Protocol
 */
export default function DeviceScanner({
  protocol,
  connectionId,
}: DeviceScannerProps) {
  const [startAddress, setStartAddress] = useState<number>(0)
  const [maxAddress, setMaxAddress] = useState<number>(100)
  const [scanStep, setScanStep] = useState<number>(1)
  const [scanning, setScanning] = useState<boolean>(false)
  const [scanResults, setScanResults] = useState<ScanResult[]>([])
  const [deviceCount, setDeviceCount] = useState<number>(0)
  
  // 協議特定參數
  const [symbol, setSymbol] = useState<string>('D') // FATEK
  const [device, setDevice] = useState<string>('D') // MC Protocol
  const [operation, setOperation] = useState<string>('')
  
  const { read } = useTestAPI()
  const scanAbortRef = useRef<boolean>(false)

  const isModbus = protocol.includes('modbus')
  const isFatek = protocol.includes('fatek')
  const isMCProtocol = protocol.includes('mcprotocol') || protocol.includes('mc_')

  /**
   * 根據協議獲取預設操作類型
   */
  const getDefaultOperation = useCallback(() => {
    if (isModbus) {
      return 'read_holding_registers'
    }
    if (isFatek) {
      return 'read_registers'
    }
    if (isMCProtocol) {
      return 'batch_read_word'
    }
    return ''
  }, [isModbus, isFatek, isMCProtocol])

  /**
   * 獲取協議支援的操作列表
   */
  const getAvailableOperations = () => {
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
   * 執行單次讀取掃描
   */
  const scanAddress = useCallback(async (address: number): Promise<ScanResult> => {
    if (!connectionId) {
      return { address, success: false, error: '未連線' }
    }

    try {
      const readParams: any = {
        operation,
        address,
        count: 1,
      }

      // 根據協議添加特定參數
      if (isFatek) {
        readParams.symbol = symbol
      }
      if (isMCProtocol) {
        readParams.device = device
      }

      const result = await read(connectionId, readParams)
      
      return {
        address,
        success: true,
        data: result,
      }
    } catch (error: any) {
      return {
        address,
        success: false,
        error: error.message || '讀取失敗',
      }
    }
  }, [connectionId, operation, isFatek, isMCProtocol, symbol, device, read])

  /**
   * 執行掃描
   */
  const handleScan = async () => {
    if (!connectionId) {
      alert('請先建立連線')
      return
    }

    if (!operation) {
      alert('請選擇操作類型')
      return
    }

    setScanning(true)
    setScanResults([])
    setDeviceCount(0)
    scanAbortRef.current = false

    const results: ScanResult[] = []
    let successCount = 0

    try {
      // 從起始地址開始，逐步掃描到最大地址
      for (let addr = startAddress; addr <= maxAddress && !scanAbortRef.current; addr += scanStep) {
        const result = await scanAddress(addr)
        results.push(result)
        
        if (result.success) {
          successCount++
        }

        // 更新結果（實時顯示）
        setScanResults([...results])
        setDeviceCount(successCount)

        // 添加小延遲，避免過快掃描
        await new Promise(resolve => setTimeout(resolve, 50))
      }
    } catch (error: any) {
      console.error('掃描錯誤:', error)
    } finally {
      setScanning(false)
      scanAbortRef.current = false
    }
  }

  /**
   * 停止掃描
   */
  const handleStopScan = () => {
    scanAbortRef.current = true
    setScanning(false)
  }

  /**
   * 清除結果
   */
  const handleClearResults = () => {
    setScanResults([])
    setDeviceCount(0)
  }

  /**
   * 當協議變更時，自動更新操作類型
   * 如果當前操作不在新協議的操作列表中，則切換到預設操作
   */
  useEffect(() => {
    const defaultOp = getDefaultOperation()
    if (!defaultOp) return

    // 獲取當前協議支援的操作列表
    let availableOps: { value: string; label: string }[] = []
    if (isModbus) {
      availableOps = [
        { value: 'read_coils', label: '讀取線圈 (Coils)' },
        { value: 'read_discrete_inputs', label: '讀取離散輸入 (Discrete Inputs)' },
        { value: 'read_holding_registers', label: '讀取保持暫存器 (Holding Registers)' },
        { value: 'read_input_registers', label: '讀取輸入暫存器 (Input Registers)' },
      ]
    } else if (isFatek) {
      availableOps = [
        { value: 'read_registers', label: '讀取暫存器' },
        { value: 'read_status', label: '讀取狀態' },
      ]
    } else if (isMCProtocol) {
      availableOps = [
        { value: 'batch_read_word', label: '批量讀取字組' },
        { value: 'batch_read_bit', label: '批量讀取位元' },
      ]
    }

    const isValidOp = availableOps.some(op => op.value === operation)
    
    // 如果沒有操作類型，或當前操作無效，則設置為預設操作
    if (!operation || !isValidOp) {
      setOperation(defaultOp)
    }
  }, [protocol, getDefaultOperation, operation, isModbus, isFatek, isMCProtocol])

  return (
    <div className="space-y-6">
      {/* 配置區域 */}
      <div className="space-y-4">
        <h3 className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider border-b dark:border-gray-700 pb-2">
          掃描配置
        </h3>

        {/* 操作類型選擇 */}
        <div className="space-y-1">
          <label className="text-sm font-medium text-gray-700 dark:text-gray-300">操作類型</label>
          <div className="relative">
            <select
              value={operation}
              onChange={(e) => setOperation(e.target.value)}
              disabled={scanning}
              className="w-full px-4 py-2 appearance-none border border-gray-300 dark:border-gray-600 rounded-lg shadow-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-white hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400 focus:border-blue-500 dark:focus:border-blue-400 transition-all disabled:bg-gray-100 dark:disabled:bg-gray-700 disabled:cursor-not-allowed"
            >
              <option value="">請選擇操作類型</option>
              {getAvailableOperations().map((op) => (
                <option key={op.value} value={op.value}>
                  {op.label}
                </option>
              ))}
            </select>
            <div className="absolute inset-y-0 right-0 flex items-center px-2 pointer-events-none text-gray-500 dark:text-gray-400">
              <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
              </svg>
            </div>
          </div>
        </div>

        {/* 協議特定參數 */}
        <div className="grid grid-cols-2 gap-4">
          {isFatek && (
            <div className="space-y-1">
              <label className="text-sm font-medium text-gray-700 dark:text-gray-300">組件符號 (Symbol)</label>
              <input
                type="text"
                value={symbol}
                onChange={(e) => setSymbol(e.target.value.toUpperCase())}
                disabled={scanning}
                placeholder="D, M, X, Y"
                className="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white rounded-lg shadow-sm hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400 focus:border-blue-500 dark:focus:border-blue-400 transition-all disabled:bg-gray-100 dark:disabled:bg-gray-700 disabled:cursor-not-allowed"
              />
            </div>
          )}

          {isMCProtocol && (
            <div className="space-y-1">
              <label className="text-sm font-medium text-gray-700 dark:text-gray-300">設備代號 (Device)</label>
              <input
                type="text"
                value={device}
                onChange={(e) => setDevice(e.target.value.toUpperCase())}
                disabled={scanning}
                placeholder="D, M, X, Y"
                className="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white rounded-lg shadow-sm hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400 focus:border-blue-500 dark:focus:border-blue-400 transition-all disabled:bg-gray-100 dark:disabled:bg-gray-700 disabled:cursor-not-allowed"
              />
            </div>
          )}
        </div>

        {/* 掃描範圍配置 */}
        <div className="grid grid-cols-3 gap-4">
          <div className="space-y-1">
            <label className="text-sm font-medium text-gray-700 dark:text-gray-300">起始位置</label>
            <input
              type="number"
              value={startAddress}
              onChange={(e) => setStartAddress(parseInt(e.target.value) || 0)}
              disabled={scanning}
              min={0}
              className="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white rounded-lg shadow-sm hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400 focus:border-blue-500 dark:focus:border-blue-400 transition-all disabled:bg-gray-100 dark:disabled:bg-gray-700 disabled:cursor-not-allowed"
            />
          </div>

          <div className="space-y-1">
            <label className="text-sm font-medium text-gray-700 dark:text-gray-300">最大位置</label>
            <input
              type="number"
              value={maxAddress}
              onChange={(e) => setMaxAddress(parseInt(e.target.value) || 100)}
              disabled={scanning}
              min={startAddress}
              className="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white rounded-lg shadow-sm hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400 focus:border-blue-500 dark:focus:border-blue-400 transition-all disabled:bg-gray-100 dark:disabled:bg-gray-700 disabled:cursor-not-allowed"
            />
          </div>

          <div className="space-y-1">
            <label className="text-sm font-medium text-gray-700 dark:text-gray-300">掃描步進</label>
            <input
              type="number"
              value={scanStep}
              onChange={(e) => setScanStep(Math.max(1, parseInt(e.target.value) || 1))}
              disabled={scanning}
              min={1}
              className="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white rounded-lg shadow-sm hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400 focus:border-blue-500 dark:focus:border-blue-400 transition-all disabled:bg-gray-100 dark:disabled:bg-gray-700 disabled:cursor-not-allowed"
            />
          </div>
        </div>
      </div>

      {/* 操作按鈕 */}
      <div className="flex gap-3">
        {!scanning ? (
          <button
            onClick={handleScan}
            disabled={!connectionId || !operation}
            className="flex-1 bg-blue-600 text-white px-4 py-3 rounded-xl font-semibold shadow-lg shadow-blue-500/30 hover:bg-blue-700 hover:shadow-blue-500/40 active:scale-[0.98] transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
            開始掃描
          </button>
        ) : (
          <button
            onClick={handleStopScan}
            className="flex-1 bg-red-600 text-white px-4 py-3 rounded-xl font-semibold shadow-lg shadow-red-500/30 hover:bg-red-700 hover:shadow-red-500/40 active:scale-[0.98] transition-all flex items-center justify-center gap-2"
          >
            <span className="animate-pulse w-2 h-2 rounded-full bg-white"></span>
            停止掃描
          </button>
        )}
        <button
          onClick={handleClearResults}
          disabled={scanResults.length === 0}
          className="px-4 py-3 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded-xl font-medium hover:bg-gray-50 dark:hover:bg-gray-700 active:scale-[0.98] transition-all disabled:opacity-50 disabled:cursor-not-allowed"
        >
          清除結果
        </button>
      </div>

      {/* 掃描統計 */}
      {scanResults.length > 0 && (
        <div className="bg-gradient-to-r from-blue-50 to-purple-50 dark:from-blue-900/20 dark:to-purple-900/20 border border-blue-200 dark:border-blue-800 rounded-xl p-4 animate-fade-in">
          <div className="flex items-center justify-between mb-2">
            <h4 className="text-sm font-semibold text-gray-800 dark:text-gray-200">掃描統計</h4>
            <span className="text-xs text-gray-500 dark:text-gray-400">
              進度: {scanResults.length} / {Math.floor((maxAddress - startAddress) / scanStep) + 1}
            </span>
          </div>
          <div className="grid grid-cols-3 gap-4">
            <div className="bg-white dark:bg-gray-800 rounded-lg p-3 shadow-sm">
              <div className="text-xs text-gray-500 dark:text-gray-400 mb-1">成功設備</div>
              <div className="text-2xl font-bold text-green-600 dark:text-green-400">{deviceCount}</div>
            </div>
            <div className="bg-white dark:bg-gray-800 rounded-lg p-3 shadow-sm">
              <div className="text-xs text-gray-500 dark:text-gray-400 mb-1">失敗地址</div>
              <div className="text-2xl font-bold text-red-600 dark:text-red-400">
                {scanResults.length - deviceCount}
              </div>
            </div>
            <div className="bg-white dark:bg-gray-800 rounded-lg p-3 shadow-sm">
              <div className="text-xs text-gray-500 dark:text-gray-400 mb-1">總掃描數</div>
              <div className="text-2xl font-bold text-gray-800 dark:text-white">{scanResults.length}</div>
            </div>
          </div>
        </div>
      )}

      {/* 掃描結果列表 */}
      {scanResults.length > 0 && (
        <div className="bg-gray-50 dark:bg-gray-900/50 p-4 rounded-xl border border-gray-200 dark:border-gray-700 animate-fade-in">
          <div className="flex justify-between items-center mb-3">
            <h4 className="text-xs font-bold text-gray-500 dark:text-gray-400 uppercase">掃描結果</h4>
            <span className="text-xs text-gray-400 dark:text-gray-500">
              顯示前 50 筆結果
            </span>
          </div>
          <div className="max-h-64 overflow-y-auto scrollbar-thin space-y-2">
            {scanResults.slice(0, 50).map((result, idx) => (
              <div
                key={idx}
                className={`bg-white dark:bg-gray-800 p-3 rounded-lg border shadow-sm flex items-center justify-between ${
                  result.success
                    ? 'border-green-200 dark:border-green-800 bg-green-50 dark:bg-green-900/20'
                    : 'border-red-200 dark:border-red-800 bg-red-50 dark:bg-red-900/20'
                }`}
              >
                <div className="flex items-center gap-3">
                  <div
                    className={`w-2 h-2 rounded-full ${
                      result.success ? 'bg-green-500' : 'bg-red-500'
                    }`}
                  ></div>
                  <div>
                    <div className="text-sm font-mono font-semibold text-gray-800 dark:text-gray-200">
                      位置: {result.address}
                    </div>
                    {result.error && (
                      <div className="text-xs text-red-600 dark:text-red-400 mt-1">{result.error}</div>
                    )}
                    {result.success && result.data && (
                      <div className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                        數據: {JSON.stringify(result.data.values || result.data).slice(0, 50)}
                      </div>
                    )}
                  </div>
                </div>
                {result.success && (
                  <span className="text-xs font-medium text-green-600 dark:text-green-400 bg-green-100 dark:bg-green-900/30 px-2 py-1 rounded">
                    成功
                  </span>
                )}
              </div>
            ))}
          </div>
          {scanResults.length > 50 && (
            <div className="mt-3 text-center text-xs text-gray-500 dark:text-gray-400">
              還有 {scanResults.length - 50} 筆結果未顯示
            </div>
          )}
        </div>
      )}

      {/* 提示信息 */}
      {!connectionId && (
        <div className="bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-lg p-3 text-sm text-amber-800 dark:text-amber-400">
          <span className="font-medium">提示：</span>請先建立連線後再進行掃描。
        </div>
      )}
    </div>
  )
}