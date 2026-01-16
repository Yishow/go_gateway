import { useState } from 'react'
import { useTestAPI } from '../services/api'

interface ConfigFormProps {
  protocol: string
  mode: string
  config: Record<string, any>
  onConfigChange: (config: Record<string, any>) => void
  connectionId: string | null
  onConnectionChange: (id: string | null) => void
}

export default function ConfigForm({
  protocol,
  mode,
  config,
  onConfigChange,
  connectionId,
  onConnectionChange,
}: ConfigFormProps) {
  const [loading, setLoading] = useState(false)
  const { connect, disconnect } = useTestAPI()

  const isTCP = mode === 'tcp' || mode === 'udp'
  const isSerial = mode === 'serial'

  const handleConfigChange = (key: string, value: any) => {
    onConfigChange({ ...config, [key]: value })
  }

  const handleConnect = async () => {
    if (!config.host && isTCP) {
      alert('請輸入主機位址')
      return
    }
    if (!config.port && isTCP) {
      alert('請輸入埠號')
      return
    }
    if (!config.port && isSerial) {
      alert('請輸入串列埠名稱')
      return
    }

    setLoading(true)
    try {
      const result = await connect(protocol, config)
      onConnectionChange(result.connection_id)
    } catch (error: any) {
      alert(`連線失敗: ${error.message}`)
    } finally {
      setLoading(false)
    }
  }

  const handleDisconnect = async () => {
    if (!connectionId) return

    setLoading(true)
    try {
      await disconnect(connectionId)
      onConnectionChange(null)
    } catch (error: any) {
      alert(`斷線失敗: ${error.message}`)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="space-y-4">
      {/* TCP/UDP 配置 */}
      {isTCP && (
        <>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              主機位址
            </label>
            <input
              type="text"
              value={config.host || ''}
              onChange={(e) => handleConfigChange('host', e.target.value)}
              placeholder="192.168.1.100"
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              埠號
            </label>
            <input
              type="number"
              value={config.port || ''}
              onChange={(e) => handleConfigChange('port', parseInt(e.target.value))}
              placeholder={protocol.includes('modbus') ? '502' : protocol.includes('fatek') ? '500' : '5000'}
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div>
          {(protocol.includes('modbus') || protocol.includes('fatek')) && (
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                單元 ID / 站號
              </label>
              <input
                type="number"
                value={config.unitID || config.station || ''}
                onChange={(e) => {
                  const key = protocol.includes('modbus') ? 'unitID' : 'station'
                  handleConfigChange(key, parseInt(e.target.value))
                }}
                placeholder="1"
                className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              />
            </div>
          )}
        </>
      )}

      {/* Serial 配置 */}
      {isSerial && (
        <>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              串列埠名稱
            </label>
            <input
              type="text"
              value={config.port || ''}
              onChange={(e) => handleConfigChange('port', e.target.value)}
              placeholder="COM3 或 /dev/ttyUSB0"
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              波特率
            </label>
            <select
              value={config.baudRate || 9600}
              onChange={(e) => handleConfigChange('baudRate', parseInt(e.target.value))}
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            >
              <option value={9600}>9600</option>
              <option value={19200}>19200</option>
              <option value={38400}>38400</option>
              <option value={57600}>57600</option>
              <option value={115200}>115200</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              資料位元
            </label>
            <select
              value={config.dataBits || (protocol.includes('fatek') || protocol.includes('mcprotocol') ? 7 : 8)}
              onChange={(e) => handleConfigChange('dataBits', parseInt(e.target.value))}
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            >
              <option value={7}>7</option>
              <option value={8}>8</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              停止位元
            </label>
            <select
              value={config.stopBits || (protocol.includes('mcprotocol') ? 2 : 1)}
              onChange={(e) => handleConfigChange('stopBits', parseInt(e.target.value))}
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            >
              <option value={1}>1</option>
              <option value={2}>2</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              同位檢查
            </label>
            <select
              value={config.parity || (protocol.includes('fatek') || protocol.includes('mcprotocol') ? 'E' : 'N')}
              onChange={(e) => handleConfigChange('parity', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="N">None</option>
              <option value="E">Even</option>
              <option value="O">Odd</option>
            </select>
          </div>
          {(protocol.includes('modbus') || protocol.includes('fatek')) && (
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                單元 ID / 站號
              </label>
              <input
                type="number"
                value={config.unitID || config.station || ''}
                onChange={(e) => {
                  const key = protocol.includes('modbus') ? 'unitID' : 'station'
                  handleConfigChange(key, parseInt(e.target.value))
                }}
                placeholder="1"
                className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              />
            </div>
          )}
        </>
      )}

      {/* 連線按鈕 */}
      <div className="pt-4">
        {!connectionId ? (
          <button
            onClick={handleConnect}
            disabled={loading}
            className="w-full bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {loading ? '連線中...' : '建立連線'}
          </button>
        ) : (
          <div className="space-y-2">
            <div className="text-sm text-green-600 font-medium">
              ✓ 已連線 ({connectionId.slice(0, 8)}...)
            </div>
            <button
              onClick={handleDisconnect}
              disabled={loading}
              className="w-full bg-red-600 text-white px-4 py-2 rounded-md hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {loading ? '斷線中...' : '斷開連線'}
            </button>
          </div>
        )}
      </div>
    </div>
  )
}
