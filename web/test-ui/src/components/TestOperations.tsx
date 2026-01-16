import { useState } from 'react'
import { useTestAPI } from '../services/api'
import type { BatchOperation } from '../types/api'

interface TestOperationsProps {
  protocol: string
  connectionId: string | null
}

export default function TestOperations({
  protocol,
  connectionId,
}: TestOperationsProps) {
  const [activeTab, setActiveTab] = useState<'single' | 'batch'>('single')
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
  
  const { read, write, batch } = useTestAPI()

  const isModbus = protocol.includes('modbus')
  const isFatek = protocol.includes('fatek')
  const isMCProtocol = protocol.includes('mcprotocol')

  const getOperations = () => {
    if (isModbus) {
      return [
        { value: 'read_coils', label: '讀取線圈 (Coils)' },
        { value: 'read_discrete_inputs', label: '讀取離散輸入 (Discrete Inputs)' },
        { value: 'read_holding_registers', label: '讀取保持暫存器 (Holding Registers)' },
        { value: 'read_input_registers', label: '讀取輸入暫存器 (Input Registers)' },
        { value: 'write_single_coil', label: '寫入單個線圈' },
        { value: 'write_single_register', label: '寫入單個暫存器' },
        { value: 'write_multiple_coils', label: '寫入多個線圈' },
        { value: 'write_multiple_registers', label: '寫入多個暫存器' },
      ]
    }
    if (isFatek) {
      return [
        { value: 'read_registers', label: '讀取暫存器' },
        { value: 'read_status', label: '讀取狀態' },
        { value: 'write_registers', label: '寫入暫存器' },
        { value: 'write_status', label: '寫入狀態' },
      ]
    }
    if (isMCProtocol) {
      return [
        { value: 'batch_read_word', label: '批量讀取字組' },
        { value: 'batch_read_bit', label: '批量讀取位元' },
        { value: 'batch_write_word', label: '批量寫入字組' },
        { value: 'batch_write_bit', label: '批量寫入位元' },
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
      if (operation.startsWith('read') || operation.includes('read')) {
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
    const isRead = operation.startsWith('read') || operation.includes('read')
    const newOp: BatchOperation = {
      type: isRead ? 'read' : 'write',
      read_request: isRead ? {
        connection_id: connectionId!,
        operation,
        address,
        count,
        symbol: isFatek ? symbol : undefined,
        device: isMCProtocol ? device : undefined,
      } : undefined,
      write_request: !isRead ? {
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

  return (
    <div className="space-y-4">
      <div className="flex border-b border-gray-200">
        <button
          className={`py-2 px-4 ${activeTab === 'single' ? 'border-b-2 border-blue-500 text-blue-600 font-medium' : 'text-gray-500'}`}
          onClick={() => setActiveTab('single')}
        >
          單次操作
        </button>
        <button
          className={`py-2 px-4 ${activeTab === 'batch' ? 'border-b-2 border-blue-500 text-blue-600 font-medium' : 'text-gray-500'}`}
          onClick={() => setActiveTab('batch')}
        >
          批量操作 ({batchQueue.length})
        </button>
      </div>

      <div className="bg-white p-4 rounded-lg border border-gray-200 shadow-sm space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            操作類型
          </label>
          <select
            value={operation}
            onChange={(e) => setOperation(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          >
            {getOperations().map((op) => (
              <option key={op.value} value={op.value}>
                {op.label}
              </option>
            ))}
          </select>
        </div>

        {isFatek && (
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              組件符號
            </label>
            <input
              type="text"
              value={symbol}
              onChange={(e) => setSymbol(e.target.value)}
              placeholder="D, M, X, Y 等"
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div>
        )}

        {isMCProtocol && (
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              設備代號
            </label>
            <input
              type="text"
              value={device}
              onChange={(e) => setDevice(e.target.value)}
              placeholder="D, M, X, Y 等"
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div>
        )}

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            起始位址
          </label>
          <input
            type="number"
            value={address}
            onChange={(e) => setAddress(parseInt(e.target.value) || 0)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          />
        </div>

        {(operation.startsWith('read') || operation.includes('read')) && (
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              讀取數量
            </label>
            <input
              type="number"
              value={count}
              onChange={(e) => setCount(parseInt(e.target.value) || 1)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div>
        )}

        {(operation.startsWith('write') || operation.includes('write')) && (
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              寫入值 (逗號分隔)
            </label>
            <input
              type="text"
              value={values}
              onChange={(e) => setValues(e.target.value)}
              placeholder="1,2,3 或 true,false,true"
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div>
        )}

        {activeTab === 'single' ? (
          <button
            onClick={handleExecute}
            disabled={loading || !connectionId}
            className="w-full bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed font-medium transition"
          >
            {loading ? '執行中...' : '執行單次操作'}
          </button>
        ) : (
          <div className="flex space-x-2">
            <button
              onClick={addToBatch}
              disabled={!connectionId}
              className="flex-1 bg-gray-100 text-gray-700 px-4 py-2 rounded-md hover:bg-gray-200 disabled:opacity-50 font-medium"
            >
              加入批量清單
            </button>
            <button
              onClick={handleExecuteBatch}
              disabled={loading || !connectionId || batchQueue.length === 0}
              className="flex-1 bg-green-600 text-white px-4 py-2 rounded-md hover:bg-green-700 disabled:opacity-50 font-medium transition"
            >
              {loading ? '批次執行中...' : `執行批量 (${batchQueue.length})`}
            </button>
            <button
              onClick={() => setBatchQueue([])}
              className="bg-red-50 text-red-600 px-4 py-2 rounded-md hover:bg-red-100"
            >
              重設
            </button>
          </div>
        )}
      </div>

      {activeTab === 'batch' && batchQueue.length > 0 && (
        <div className="bg-gray-50 p-3 rounded border border-dashed border-gray-300">
          <h4 className="text-xs font-bold text-gray-500 uppercase mb-2">批量清單項目</h4>
          <ul className="text-xs space-y-1">
            {batchQueue.map((op, idx) => (
              <li key={idx} className="flex justify-between text-gray-600">
                <span>{idx + 1}. {op.type === 'read' ? op.read_request?.operation : op.write_request?.operation} @ {op.type === 'read' ? op.read_request?.address : op.write_request?.address}</span>
              </li>
            ))}
          </ul>
        </div>
      )}

      {result && (
        <div className="mt-4">
          <div className="flex justify-between items-center mb-1">
            <label className="text-xs font-bold text-gray-500 uppercase">執行結果</label>
            <button onClick={() => setResult(null)} className="text-xs text-blue-500">清除</button>
          </div>
          <div className="p-3 bg-gray-900 text-green-400 rounded-md shadow-inner">
            <pre className="text-xs overflow-auto max-h-64 font-mono">
              {JSON.stringify(result, null, 2)}
            </pre>
          </div>
        </div>
      )}
    </div>
  )
}
