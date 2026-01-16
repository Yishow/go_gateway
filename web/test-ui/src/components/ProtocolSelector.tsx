interface ProtocolSelectorProps {
  selectedProtocol: string
  onProtocolChange: (protocol: string) => void
  connectionMode: string
  onConnectionModeChange: (mode: string) => void
}

export default function ProtocolSelector({
  selectedProtocol,
  onProtocolChange,
  connectionMode,
  onConnectionModeChange,
}: ProtocolSelectorProps) {
  const protocols = [
    { value: 'modbus_tcp', label: 'Modbus TCP' },
    { value: 'modbus_udp', label: 'Modbus UDP' },
    { value: 'modbus_rtu', label: 'Modbus RTU' },
    { value: 'fatek_tcp', label: 'FATEK TCP' },
    { value: 'fatek_serial', label: 'FATEK Serial' },
    { value: 'mcprotocol_tcp', label: 'MC Protocol TCP' },
    { value: 'mcprotocol_serial', label: 'MC Protocol Serial' },
  ]

  const getModes = (protocol: string) => {
    if (protocol.startsWith('modbus_')) {
      return protocol === 'modbus_rtu'
        ? [{ value: 'serial', label: 'Serial (RTU)' }]
        : [
            { value: 'tcp', label: 'TCP' },
            { value: 'udp', label: 'UDP' },
          ]
    }
    if (protocol.startsWith('fatek_')) {
      return protocol === 'fatek_serial'
        ? [{ value: 'serial', label: 'Serial' }]
        : [{ value: 'tcp', label: 'TCP' }]
    }
    if (protocol.startsWith('mcprotocol_')) {
      return protocol === 'mcprotocol_serial'
        ? [{ value: 'serial', label: 'Serial' }]
        : [{ value: 'tcp', label: 'TCP' }]
    }
    return []
  }

  const modes = getModes(selectedProtocol)

  return (
    <div className="flex items-center space-x-4">
      <div className="flex-1">
        <label className="block text-sm font-medium text-gray-700 mb-1">
          協議類型
        </label>
        <select
          value={selectedProtocol}
          onChange={(e) => {
            onProtocolChange(e.target.value)
            const newModes = getModes(e.target.value)
            if (newModes.length > 0) {
              onConnectionModeChange(newModes[0].value)
            }
          }}
          className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
        >
          {protocols.map((protocol) => (
            <option key={protocol.value} value={protocol.value}>
              {protocol.label}
            </option>
          ))}
        </select>
      </div>

      {modes.length > 1 && (
        <div className="flex-1">
          <label className="block text-sm font-medium text-gray-700 mb-1">
            連線模式
          </label>
          <select
            value={connectionMode}
            onChange={(e) => onConnectionModeChange(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          >
            {modes.map((mode) => (
              <option key={mode.value} value={mode.value}>
                {mode.label}
              </option>
            ))}
          </select>
        </div>
      )}
    </div>
  )
}
