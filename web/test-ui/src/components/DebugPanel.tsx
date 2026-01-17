import { useEffect, useState, useRef, useMemo } from 'react'
import type { ReactElement } from 'react'
import { useDebugAPI } from '../services/api'
import { useToast } from '../contexts/ToastContext'
import { normalizeRawDataFromPacket } from '../utils/modbus'

interface DebugPanelProps {
  connectionId: string | null
}

type DisplayMode = 'hex' | 'ascii' | 'parsed'
type FilterDirection = 'all' | 'request' | 'response'

/**
 * 資料部分類型
 */
type DataPartType = 'address' | 'quantity' | 'byteCount' | 'data' | 'value'

/**
 * Modbus TCP 標頭部分類型
 */
type ModbusHeaderPartType = 'transactionId' | 'protocolId' | 'length' | 'unitId'

/**
 * 資料部分介面
 */
interface DataPart {
  start: number
  end: number
  type: DataPartType
}

export default function DebugPanel({ connectionId }: DebugPanelProps) {
  const [activeTab, setActiveTab] = useState<'packets' | 'logs'>('packets')
  const [autoScroll, setAutoScroll] = useState(true)
  const [displayMode, setDisplayMode] = useState<DisplayMode>('hex')
  const [filterDirection, setFilterDirection] = useState<FilterDirection>('all')
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedPacket, setSelectedPacket] = useState<string | null>(null)
  const [showFilters, setShowFilters] = useState(false)
  const [showClearConfirm, setShowClearConfirm] = useState(false)
  const scrollRef = useRef<HTMLDivElement>(null)
  const { packets, logs, refresh, clear } = useDebugAPI()
  const { showInfo } = useToast()

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

  /**
   * 解析 Modbus 資料結構並返回各部分範圍
   * @param functionCode 功能碼
   * @param isTx 是否為發送
   * @param dataStart 資料起始索引
   * @param dataEnd 資料結束索引
   */
  const parseModbusDataStructure = (
    _rawData: number[],
    functionCode: number,
    isTx: boolean,
    dataStart: number,
    dataEnd: number
  ): DataPart[] => {
    const parts: DataPart[] = []
    let currentIndex = dataStart

    if (isTx) {
      // TX (請求) 資料結構
      switch (functionCode) {
        case 0x01: // Read Coils
        case 0x02: // Read Discrete Inputs
        case 0x03: // Read Holding Registers
        case 0x04: // Read Input Registers
          // Starting Address(2) + Quantity(2)
          if (currentIndex + 4 <= dataEnd) {
            parts.push({ start: currentIndex, end: currentIndex + 1, type: 'address' })
            currentIndex += 2
            parts.push({ start: currentIndex, end: currentIndex + 1, type: 'quantity' })
            currentIndex += 2
          }
          // 處理可能的額外資料
          if (currentIndex < dataEnd) {
            parts.push({ start: currentIndex, end: dataEnd - 1, type: 'data' })
          }
          break

        case 0x05: // Write Single Coil
        case 0x06: // Write Single Register
          // Address(2) + Value(2)
          if (currentIndex + 4 <= dataEnd) {
            parts.push({ start: currentIndex, end: currentIndex + 1, type: 'address' })
            currentIndex += 2
            parts.push({ start: currentIndex, end: currentIndex + 1, type: 'value' })
            currentIndex += 2
          }
          // 處理可能的額外資料
          if (currentIndex < dataEnd) {
            parts.push({ start: currentIndex, end: dataEnd - 1, type: 'data' })
          }
          break

        case 0x0F: // Write Multiple Coils
        case 0x10: // Write Multiple Registers
          // Starting Address(2) + Quantity(2) + Byte Count(1) + Values(N)
          if (currentIndex + 5 <= dataEnd) {
            parts.push({ start: currentIndex, end: currentIndex + 1, type: 'address' })
            currentIndex += 2
            parts.push({ start: currentIndex, end: currentIndex + 1, type: 'quantity' })
            currentIndex += 2
            parts.push({ start: currentIndex, end: currentIndex, type: 'byteCount' })
            currentIndex += 1
            // 確保所有剩餘資料都被標記
            if (currentIndex < dataEnd) {
              parts.push({ start: currentIndex, end: dataEnd - 1, type: 'data' })
            }
          } else {
            // 如果資料不足，至少標記現有的
            if (currentIndex < dataEnd) {
              parts.push({ start: currentIndex, end: dataEnd - 1, type: 'data' })
            }
          }
          break

        default:
          // 未知功能碼，將所有資料標記為 data
          if (currentIndex < dataEnd) {
            parts.push({ start: currentIndex, end: dataEnd - 1, type: 'data' })
          }
          break
      }
    } else {
      // RX (回應) 資料結構
      switch (functionCode) {
        case 0x01: // Read Coils
        case 0x02: // Read Discrete Inputs
        case 0x03: // Read Holding Registers
        case 0x04: // Read Input Registers
          // Byte Count(1) + Data(N)
          if (currentIndex < dataEnd) {
            parts.push({ start: currentIndex, end: currentIndex, type: 'byteCount' })
            currentIndex += 1
            // 確保所有剩餘資料都被標記為 data
            if (currentIndex < dataEnd) {
              parts.push({ start: currentIndex, end: dataEnd - 1, type: 'data' })
            }
          }
          break

        case 0x05: // Write Single Coil
        case 0x06: // Write Single Register
          // Address(2) + Value(2) - 回應與請求相同
          if (currentIndex + 4 <= dataEnd) {
            parts.push({ start: currentIndex, end: currentIndex + 1, type: 'address' })
            currentIndex += 2
            parts.push({ start: currentIndex, end: currentIndex + 1, type: 'value' })
            currentIndex += 2
          }
          // 處理可能的額外資料
          if (currentIndex < dataEnd) {
            parts.push({ start: currentIndex, end: dataEnd - 1, type: 'data' })
          }
          break

        case 0x0F: // Write Multiple Coils
        case 0x10: // Write Multiple Registers
          // Starting Address(2) + Quantity(2)
          if (currentIndex + 4 <= dataEnd) {
            parts.push({ start: currentIndex, end: currentIndex + 1, type: 'address' })
            currentIndex += 2
            parts.push({ start: currentIndex, end: currentIndex + 1, type: 'quantity' })
            currentIndex += 2
          }
          // 處理可能的額外資料
          if (currentIndex < dataEnd) {
            parts.push({ start: currentIndex, end: dataEnd - 1, type: 'data' })
          }
          break

        default:
          // 未知功能碼，將所有資料標記為 data
          if (currentIndex < dataEnd) {
            parts.push({ start: currentIndex, end: dataEnd - 1, type: 'data' })
          }
          break
      }
    }

    return parts
  }

  /**
   * 獲取資料部分的顏色
   * @param type 資料類型
   * @param isTx 是否為發送
   */
  const getDataPartColor = (type: DataPartType, isTx: boolean): string => {
    if (isTx) {
      // TX 顏色方案
      switch (type) {
        case 'address': return 'text-purple-400'
        case 'quantity': return 'text-pink-400'
        case 'byteCount': return 'text-indigo-400'
        case 'value': return 'text-cyan-400'
        case 'data': return 'text-cyan-400'
        default: return 'text-cyan-400'
      }
    } else {
      // RX 顏色方案
      switch (type) {
        case 'address': return 'text-orange-400'
        case 'quantity': return 'text-amber-400'
        case 'byteCount': return 'text-yellow-400'
        case 'value': return 'text-yellow-400'
        case 'data': return 'text-yellow-400'
        default: return 'text-yellow-400'
      }
    }
  }

  /**
   * 格式化 hex 數據並高亮 Modbus 指令與資料
   * @param hexData hex 字符串
   * @param rawData 原始字節數組
   * @param direction 方向（request=TX, response=RX）
   */
  const formatHexWithModbusHighlight = (
    hexData: string, 
    rawData?: any, 
    direction?: string
  ): ReactElement => {
    const dataArray = normalizeRawDataFromPacket(rawData, hexData)

    if (dataArray.length === 0) {
      return <span>{hexData}</span>
    }

    const isTx = direction === 'request'
    
    // 將 hex 字符串分割成字節對，處理多種分隔符
    let hexBytes = hexData.split(/\s+/).filter(b => b.length > 0)
    
    // 如果分割後長度不一致，嘗試從 dataArray 生成 hex 字符串
    if (hexBytes.length !== dataArray.length) {
      hexBytes = dataArray.map((b: number) => {
        const num = typeof b === 'number' ? b : parseInt(b, 10)
        return num.toString(16).toUpperCase().padStart(2, '0')
      })
    }
    
    // 檢查是否為 Modbus TCP (至少 8 bytes: MBAP Header(7) + Function Code(1))
    const isModbusTCP = dataArray.length >= 8 && 
                       dataArray[2] === 0x00 && dataArray[3] === 0x00 // Protocol ID = 0x0000
    
    // 檢查是否為 Modbus RTU (至少 4 bytes: Address(1) + Function Code(1) + CRC(2))
    const isModbusRTU = dataArray.length >= 4 && 
                       !isModbusTCP // 不是 TCP 且長度足夠
    
    if (!isModbusTCP && !isModbusRTU) {
      return <span>{hexData}</span>
    }

    // 定義指令部分和資料部分的範圍
    let headerStart = 0
    let headerEnd = 0
    let functionCodeIndex = -1
    let dataStart = 0
    let dataEnd = dataArray.length  // 預設包含所有資料
    let crcStart = -1
    let crcEnd = -1

    /**
     * 取得 MBAP 標頭區段的顏色
     * @param partType 標頭區段類型
     * @param tx 是否為發送
     */
    const getMbapHeaderColor = (partType: ModbusHeaderPartType, tx: boolean): string => {
      if (tx) {
        switch (partType) {
          case 'transactionId': return 'text-blue-300'
          case 'protocolId': return 'text-blue-400'
          case 'length': return 'text-blue-500'
          case 'unitId': return 'text-blue-600'
          default: return 'text-blue-400'
        }
      }
      switch (partType) {
        case 'transactionId': return 'text-green-300'
        case 'protocolId': return 'text-green-400'
        case 'length': return 'text-green-500'
        case 'unitId': return 'text-green-600'
        default: return 'text-green-400'
      }
    }

    /**
     * 建立 MBAP 標頭索引對應資訊
     * @param tx 是否為發送
     */
    const buildMbapHeaderIndexMap = (tx: boolean): Map<number, { label: string; color: string }> => {
      const headerMap = new Map<number, { label: string; color: string }>()
      const headerParts: Array<{ start: number; end: number; type: ModbusHeaderPartType; label: string }> = [
        { start: 0, end: 1, type: 'transactionId', label: 'Transaction ID' },
        { start: 2, end: 3, type: 'protocolId', label: 'Protocol ID' },
        { start: 4, end: 5, type: 'length', label: 'Length' },
        { start: 6, end: 6, type: 'unitId', label: 'Unit ID' }
      ]

      headerParts.forEach(part => {
        for (let i = part.start; i <= part.end; i++) {
          headerMap.set(i, { label: part.label, color: getMbapHeaderColor(part.type, tx) })
        }
      })

      return headerMap
    }
    
    const mbapHeaderMap = isModbusTCP ? buildMbapHeaderIndexMap(isTx) : new Map<number, { label: string; color: string }>()

    if (isModbusTCP) {
      // Modbus TCP: MBAP Header (7 bytes) + Function Code(1) + Data(N)
      headerStart = 0
      headerEnd = 6  // MBAP Header (0-6)
      functionCodeIndex = 7
      dataStart = 8  // Function Code 之後開始是資料
      dataEnd = dataArray.length  // 包含所有資料（確保 RX 完整顯示）
    } else if (isModbusRTU) {
      // Modbus RTU: Address(1) + Function Code(1) + Data(N) + CRC(2)
      headerStart = 0
      headerEnd = 0  // Address
      functionCodeIndex = 1
      dataStart = 2  // Function Code 之後開始是資料
      // 確保至少有 CRC 的空間
      if (dataArray.length >= 4) {
        dataEnd = dataArray.length - 2  // 排除 CRC
        crcStart = dataArray.length - 2
        crcEnd = dataArray.length - 1
      } else {
        dataEnd = dataArray.length
      }
    }
    
    // 獲取功能碼（處理異常回應）
    let functionCode = dataArray[functionCodeIndex]
    const isException = functionCode >= 0x81 && functionCode <= 0x90
    if (isException) {
      functionCode = functionCode - 0x80  // 還原原始功能碼
    }
    
    // 解析資料結構（只在有資料時解析）
    const dataParts: DataPart[] = []
    if (dataStart < dataEnd) {
      const parsedParts = parseModbusDataStructure(dataArray, functionCode, isTx, dataStart, dataEnd)
      dataParts.push(...parsedParts)
    }
    
    // 指令部分顏色（Header + Function Code）
    const commandColor = isTx ? 'text-blue-400' : 'text-green-400'
    const functionCodeColor = isTx ? 'text-blue-500 font-bold' : 'text-green-500 font-bold'
    
    // 創建索引到資料類型的映射
    const indexToType = new Map<number, DataPartType>()
    dataParts.forEach(part => {
      for (let i = part.start; i <= part.end && i < dataArray.length; i++) {
        indexToType.set(i, part.type)
      }
    })
    
    // 確保所有資料部分都被標記（如果解析失敗，使用默認類型）
    if (dataParts.length === 0 && dataStart < dataEnd) {
      // 如果沒有解析出任何部分，將整個資料區域標記為 data
      for (let i = dataStart; i < dataEnd; i++) {
        if (!indexToType.has(i)) {
          indexToType.set(i, 'data')
        }
      }
    }
    
    return (
      <span>
        {hexBytes.map((byte, index) => {
          // 確保索引在範圍內
          if (index >= dataArray.length) {
            return (
              <span key={index} className="text-gray-400">
                {byte}
                {index < hexBytes.length - 1 && ' '}
              </span>
            )
          }
          
          // 指令部分（Header - MBAP）
          if (index >= headerStart && index <= headerEnd) {
            const headerInfo = mbapHeaderMap.get(index)
            return (
              <span
                key={index}
                className={`${headerInfo?.color || commandColor} font-semibold`}
                title={headerInfo?.label || 'Header'}
              >
                {byte}
                {index < hexBytes.length - 1 && ' '}
                {isModbusTCP && index === headerEnd && (
                  <span className="ml-1 text-[10px] text-gray-500">MBAP(TID/PID/LEN/UID)</span>
                )}
              </span>
            )
          }
          
          // 功能碼
          if (index === functionCodeIndex) {
            return (
              <span key={index} className={`${functionCodeColor}`}>
                {byte}
                {index < hexBytes.length - 1 && ' '}
              </span>
            )
          }
          
          // CRC 部分（RTU）- 優先處理，避免被當作資料
          if (crcStart >= 0 && index >= crcStart && index <= crcEnd) {
            return (
              <span key={index} className="text-gray-500" title="CRC">
                {byte}
                {index < hexBytes.length - 1 && ' '}
                {index === crcEnd && (
                  <span className="ml-1 text-[10px] text-gray-500">CRC</span>
                )}
              </span>
            )
          }
          
          // 資料部分（根據類型使用不同顏色）
          if (index >= dataStart && index < dataEnd) {
            const dataType = indexToType.get(index) || 'data'
            const dataColor = getDataPartColor(dataType, isTx)
            return (
              <span key={index} className={dataColor}>
                {byte}
                {index < hexBytes.length - 1 && ' '}
              </span>
            )
          }
          
          // 其他部分（確保所有資料都被顯示）
          const dataColor = getDataPartColor('data', isTx)
          return (
            <span key={index} className={dataColor}>
              {byte}
              {index < hexBytes.length - 1 && ' '}
            </span>
          )
        })}
      </span>
    )
  }

  /**
   * 格式化數據包顯示
   * 返回字符串或 React 元素
   */
  const formatPacketData = (packet: any): string | ReactElement => {
    switch (displayMode) {
      case 'ascii':
        return Array.from(packet.raw_data || [])
          .map((b: any) => {
            const char = String.fromCharCode(b)
            return char >= ' ' && char <= '~' ? char : '.'
          })
          .join('')
      case 'parsed':
        return parseProtocol(packet)
      default:
        // 返回 JSX 元素以支持顏色高亮
        return formatHexWithModbusHighlight(packet.hex_data, packet.raw_data, packet.direction)
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
    <div className="flex flex-col h-[500px] bg-gray-900 dark:bg-gray-950 rounded-xl shadow-xl border border-gray-700 dark:border-gray-800 overflow-hidden font-mono text-xs transition-colors">
      {/* Header / Tabs */}
      <div className="flex items-center justify-between px-2 py-1 bg-gray-800 dark:bg-gray-900 border-b border-gray-700 dark:border-gray-800">
        <div className="flex space-x-1">
          <button
            onClick={() => setActiveTab('packets')}
            className={`px-3 py-1.5 rounded-t-md transition-colors ${
              activeTab === 'packets'
                ? 'bg-gray-700 dark:bg-gray-800 text-blue-400 font-bold'
                : 'text-gray-400 dark:text-gray-500 hover:text-gray-200 dark:hover:text-gray-300 hover:bg-gray-700/50 dark:hover:bg-gray-800/50'
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
                ? 'bg-gray-700 dark:bg-gray-800 text-amber-400 font-bold'
                : 'text-gray-400 dark:text-gray-500 hover:text-gray-200 dark:hover:text-gray-300 hover:bg-gray-700/50 dark:hover:bg-gray-800/50'
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
            className={`px-2 py-1 rounded hover:bg-gray-700 dark:hover:bg-gray-800 transition-colors ${
              showFilters ? 'bg-gray-700 dark:bg-gray-800 text-blue-400' : ''
            }`}
            title="Toggle Filters"
          >
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z" />
            </svg>
          </button>

          {/* 顯示模式切換（僅數據包） */}
          {activeTab === 'packets' && (
            <div className="flex items-center space-x-1 bg-gray-700 dark:bg-gray-800 rounded px-1">
              <button
                onClick={() => setDisplayMode('hex')}
                className={`px-2 py-0.5 rounded text-xs transition-colors ${
                  displayMode === 'hex' ? 'bg-blue-600 text-white' : 'hover:bg-gray-600 dark:hover:bg-gray-700'
                }`}
                title="Hexadecimal"
              >
                HEX
              </button>
              <button
                onClick={() => setDisplayMode('ascii')}
                className={`px-2 py-0.5 rounded text-xs transition-colors ${
                  displayMode === 'ascii' ? 'bg-blue-600 text-white' : 'hover:bg-gray-600 dark:hover:bg-gray-700'
                }`}
                title="ASCII"
              >
                ASCII
              </button>
              <button
                onClick={() => setDisplayMode('parsed')}
                className={`px-2 py-0.5 rounded text-xs transition-colors ${
                  displayMode === 'parsed' ? 'bg-blue-600 text-white' : 'hover:bg-gray-600 dark:hover:bg-gray-700'
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
            className="px-2 py-1 rounded hover:bg-gray-700 dark:hover:bg-gray-800 transition-colors"
            title="Export Data"
          >
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </button>

          {/* 清空按鈕 */}
          {showClearConfirm ? (
            <div className="flex items-center gap-2 px-2 py-1">
              <button
                onClick={async () => {
                  await clear(connectionId)
                  setShowClearConfirm(false)
                  showInfo('數據已清空')
                }}
                className="text-xs text-red-400 hover:text-red-300 transition-colors"
              >
                確認
              </button>
              <button
                onClick={() => setShowClearConfirm(false)}
                className="text-xs text-gray-400 hover:text-gray-300 transition-colors"
              >
                取消
              </button>
            </div>
          ) : (
            <button
              onClick={() => setShowClearConfirm(true)}
              className="px-2 py-1 rounded hover:bg-gray-700 dark:hover:bg-gray-800 hover:text-red-400 transition-colors"
              title="Clear Data"
            >
              <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </button>
          )}

          {/* Auto-scroll */}
          <label className="flex items-center space-x-1 cursor-pointer hover:text-gray-200">
            <input 
              type="checkbox" 
              checked={autoScroll} 
              onChange={(e) => setAutoScroll(e.target.checked)}
              className="rounded border-gray-600 dark:border-gray-700 bg-gray-700 dark:bg-gray-800 text-blue-500 focus:ring-0 w-3 h-3"
            />
            <span className="text-xs">Auto</span>
          </label>
        </div>
      </div>

      {/* 過濾器面板 */}
      {showFilters && (
        <div className="px-2 py-1.5 bg-gray-800/50 dark:bg-gray-900/50 border-b border-gray-700 dark:border-gray-800 flex items-center gap-2 flex-wrap">
          {/* 方向過濾 */}
          {activeTab === 'packets' && (
            <div className="flex items-center space-x-1">
              <span className="text-gray-500 text-xs">方向:</span>
              <select
                value={filterDirection}
                onChange={(e) => setFilterDirection(e.target.value as FilterDirection)}
                className="bg-gray-700 dark:bg-gray-800 text-gray-300 dark:text-gray-200 text-xs px-2 py-0.5 rounded border border-gray-600 dark:border-gray-700 focus:outline-none focus:border-blue-500"
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
              className="bg-gray-700 dark:bg-gray-800 text-gray-300 dark:text-gray-200 text-xs px-2 py-0.5 rounded border border-gray-600 dark:border-gray-700 focus:outline-none focus:border-blue-500 flex-1"
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
        className="flex-1 overflow-auto p-2 space-y-1 bg-gray-900 dark:bg-gray-950 scrollbar-thin scrollbar-thumb-gray-700 dark:scrollbar-thumb-gray-800 scrollbar-track-transparent"
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
                  className={`group flex gap-3 p-1.5 hover:bg-gray-800 dark:hover:bg-gray-900 rounded border transition-all cursor-pointer ${
                    selectedPacket === packet.id ? 'bg-gray-800 dark:bg-gray-900 border-blue-500' : 'border-transparent hover:border-gray-700 dark:hover:border-gray-800'
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
                      {(() => {
                        const formatted = formatPacketData(packet)
                        return typeof formatted === 'string' ? formatted : formatted
                      })()}
                    </div>
                    {selectedPacket === packet.id && (
                      <div className="mt-1 pt-1 border-t border-gray-700 dark:border-gray-800 text-xs text-gray-400">
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
                      : 'border-blue-500 bg-gray-800/30 dark:bg-gray-900/30 text-gray-300'
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