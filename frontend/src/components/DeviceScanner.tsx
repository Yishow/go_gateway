import { useState, useCallback, useRef, useEffect, useMemo } from 'react'
import { useTestAPI } from '../services/api'

interface DeviceScannerProps {
  protocol: string
  connectionId: string | null
  baseConfig?: Record<string, any> // 基礎連線配置，用於站號掃描時創建臨時連線
}

/**
 * 掃描結果介面
 */
interface ScanResult {
  stationOrIp: number | string // 站號或 IP 地址
  success: boolean
  data?: any
  error?: string
  responseTime?: number // 響應時間（毫秒）
  testAddress?: number // 測試位置（站號掃描時）
}

/**
 * 掃描統計介面
 */
interface ScanStatistics {
  totalScanned: number
  foundDevices: number
  failedScans: number
  averageResponseTime: number
  scanSpeed: number // 站號/秒 或 IP/秒
  startTime?: number
}

/**
 * FATEK Symbol 選項
 */
const FATEK_SYMBOLS = [
  { value: 'X', label: 'X (輸入)' },
  { value: 'Y', label: 'Y (輸出)' },
  { value: 'M', label: 'M (內部繼電器)' },
  { value: 'S', label: 'S (步進繼電器)' },
  { value: 'T', label: 'T (計時器狀態)' },
  { value: 'C', label: 'C (計數器狀態)' },
  { value: 'R', label: 'R (暫存器)' },
  { value: 'D', label: 'D (資料暫存器)' },
  { value: 'RT', label: 'RT (計時器值)' },
  { value: 'RC', label: 'RC (計數器值)' },
  { value: 'DR', label: 'DR (32位資料暫存器)' },
  { value: 'DW', label: 'DW (32位離散)' },
  { value: 'DWM', label: 'DWM' },
]

/**
 * 設備掃描器組件
 * 支援掃描設備站號（Modbus、FATEK）或 IP 網段（MC Protocol）
 * 適配所有協議：Modbus、FATEK、MC Protocol
 */
