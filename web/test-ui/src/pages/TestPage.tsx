import { useState } from 'react'
import ProtocolSelector from '../components/ProtocolSelector'
import ConfigForm from '../components/ConfigForm'
import TestOperations from '../components/TestOperations'
import DebugPanel from '../components/DebugPanel'
import MonitorControl from '../components/MonitorControl'

export default function TestPage() {
  const [selectedProtocol, setSelectedProtocol] = useState<string>('modbus_tcp')
  const [connectionMode, setConnectionMode] = useState<string>('tcp')
  const [config, setConfig] = useState<Record<string, any>>({})
  const [connectionId, setConnectionId] = useState<string | null>(null)

  return (
    <div className="space-y-6 max-w-[1600px] mx-auto">
      {/* 頂部：協議選擇器 */}
      <div className="bg-white shadow-sm border border-gray-100 rounded-2xl p-6">
        <ProtocolSelector
          selectedProtocol={selectedProtocol}
          onProtocolChange={setSelectedProtocol}
          connectionMode={connectionMode}
          onConnectionModeChange={setConnectionMode}
        />
      </div>

      <div className="grid grid-cols-1 xl:grid-cols-12 gap-6 items-start">
        {/* 左側：配置與操作 (佔 7/12) */}
        <div className="xl:col-span-7 grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* 配置表單 */}
          <div className="bg-white shadow-sm border border-gray-100 rounded-2xl p-6 h-fit">
            <h2 className="text-lg font-bold text-gray-800 mb-6 flex items-center gap-2">
              <span className="w-1.5 h-6 bg-blue-500 rounded-full"></span>
              連線配置
            </h2>
            <ConfigForm
              protocol={selectedProtocol}
              mode={connectionMode}
              config={config}
              onConfigChange={setConfig}
              connectionId={connectionId}
              onConnectionChange={setConnectionId}
            />
          </div>

          {/* 測試操作 */}
          <div className="bg-white shadow-sm border border-gray-100 rounded-2xl p-6 h-fit">
            <h2 className="text-lg font-bold text-gray-800 mb-6 flex items-center gap-2">
              <span className="w-1.5 h-6 bg-amber-500 rounded-full"></span>
              測試操作
            </h2>
            <TestOperations
              protocol={selectedProtocol}
              connectionId={connectionId}
            />
          </div>

          {/* 監控模式 (放在操作下方) */}
           <div className="lg:col-span-2 bg-white shadow-sm border border-gray-100 rounded-2xl p-6">
            <MonitorControl connectionId={connectionId} />
          </div>
        </div>

        {/* 右側：Debug 面板 (佔 5/12) - 獨立顯示以獲得更好的寬度 */}
        <div className="xl:col-span-5 flex flex-col gap-6 sticky top-6">
           {/* 直接渲染 DebugPanel，不加額外 Wrapper */}
           <DebugPanel connectionId={connectionId} />
        </div>
      </div>
    </div>
  )
}
