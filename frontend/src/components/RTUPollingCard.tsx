import { useState, useEffect, useMemo } from 'react'

import CardMinimizeButton from './CardMinimizeButton'
import { useCardMinimize } from '../hooks/useCardMinimize'
import PollingOperationCard, { type PollingOperation } from './PollingOperationCard'

interface RTUPollingCardProps {
  protocol: string
  connectionId: string | null
  config: Record<string, string | number | string[] | undefined>
  onConfigChange?: (config: Record<string, string | number | string[] | undefined>) => void
  onConnectionChange?: (id: string | null) => void
}

/**
 * RTU Polling Card 組件
 * 提供快速波特率設定和無限新增的 polling 操作
 */
export default function RTUPollingCard({
  protocol,
  connectionId,
  config,
}: RTUPollingCardProps) {
  const baudRate =
    typeof config.baudRate === 'number' ? config.baudRate : Number(config.baudRate ?? 9600)
  const [operations, setOperations] = useState<PollingOperation[]>([])
  const [nextId, setNextId] = useState<number>(1)
  
  const { isMinimized } = useCardMinimize()

  // 判斷協議類型
  const isModbus = useMemo(() => protocol.includes('modbus'), [protocol])
  const isFatek = useMemo(() => protocol.includes('fatek'), [protocol])
  const isMCProtocol = useMemo(() => protocol.includes('mcprotocol'), [protocol])

  /**
   * 獲取可用的操作選項
   */
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

  /**
   * 判斷操作是否為讀取操作
   */
  const isReadOperation = (operation: string) => {
    return operation.startsWith('read') || operation.includes('read') || operation.includes('batch_read')
  }


  /**
   * 新增操作
   */
  const handleAddOperation = () => {
    const operationsList = getOperations()
    if (operationsList.length === 0) return
    
    const firstOp = operationsList[0]
    const newOperation: PollingOperation = {
      id: `op-${nextId}`,
      operation: firstOp.value,
      address: 0,
      count: 10,
      values: '',
      interval: 1000,
      baudRate: Number.isNaN(baudRate) ? 9600 : baudRate, // 使用配置中的波特率作為默認值
      enabled: false,
    }
    
    setOperations(prev => [...prev, newOperation])
    setNextId(prev => prev + 1)
  }

  /**
   * 更新操作
   */
  const handleUpdateOperation = (updatedOperation: PollingOperation) => {
    setOperations(prev => 
      prev.map(op => op.id === updatedOperation.id ? updatedOperation : op)
    )
  }

  /**
   * 刪除操作
   */
  const handleDeleteOperation = (operationId: string) => {
    setOperations(prev => prev.filter(op => op.id !== operationId))
  }

  /**
   * 協議切換時重置操作類型
   */
  useEffect(() => {
    const operationsList = getOperations()
    if (operationsList.length === 0) return
    
    const firstOp = operationsList[0]
    setOperations(prev => 
      prev.map(op => ({ ...op, operation: firstOp.value }))
    )
  }, [protocol]) // eslint-disable-line react-hooks/exhaustive-deps


  /**
   * 初始化時添加一個默認操作
   */
  useEffect(() => {
    if (operations.length === 0) {
      const operationsList = getOperations()
      if (operationsList.length > 0) {
        const firstOp = operationsList[0]
        const defaultOperation: PollingOperation = {
          id: 'op-1',
          operation: firstOp.value,
          address: 0,
          count: 10,
          values: '',
          interval: 1000,
          baudRate: Number.isNaN(baudRate) ? 9600 : baudRate,
          enabled: false,
        }
        setOperations([defaultOperation])
        setNextId(2)
      }
    }
  }, []) // eslint-disable-line react-hooks/exhaustive-deps

  if (isMinimized('rtu-polling')) {
    return null
  }

  return (
    <div className="relative bg-white dark:bg-gray-800 shadow-sm border border-gray-100 dark:border-gray-700 rounded-xl p-4 transition-colors">
      {/* 最小化按鈕 */}
      <CardMinimizeButton
        cardType="rtu-polling"
        title="RTU Polling"
        summary={`${operations.length} 個操作`}
        status={connectionId ? 'connected' : 'disconnected'}
      />
      
      <h2 className="text-base font-bold text-gray-800 dark:text-white mb-3 flex items-center gap-2">
        <span className="w-1 h-5 bg-purple-500 rounded-full"></span>
        RTU Polling
      </h2>

      {/* 新增操作按鈕 - 外層 */}
      <div className="mb-2.5 flex justify-end">
        <button
          onClick={handleAddOperation}
          className="bg-purple-600 text-white px-3 py-1 rounded-md font-semibold shadow-md shadow-purple-500/30 hover:bg-purple-700 hover:shadow-purple-500/40 active:scale-[0.98] transition-all flex items-center gap-1.5 text-xs flex-shrink-0"
        >
          <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
          新增操作
        </button>
      </div>

      <div className="space-y-2.5">

        {/* Polling 操作列表 */}
        <div className="space-y-2.5">
          {operations.map((operation) => (
            <PollingOperationCard
              key={operation.id}
              operation={operation}
              onUpdate={handleUpdateOperation}
              onDelete={() => handleDeleteOperation(operation.id)}
              connectionId={connectionId}
              protocol={protocol}
              isModbus={isModbus}
              isFatek={isFatek}
              isMCProtocol={isMCProtocol}
              getOperations={getOperations}
              isReadOperation={isReadOperation}
            />
          ))}
        </div>
      </div>
    </div>
  )
}