export default function DeviceScanner({
  protocol,
  connectionId,
  baseConfig = {},
}: DeviceScannerProps) {
  // 站號掃描配置（Modbus、FATEK）
  const [startStation, setStartStation] = useState<number>(1)
  const [maxStation, setMaxStation] = useState<number>(247)
  const [testAddress, setTestAddress] = useState<number>(0)
  
  // MC Protocol IP 網段掃描配置
  const [ipNetwork, setIpNetwork] = useState<string>('192.168.1.0/24')
  const [mcPort, setMcPort] = useState<number>(5000)
  const [mcTimeout, setMcTimeout] = useState<number>(2000)
  
  // 協議特定參數
  const [symbol, setSymbol] = useState<string>('D') // FATEK
  const [operation, setOperation] = useState<string>('')
  
  // 掃描狀態
  const [scanning, setScanning] = useState<boolean>(false)
  const [scanResults, setScanResults] = useState<ScanResult[]>([])
  const [scanStatistics, setScanStatistics] = useState<ScanStatistics>({
    totalScanned: 0,
    foundDevices: 0,
    failedScans: 0,
    averageResponseTime: 0,
    scanSpeed: 0,
  })
  const [currentScanning, setCurrentScanning] = useState<number | string>('')
  const [scanProgress, setScanProgress] = useState<number>(0)
  
  // 掃描速度控制
  const [scanDelay, setScanDelay] = useState<number>(50)
  const [concurrentScan, setConcurrentScan] = useState<boolean>(false)
  const [concurrency, setConcurrency] = useState<number>(3)
  
  // 結果過濾和排序
  const [filterType, setFilterType] = useState<'all' | 'success' | 'failed'>('all')
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc')
  const [searchTerm, setSearchTerm] = useState<string>('')
  
  // 驗證錯誤
  const [validationErrors, setValidationErrors] = useState<Record<string, string>>({})
  
  const { read, connect: connectAPI, disconnect } = useTestAPI()
  
  // 包裝 connect 函數以處理錯誤
  const connect = useCallback(async (protocol: string, config: Record<string, any>) => {
    return await connectAPI(protocol, config)
  }, [connectAPI])
  const scanAbortRef = useRef<boolean>(false)
  const scanStartTimeRef = useRef<number>(0)

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
   * 獲取可用 Symbol 列表（根據操作類型過濾）
   */
  const getAvailableSymbols = () => {
    if (!isFatek) return []
    
    // read_status 只能用離散型 Symbol
    if (operation === 'read_status') {
      return FATEK_SYMBOLS.filter(s => ['X', 'Y', 'M', 'S', 'T', 'C'].includes(s.value))
    }
    
    // read_registers 可以用所有 Symbol
    return FATEK_SYMBOLS
  }

  /**
   * 驗證輸入
   */
  const validateInputs = useCallback((): boolean => {
    const errors: Record<string, string> = {}
    
    if (isMCProtocol) {
      // MC Protocol: 驗證 IP 網段
      const cidrRegex = /^(\d{1,3}\.){3}\d{1,3}\/\d{1,2}$/
      if (!cidrRegex.test(ipNetwork)) {
        errors.ipNetwork = 'IP 網段格式錯誤，請使用 CIDR 格式（如 192.168.1.0/24）'
      } else {
        const [, prefix] = ipNetwork.split('/')
        const prefixLen = parseInt(prefix)
        if (prefixLen < 16 || prefixLen > 30) {
          errors.ipNetwork = '網段範圍過大，建議使用 /16 到 /30'
        }
      }
      
      if (mcPort < 1 || mcPort > 65535) {
        errors.mcPort = '埠號必須在 1-65535 之間'
      }
      
      if (mcTimeout < 100 || mcTimeout > 10000) {
        errors.mcTimeout = '超時時間必須在 100-10000 毫秒之間'
      }
    } else {
      // Modbus 或 FATEK: 驗證站號範圍
      const minStation = isModbus ? 1 : 0
      const maxStationLimit = isModbus ? 247 : 254
      
      if (startStation < minStation || startStation > maxStationLimit) {
        errors.startStation = `起始站號必須在 ${minStation}-${maxStationLimit} 之間`
      }
      
      if (maxStation < minStation || maxStation > maxStationLimit) {
        errors.maxStation = `最大站號必須在 ${minStation}-${maxStationLimit} 之間`
      }
      
      if (startStation > maxStation) {
        errors.maxStation = '最大站號必須 >= 起始站號'
      }
      
      if (testAddress < 0) {
        errors.testAddress = '測試位置必須 >= 0'
      }
      
      if (!operation) {
        errors.operation = '請選擇操作類型'
      }
      
      if (isFatek && !symbol) {
        errors.symbol = '請選擇組件符號'
      }
    }
    
    setValidationErrors(errors)
    return Object.keys(errors).length === 0
  }, [isMCProtocol, isModbus, isFatek, ipNetwork, mcPort, mcTimeout, startStation, maxStation, testAddress, operation, symbol])

  /**
   * 解析 CIDR 格式，生成 IP 列表
   */
  const parseCIDR = useCallback((cidr: string): string[] => {
    const [ip, prefix] = cidr.split('/')
    const prefixLen = parseInt(prefix)
    
    if (prefixLen < 0 || prefixLen > 32) {
      return []
    }
    
    const ipParts = ip.split('.').map(Number)
    if (ipParts.length !== 4 || ipParts.some(p => p < 0 || p > 255)) {
      return []
    }
    
    const ipNum = (ipParts[0] << 24) + (ipParts[1] << 16) + (ipParts[2] << 8) + ipParts[3]
    const mask = (0xFFFFFFFF << (32 - prefixLen)) >>> 0
    const network = ipNum & mask
    const hostCount = Math.pow(2, 32 - prefixLen) - 2 // 減去網路地址和廣播地址
    
    const ips: string[] = []
    for (let i = 1; i <= hostCount && i < 256; i++) { // 限制最多 255 個 IP
      const hostIp = network + i
      const ipStr = [
        (hostIp >>> 24) & 0xFF,
        (hostIp >>> 16) & 0xFF,
        (hostIp >>> 8) & 0xFF,
        hostIp & 0xFF,
      ].join('.')
      ips.push(ipStr)
    }
    
    return ips
  }, [])

  /**
   * 掃描單個 Modbus 站號
   */
  const scanModbusStation = useCallback(async (station: number, testAddr: number): Promise<ScanResult> => {
    const startTime = Date.now()
    let tempConnectionId: string | null = null
    
    try {
      // 為每個站號創建臨時連線
      const scanConfig = {
        ...baseConfig,
        unitID: station, // 設置站號
      }
      
      // 創建臨時連線
      const connectResult = await connect(protocol, scanConfig)
      tempConnectionId = connectResult.connection_id
      
      // 執行讀取
      const readParams: any = {
        operation,
        address: testAddr,
        count: 1,
      }
      
      const result = await read(tempConnectionId, readParams)
      const responseTime = Date.now() - startTime
      
      // 斷開臨時連線
      if (tempConnectionId) {
        try {
          await disconnect(tempConnectionId)
        } catch (_e) {
          // 忽略斷線錯誤
        }
      }
      
      return {
        stationOrIp: station,
        success: true,
        data: result,
        responseTime,
        testAddress: testAddr,
      }
    } catch (error: any) {
      const responseTime = Date.now() - startTime
      
      // 確保斷開連線
      if (tempConnectionId) {
        try {
          await disconnect(tempConnectionId)
        } catch (_e) {
          // 忽略斷線錯誤
        }
      }
      
      return {
        stationOrIp: station,
        success: false,
        error: error.message || '讀取失敗',
        responseTime,
        testAddress: testAddr,
      }
    }
  }, [protocol, baseConfig, connect, disconnect, read, operation])

  /**
   * 掃描單個 FATEK 站號
   */
  const scanFatekStation = useCallback(async (station: number, testAddr: number): Promise<ScanResult> => {
    const startTime = Date.now()
    let tempConnectionId: string | null = null
    
    try {
      // 為每個站號創建臨時連線
      const scanConfig = {
        ...baseConfig,
        station: station, // 設置站號
      }
      
      // 創建臨時連線
      const connectResult = await connect(protocol, scanConfig)
      tempConnectionId = connectResult.connection_id
      
      // 執行讀取
      const readParams: any = {
        operation,
        address: testAddr,
        count: 1,
        symbol,
      }
      
      const result = await read(tempConnectionId, readParams)
      const responseTime = Date.now() - startTime
      
      // 斷開臨時連線
      if (tempConnectionId) {
        try {
          await disconnect(tempConnectionId)
        } catch (_e) {
          // 忽略斷線錯誤
        }
      }
      
      return {
        stationOrIp: station,
        success: true,
        data: result,
        responseTime,
        testAddress: testAddr,
      }
    } catch (error: any) {
      const responseTime = Date.now() - startTime
      
      // 確保斷開連線
      if (tempConnectionId) {
        try {
          await disconnect(tempConnectionId)
        } catch (_e) {
          // 忽略斷線錯誤
        }
      }
      
      return {
        stationOrIp: station,
        success: false,
        error: error.message || '讀取失敗',
        responseTime,
        testAddress: testAddr,
      }
    }
  }, [protocol, baseConfig, connect, disconnect, read, operation, symbol])

  /**
   * 掃描單個 MC Protocol IP
   */
  const scanMCProtocolIP = useCallback(async (ip: string): Promise<ScanResult> => {
    const startTime = Date.now()
    let tempConnectionId: string | null = null
    
    try {
      // 為每個 IP 創建臨時連線
      const config = {
        host: ip,
        port: mcPort,
        timeout: mcTimeout / 1000, // 轉換為秒
      }
      
      const connectResult = await connect('mcprotocol_tcp', config)
      tempConnectionId = connectResult.connection_id
      
      // 嘗試輕量級讀取測試
      const readParams: any = {
        operation: 'batch_read_word',
        address: 0,
        count: 1,
        device: 'D',
      }
      
      const result = await read(tempConnectionId, readParams)
      const responseTime = Date.now() - startTime
      
      // 斷開臨時連線
      if (tempConnectionId) {
        try {
          await disconnect(tempConnectionId)
        } catch (_e) {
          // 忽略斷線錯誤
        }
      }
      
      return {
        stationOrIp: ip,
        success: true,
        data: result,
        responseTime,
      }
    } catch (error: any) {
      const responseTime = Date.now() - startTime
      
      // 確保斷開連線
      if (tempConnectionId) {
        try {
          await disconnect(tempConnectionId)
        } catch (_e) {
          // 忽略斷線錯誤
        }
      }
      
      return {
        stationOrIp: ip,
        success: false,
        error: error.message || '連線失敗',
        responseTime,
      }
    }
  }, [connect, disconnect, read, mcPort, mcTimeout])

  /**
   * 並發掃描實現
   */
  const scanWithConcurrency = useCallback(async (
    items: (number | string)[],
    scanFn: (item: number | string) => Promise<ScanResult>
  ): Promise<ScanResult[]> => {
    const results: ScanResult[] = []
    const batches: (number | string)[][] = []
    
    // 將項目分成批次
    for (let i = 0; i < items.length; i += concurrency) {
      batches.push(items.slice(i, i + concurrency))
    }
    
    // 依序處理每個批次
    for (const batch of batches) {
      if (scanAbortRef.current) break
      
      const batchResults = await Promise.all(
        batch.map(item => scanFn(item))
      )
      
      results.push(...batchResults)
      
      // 更新進度
      const progress = (results.length / items.length) * 100
      setScanProgress(progress)
      setScanResults([...results])
      
      // 批次間延遲
      if (scanDelay > 0) {
        await new Promise(resolve => setTimeout(resolve, scanDelay))
      }
    }
    
    return results
  }, [concurrency, scanDelay])

  /**
   * 執行掃描
   */
  const handleScan = async () => {
    if (!validateInputs()) {
      return
    }
    
    // 對於 Modbus 和 FATEK，需要基礎配置來創建臨時連線
    if (!isMCProtocol) {
      // 檢查是否有基礎配置（從現有連線或 baseConfig）
      const hasBaseConfig = Object.keys(baseConfig).length > 0 || connectionId !== null
      if (!hasBaseConfig) {
        setValidationErrors({ connection: '請先建立連線以獲取基礎配置，或提供基礎連線配置' })
        return
      }
    }

    setScanning(true)
    setScanResults([])
    setScanProgress(0)
    setCurrentScanning('')
    scanAbortRef.current = false
    scanStartTimeRef.current = Date.now()
    
    const results: ScanResult[] = []
    const responseTimes: number[] = []

    try {
      if (isMCProtocol) {
        // MC Protocol: IP 網段掃描
        const ips = parseCIDR(ipNetwork)
        if (ips.length === 0) {
          setValidationErrors({ ipNetwork: '無法解析 IP 網段' })
          setScanning(false)
          return
        }
        
        const totalIPs = ips.length
        setScanStatistics(prev => ({ ...prev, totalScanned: totalIPs }))
        
        if (concurrentScan) {
          // 並發掃描
          const scanResults = await scanWithConcurrency(ips, (ip) => scanMCProtocolIP(ip as string))
          results.push(...scanResults)
        } else {
          // 順序掃描
          for (let i = 0; i < ips.length && !scanAbortRef.current; i++) {
            const ip = ips[i]
            setCurrentScanning(ip)
            
            const result = await scanMCProtocolIP(ip)
            results.push(result)
            
            if (result.responseTime) {
              responseTimes.push(result.responseTime)
            }
            
            // 更新進度和結果
            const progress = ((i + 1) / totalIPs) * 100
            setScanProgress(progress)
            setScanResults([...results])
            
            // 延遲
            if (scanDelay > 0 && i < ips.length - 1) {
              await new Promise(resolve => setTimeout(resolve, scanDelay))
            }
          }
        }
      } else {
        // Modbus 或 FATEK: 站號掃描
        const stations: number[] = []
        for (let i = startStation; i <= maxStation; i++) {
          stations.push(i)
        }
        
        const totalStations = stations.length
        setScanStatistics(prev => ({ ...prev, totalScanned: totalStations }))
        
        if (concurrentScan) {
          // 並發掃描
          const scanFn = isFatek 
            ? (station: number | string) => scanFatekStation(station as number, testAddress)
            : (station: number | string) => scanModbusStation(station as number, testAddress)
          
          const scanResults = await scanWithConcurrency(stations, scanFn)
          results.push(...scanResults)
        } else {
          // 順序掃描
          for (let i = 0; i < stations.length && !scanAbortRef.current; i++) {
            const station = stations[i]
            setCurrentScanning(station)
            
            const result = isFatek
              ? await scanFatekStation(station, testAddress)
              : await scanModbusStation(station, testAddress)
            
            results.push(result)
            
            if (result.responseTime) {
              responseTimes.push(result.responseTime)
            }
            
            // 更新進度和結果
            const progress = ((i + 1) / totalStations) * 100
            setScanProgress(progress)
            setScanResults([...results])
            
            // 延遲
            if (scanDelay > 0 && i < stations.length - 1) {
              await new Promise(resolve => setTimeout(resolve, scanDelay))
            }
          }
        }
      }
      
      // 計算統計信息
      const totalScanned = results.length
      const foundDevices = results.filter(r => r.success).length
      const failedScans = totalScanned - foundDevices
      const averageResponseTime = responseTimes.length > 0
        ? responseTimes.reduce((a, b) => a + b, 0) / responseTimes.length
        : 0
      const elapsedTime = (Date.now() - scanStartTimeRef.current) / 1000 // 秒
      const scanSpeed = elapsedTime > 0 ? totalScanned / elapsedTime : 0
      
      setScanStatistics({
        totalScanned,
        foundDevices,
        failedScans,
        averageResponseTime,
        scanSpeed,
        startTime: scanStartTimeRef.current,
      })
      
    } catch (error: any) {
      console.error('掃描錯誤:', error)
      setValidationErrors({ scan: error.message || '掃描過程中發生錯誤' })
    } finally {
      setScanning(false)
      setCurrentScanning('')
      scanAbortRef.current = false
    }
  }

  /**
   * 停止掃描
   */
  const handleStopScan = () => {
    scanAbortRef.current = true
    setScanning(false)
    setCurrentScanning('')
  }

  /**
   * 清除結果
   */
  const handleClearResults = () => {
    setScanResults([])
    setScanStatistics({
      totalScanned: 0,
      foundDevices: 0,
      failedScans: 0,
      averageResponseTime: 0,
      scanSpeed: 0,
    })
    setScanProgress(0)
  }

  /**
   * 過濾結果
   */
  const filteredResults = useMemo(() => {
    let filtered = scanResults
    
    // 按成功/失敗過濾
    if (filterType === 'success') {
      filtered = filtered.filter(r => r.success)
    } else if (filterType === 'failed') {
      filtered = filtered.filter(r => !r.success)
    }
    
    // 搜索過濾
    if (searchTerm) {
      const term = searchTerm.toLowerCase()
      filtered = filtered.filter(r => {
        const stationOrIp = String(r.stationOrIp).toLowerCase()
        return stationOrIp.includes(term)
      })
    }
    
    // 排序
    filtered.sort((a, b) => {
      const aVal = typeof a.stationOrIp === 'number' ? a.stationOrIp : parseInt(a.stationOrIp.split('.').join(''))
      const bVal = typeof b.stationOrIp === 'number' ? b.stationOrIp : parseInt(b.stationOrIp.split('.').join(''))
      
      return sortOrder === 'asc' ? aVal - bVal : bVal - aVal
    })
    
    return filtered
  }, [scanResults, filterType, searchTerm, sortOrder])

  /**
   * 當協議變更時，自動更新操作類型和站號範圍
   */
  useEffect(() => {
    const defaultOp = getDefaultOperation()
    if (defaultOp && !operation) {
      setOperation(defaultOp)
    }
    
    // 設置預設站號範圍
    if (isModbus) {
      setStartStation(1)
      setMaxStation(247)
    } else if (isFatek) {
      setStartStation(0)
      setMaxStation(254)
    }
  }, [protocol, getDefaultOperation, operation, isModbus, isFatek])

  // 計算總掃描數（用於進度顯示）
  const totalToScan = useMemo(() => {
    if (isMCProtocol) {
      const ips = parseCIDR(ipNetwork)
      return ips.length
    } else {
      return maxStation - startStation + 1
    }
  }, [isMCProtocol, ipNetwork, startStation, maxStation, parseCIDR])

  return (
    <div className="space-y-6">
      {/* 配置區域 */}
      <div className="space-y-4">
        <h3 className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider border-b dark:border-gray-700 pb-2">
          掃描配置
        </h3>

        {/* 操作類型選擇 */}
        {!isMCProtocol && (
          <div className="space-y-1">
            <label className="text-sm font-medium text-gray-700 dark:text-gray-300">操作類型</label>
            <div className="relative">
              <select
                value={operation}
                onChange={(e) => setOperation(e.target.value)}
                disabled={scanning}
                className={`w-full px-4 py-2 appearance-none border rounded-lg shadow-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-white hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400 focus:border-blue-500 dark:focus:border-blue-400 transition-all disabled:bg-gray-100 dark:disabled:bg-gray-700 disabled:cursor-not-allowed ${
                  validationErrors.operation ? 'border-red-500 dark:border-red-500' : 'border-gray-300 dark:border-gray-600'
                }`}
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
            {validationErrors.operation && (
              <p className="text-xs text-red-600 dark:text-red-400">{validationErrors.operation}</p>
            )}
          </div>
        )}

        {/* MC Protocol: IP 網段配置 */}
        {isMCProtocol && (
          <div className="space-y-4">
            <div className="space-y-1">
              <label className="text-sm font-medium text-gray-700 dark:text-gray-300">IP 網段 (CIDR)</label>
              <input
                type="text"
                value={ipNetwork}
                onChange={(e) => setIpNetwork(e.target.value)}
                disabled={scanning}
                placeholder="192.168.1.0/24"
                className={`w-full px-3 py-2 bg-white dark:bg-gray-800 border rounded-lg shadow-sm text-gray-900 dark:text-white hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400 focus:border-blue-500 dark:focus:border-blue-400 transition-all disabled:bg-gray-100 dark:disabled:bg-gray-700 disabled:cursor-not-allowed ${
                  validationErrors.ipNetwork ? 'border-red-500 dark:border-red-500' : 'border-gray-300 dark:border-gray-600'
                }`}
              />
              {validationErrors.ipNetwork && (
                <p className="text-xs text-red-600 dark:text-red-400">{validationErrors.ipNetwork}</p>
              )}
            </div>
            
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-1">
                <label className="text-sm font-medium text-gray-700 dark:text-gray-300">埠號</label>
                <input
                  type="number"
                  value={mcPort}
                  onChange={(e) => setMcPort(parseInt(e.target.value) || 5000)}
                  disabled={scanning}
                  min={1}
                  max={65535}
                  className={`w-full px-3 py-2 bg-white dark:bg-gray-800 border rounded-lg shadow-sm text-gray-900 dark:text-white hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400 focus:border-blue-500 dark:focus:border-blue-400 transition-all disabled:bg-gray-100 dark:disabled:bg-gray-700 disabled:cursor-not-allowed ${
                    validationErrors.mcPort ? 'border-red-500 dark:border-red-500' : 'border-gray-300 dark:border-gray-600'
                  }`}
                />
                {validationErrors.mcPort && (
                  <p className="text-xs text-red-600 dark:text-red-400">{validationErrors.mcPort}</p>
                )}
              </div>
              
              <div className="space-y-1">
                <label className="text-sm font-medium text-gray-700 dark:text-gray-300">超時時間 (ms)</label>
                <input
                  type="number"
                  value={mcTimeout}
                  onChange={(e) => setMcTimeout(parseInt(e.target.value) || 2000)}
                  disabled={scanning}
                  min={100}
                  max={10000}
                  className={`w-full px-3 py-2 bg-white dark:bg-gray-800 border rounded-lg shadow-sm text-gray-900 dark:text-white hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400 focus:border-blue-500 dark:focus:border-blue-400 transition-all disabled:bg-gray-100 dark:disabled:bg-gray-700 disabled:cursor-not-allowed ${
                    validationErrors.mcTimeout ? 'border-red-500 dark:border-red-500' : 'border-gray-300 dark:border-gray-600'
                  }`}
                />
                {validationErrors.mcTimeout && (
                  <p className="text-xs text-red-600 dark:text-red-400">{validationErrors.mcTimeout}</p>
                )}
              </div>
            </div>
          </div>
        )}

        {/* Modbus/FATEK: 站號掃描配置 */}
        {!isMCProtocol && (
          <div className="space-y-4">
            <div className="grid grid-cols-3 gap-4">
              <div className="space-y-1">
                <label className="text-sm font-medium text-gray-700 dark:text-gray-300">起始站號</label>
                <input
                  type="number"
                  value={startStation}
                  onChange={(e) => setStartStation(parseInt(e.target.value) || 0)}
                  disabled={scanning}
                  min={isModbus ? 1 : 0}
                  max={isModbus ? 247 : 254}
                  className={`w-full px-3 py-2 bg-white dark:bg-gray-800 border rounded-lg shadow-sm text-gray-900 dark:text-white hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400 focus:border-blue-500 dark:focus:border-blue-400 transition-all disabled:bg-gray-100 dark:disabled:bg-gray-700 disabled:cursor-not-allowed ${
                    validationErrors.startStation ? 'border-red-500 dark:border-red-500' : 'border-gray-300 dark:border-gray-600'
                  }`}
                />
                {validationErrors.startStation && (
                  <p className="text-xs text-red-600 dark:text-red-400">{validationErrors.startStation}</p>
                )}
              </div>

              <div className="space-y-1">
                <label className="text-sm font-medium text-gray-700 dark:text-gray-300">最大站號</label>
                <input
                  type="number"
                  value={maxStation}
                  onChange={(e) => setMaxStation(parseInt(e.target.value) || (isModbus ? 247 : 254))}
                  disabled={scanning}
                  min={isModbus ? 1 : 0}
                  max={isModbus ? 247 : 254}
                  className={`w-full px-3 py-2 bg-white dark:bg-gray-800 border rounded-lg shadow-sm text-gray-900 dark:text-white hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400 focus:border-blue-500 dark:focus:border-blue-400 transition-all disabled:bg-gray-100 dark:disabled:bg-gray-700 disabled:cursor-not-allowed ${
                    validationErrors.maxStation ? 'border-red-500 dark:border-red-500' : 'border-gray-300 dark:border-gray-600'
                  }`}
                />
                {validationErrors.maxStation && (
                  <p className="text-xs text-red-600 dark:text-red-400">{validationErrors.maxStation}</p>
                )}
              </div>

              <div className="space-y-1">
                <label className="text-sm font-medium text-gray-700 dark:text-gray-300">測試位置</label>
                <input
                  type="number"
                  value={testAddress}
                  onChange={(e) => setTestAddress(parseInt(e.target.value) || 0)}
                  disabled={scanning}
                  min={0}
                  className={`w-full px-3 py-2 bg-white dark:bg-gray-800 border rounded-lg shadow-sm text-gray-900 dark:text-white hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400 focus:border-blue-500 dark:focus:border-blue-400 transition-all disabled:bg-gray-100 dark:disabled:bg-gray-700 disabled:cursor-not-allowed ${
                    validationErrors.testAddress ? 'border-red-500 dark:border-red-500' : 'border-gray-300 dark:border-gray-600'
                  }`}
                />
                {validationErrors.testAddress && (
                  <p className="text-xs text-red-600 dark:text-red-400">{validationErrors.testAddress}</p>
                )}
              </div>
            </div>

            {/* FATEK Symbol 選擇 */}
            {isFatek && (
              <div className="space-y-1">
                <label className="text-sm font-medium text-gray-700 dark:text-gray-300">組件符號 (Symbol)</label>
                <div className="relative">
                  <select
                    value={symbol}
                    onChange={(e) => setSymbol(e.target.value)}
                    disabled={scanning}
                    className={`w-full px-4 py-2 appearance-none border rounded-lg shadow-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-white hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400 focus:border-blue-500 dark:focus:border-blue-400 transition-all disabled:bg-gray-100 dark:disabled:bg-gray-700 disabled:cursor-not-allowed ${
                      validationErrors.symbol ? 'border-red-500 dark:border-red-500' : 'border-gray-300 dark:border-gray-600'
                    }`}
                  >
                    {getAvailableSymbols().map((sym) => (
                      <option key={sym.value} value={sym.value}>
                        {sym.label}
                      </option>
                    ))}
                  </select>
                  <div className="absolute inset-y-0 right-0 flex items-center px-2 pointer-events-none text-gray-500 dark:text-gray-400">
                    <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                    </svg>
                  </div>
                </div>
                {validationErrors.symbol && (
                  <p className="text-xs text-red-600 dark:text-red-400">{validationErrors.symbol}</p>
                )}
              </div>
            )}
          </div>
        )}

        {/* 掃描速度控制 */}
        <div className="grid grid-cols-3 gap-4 pt-2 border-t border-gray-200 dark:border-gray-700">
          <div className="space-y-1">
            <label className="text-sm font-medium text-gray-700 dark:text-gray-300">掃描延遲 (ms)</label>
            <input
              type="number"
              value={scanDelay}
              onChange={(e) => setScanDelay(Math.max(0, parseInt(e.target.value) || 0))}
              disabled={scanning}
              min={0}
              max={1000}
              className="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-lg shadow-sm text-gray-900 dark:text-white hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400 focus:border-blue-500 dark:focus:border-blue-400 transition-all disabled:bg-gray-100 dark:disabled:bg-gray-700 disabled:cursor-not-allowed"
            />
          </div>
          
          <div className="space-y-1">
            <label className="text-sm font-medium text-gray-700 dark:text-gray-300 flex items-center gap-2">
              <input
                type="checkbox"
                checked={concurrentScan}
                onChange={(e) => setConcurrentScan(e.target.checked)}
                disabled={scanning}
                className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
              />
              並發掃描
            </label>
            {concurrentScan && (
              <input
                type="number"
                value={concurrency}
                onChange={(e) => setConcurrency(Math.max(1, Math.min(10, parseInt(e.target.value) || 3)))}
                disabled={scanning}
                min={1}
                max={10}
                placeholder="並發數量"
                className="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-lg shadow-sm text-gray-900 dark:text-white hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400 focus:border-blue-500 dark:focus:border-blue-400 transition-all disabled:bg-gray-100 dark:disabled:bg-gray-700 disabled:cursor-not-allowed"
              />
            )}
          </div>
        </div>
      </div>

      {/* 掃描進度 */}
      {scanning && (
        <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-xl p-4">
          <div className="flex items-center justify-between mb-2">
            <span className="text-sm font-medium text-gray-800 dark:text-gray-200">掃描進度</span>
            <span className="text-xs text-gray-600 dark:text-gray-400">
              {Math.round(scanProgress)}% ({scanResults.length} / {totalToScan})
            </span>
          </div>
          <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2.5 mb-2">
            <div
              className="bg-blue-600 h-2.5 rounded-full transition-all duration-300"
              style={{ width: `${scanProgress}%` }}
            ></div>
          </div>
          {currentScanning && (
            <p className="text-xs text-gray-600 dark:text-gray-400">
              正在掃描: {currentScanning}
            </p>
          )}
        </div>
      )}

      {/* 操作按鈕 */}
      <div className="flex gap-3">
        {!scanning ? (
          <button
            onClick={handleScan}
            disabled={Object.keys(validationErrors).length > 0 || (!isMCProtocol && !connectionId)}
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
              進度: {scanResults.length} / {totalToScan}
            </span>
          </div>
          <div className="grid grid-cols-4 gap-4">
            <div className="bg-white dark:bg-gray-800 rounded-lg p-3 shadow-sm">
              <div className="text-xs text-gray-500 dark:text-gray-400 mb-1">發現設備</div>
              <div className="text-2xl font-bold text-green-600 dark:text-green-400">{scanStatistics.foundDevices}</div>
            </div>
            <div className="bg-white dark:bg-gray-800 rounded-lg p-3 shadow-sm">
              <div className="text-xs text-gray-500 dark:text-gray-400 mb-1">失敗掃描</div>
              <div className="text-2xl font-bold text-red-600 dark:text-red-400">{scanStatistics.failedScans}</div>
            </div>
            <div className="bg-white dark:bg-gray-800 rounded-lg p-3 shadow-sm">
              <div className="text-xs text-gray-500 dark:text-gray-400 mb-1">平均響應</div>
              <div className="text-lg font-bold text-gray-800 dark:text-white">
                {scanStatistics.averageResponseTime > 0 
                  ? `${Math.round(scanStatistics.averageResponseTime)}ms`
                  : '-'}
              </div>
            </div>
            <div className="bg-white dark:bg-gray-800 rounded-lg p-3 shadow-sm">
              <div className="text-xs text-gray-500 dark:text-gray-400 mb-1">掃描速度</div>
              <div className="text-lg font-bold text-gray-800 dark:text-white">
                {scanStatistics.scanSpeed > 0 
                  ? `${scanStatistics.scanSpeed.toFixed(1)}/秒`
                  : '-'}
              </div>
            </div>
          </div>
        </div>
      )}

      {/* 結果過濾和搜索 */}
      {scanResults.length > 0 && (
        <div className="bg-gray-50 dark:bg-gray-900/50 p-4 rounded-xl border border-gray-200 dark:border-gray-700">
          <div className="flex flex-col sm:flex-row gap-3 mb-3">
            <div className="flex-1">
              <input
                type="text"
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                placeholder={`搜索${isMCProtocol ? 'IP' : '站號'}...`}
                className="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400"
              />
            </div>
            <div className="flex gap-2">
              <select
                value={filterType}
                onChange={(e) => setFilterType(e.target.value as 'all' | 'success' | 'failed')}
                className="px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400"
              >
                <option value="all">全部</option>
                <option value="success">成功</option>
                <option value="failed">失敗</option>
              </select>
              <select
                value={sortOrder}
                onChange={(e) => setSortOrder(e.target.value as 'asc' | 'desc')}
                className="px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400"
              >
                <option value="asc">升序</option>
                <option value="desc">降序</option>
              </select>
            </div>
          </div>
        </div>
      )}

      {/* 掃描結果列表 */}
      {filteredResults.length > 0 && (
        <div className="bg-gray-50 dark:bg-gray-900/50 p-4 rounded-xl border border-gray-200 dark:border-gray-700 animate-fade-in">
          <div className="flex justify-between items-center mb-3">
            <h4 className="text-xs font-bold text-gray-500 dark:text-gray-400 uppercase">掃描結果</h4>
            <span className="text-xs text-gray-400 dark:text-gray-500">
              顯示 {filteredResults.length} 筆結果
            </span>
          </div>
          <div className="max-h-64 overflow-y-auto scrollbar-thin space-y-2">
            {filteredResults.slice(0, 100).map((result, idx) => (
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
                      {isMCProtocol ? 'IP' : '站號'}: {result.stationOrIp}
                      {result.testAddress !== undefined && (
                        <span className="text-xs text-gray-500 dark:text-gray-400 ml-2">
                          (測試位置: {result.testAddress})
                        </span>
                      )}
                    </div>
                    {result.error && (
                      <div className="text-xs text-red-600 dark:text-red-400 mt-1">{result.error}</div>
                    )}
                    {result.success && result.data && (
                      <div className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                        數據: {JSON.stringify(result.data.values || result.data).slice(0, 50)}
                      </div>
                    )}
                    {result.responseTime && (
                      <div className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                        響應時間: {result.responseTime}ms
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
          {filteredResults.length > 100 && (
            <div className="mt-3 text-center text-xs text-gray-500 dark:text-gray-400">
              還有 {filteredResults.length - 100} 筆結果未顯示
            </div>
          )}
        </div>
      )}

      {/* 提示信息 */}
      {!isMCProtocol && !connectionId && (
        <div className="bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-lg p-3 text-sm text-amber-800 dark:text-amber-400">
          <span className="font-medium">提示：</span>請先建立連線後再進行掃描。
        </div>
      )}
      
      {validationErrors.connection && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-3 text-sm text-red-800 dark:text-red-400">
          {validationErrors.connection}
        </div>
      )}
      
      {validationErrors.scan && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-3 text-sm text-red-800 dark:text-red-400">
          {validationErrors.scan}
        </div>
      )}
    </div>
  )
}
