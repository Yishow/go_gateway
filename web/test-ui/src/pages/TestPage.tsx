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
    <div className="space-y-4">
      {/* 協議選擇器 */}
      <div className="bg-white shadow rounded-lg p-4">
        <ProtocolSelector
          selectedProtocol={selectedProtocol}
          onProtocolChange={setSelectedProtocol}
          connectionMode={connectionMode}
          onConnectionModeChange={setConnectionMode}
        />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
        {/* 左側：配置表單 */}
        <div className="lg:col-span-1">
          <div className="bg-white shadow rounded-lg p-4">
            <h2 className="text-lg font-semibold mb-4">連線配置</h2>
            <ConfigForm
              protocol={selectedProtocol}
              mode={connectionMode}
              config={config}
              onConfigChange={setConfig}
              connectionId={connectionId}
              onConnectionChange={setConnectionId}
            />
          </div>
        </div>

        {/* 中間：測試操作 */}
        <div className="lg:col-span-1">
          <div className="bg-white shadow rounded-lg p-4">
            <h2 className="text-lg font-semibold mb-4">測試操作</h2>
            <TestOperations
              protocol={selectedProtocol}
              connectionId={connectionId}
            />
          </div>
        </div>

        {/* 右側：Debug 面板 */}
        <div className="lg:col-span-1">
          <div className="bg-white shadow rounded-lg p-4">
            <h2 className="text-lg font-semibold mb-4">Debug 面板</h2>
            <DebugPanel connectionId={connectionId} />
          </div>
        </div>
      </div>

      {/* 底部：監控模式 */}
      <div className="bg-white shadow rounded-lg p-4">
        <MonitorControl connectionId={connectionId} />
      </div>
    </div>
  )
}
