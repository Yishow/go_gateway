import { useState } from 'react'
import { useTestAPI } from '../services/api'

interface TestOperationsProps {
  protocol: string
  connectionId: string | null
}

export default function TestOperations({
  protocol,
  connectionId,
}: TestOperationsProps) {
  const [operation, setOperation] = useState<string>('read')
  const [address, setAddress] = useState<number>(0)
  const [count, setCount] = useState<number>(10)
  const [symbol, setSymbol] = useState<string>('D')
  const [values, setValues] = useState<string>('')
  const [result, setResult] = useState<any>(null)
  const [loading, setLoading] = useState(false)
  const { read, write } = useTestAPI()

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
      if (operation.startsWith('read')) {
        const data = await read(connectionId, {
          operation,
          address,
          count,
          symbol: !isModbus ? symbol : undefined,
        })
        setResult(data)
      } else {
        const valuesArray = values.split(',').map((v) => v.trim())
        await write(connectionId, {
          operation,
          address,
          values: valuesArray,
          symbol: !isModbus ? symbol : undefined,
        })
        setResult({ status: 'success', message: '寫入成功' })
      }
    } catch (error: any) {
      setResult({ error: error.message })
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="space-y-4">
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

      {!isModbus && (
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

      {operation.startsWith('read') && (
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

      {operation.startsWith('write') && (
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

      <button
        onClick={handleExecute}
        disabled={loading || !connectionId}
        className="w-full bg-green-600 text-white px-4 py-2 rounded-md hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {loading ? '執行中...' : '執行操作'}
      </button>

      {result && (
        <div className="mt-4 p-3 bg-gray-50 rounded-md">
          <pre className="text-sm overflow-auto">
            {JSON.stringify(result, null, 2)}
          </pre>
        </div>
      )}
    </div>
  )
}
