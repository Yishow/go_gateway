interface ProtocolSelectorProps {
  selectedProtocol: string
  onProtocolChange: (protocol: string) => void
  connectionMode: string
  onConnectionModeChange: (mode: string) => void
}

const ProtocolIcons = {
  Modbus: () => (
    <svg className="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
    </svg>
  ),
  Fatek: () => (
    <svg className="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
    </svg>
  ),
  MC: () => (
    <svg className="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
    </svg>
  ),
}

export default function ProtocolSelector({
  selectedProtocol,
  onProtocolChange,
  onConnectionModeChange,
}: ProtocolSelectorProps) {
  const protocols = [
    { id: 'modbus', label: 'Modbus', icon: ProtocolIcons.Modbus, variants: ['tcp', 'udp', 'rtu'] },
    { id: 'fatek', label: 'FATEK', icon: ProtocolIcons.Fatek, variants: ['tcp', 'serial'] },
    { id: 'mcprotocol', label: 'MC Protocol', icon: ProtocolIcons.MC, variants: ['tcp', 'serial'] },
  ]

  // Helper to determine active root protocol
  const activeRoot = protocols.find(p => selectedProtocol.startsWith(p.id))?.id || 'modbus'

  const getModes = (protocolId: string) => {
    const p = protocols.find(p => p.id === protocolId)
    if (!p) return []
    
    // Map internal variant names to UI labels and values
    return p.variants.map(v => {
      let value = `${protocolId}_${v}`
      let label = v.toUpperCase()
      let modeValue = v === 'rtu' || v === 'serial' ? 'serial' : v // 'tcp', 'udp'

      return { fullValue: value, mode: modeValue, label }
    })
  }

  const currentVariants = getModes(activeRoot)

  return (
    <div className="space-y-6">
      {/* 1. Protocol Family Selection */}
      <div>
        <label className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-3 block">
          選擇通訊協定
        </label>
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
          {protocols.map((p) => {
            const isActive = selectedProtocol.startsWith(p.id)
            return (
              <button
                key={p.id}
                onClick={() => {
                  // Default to first variant when switching family
                  const defaultVariant = p.variants[0]
                  const fullValue = `${p.id}_${defaultVariant}`
                  onProtocolChange(fullValue)
                  
                  // Update connection mode mapping
                  const mode = (defaultVariant === 'rtu' || defaultVariant === 'serial') ? 'serial' : defaultVariant
                  onConnectionModeChange(mode)
                }}
                className={`
                  relative flex flex-col items-center justify-center p-4 rounded-xl border-2 transition-all duration-200 ease-in-out
                  ${isActive 
                    ? 'border-blue-500 bg-blue-50 dark:bg-blue-900/20 text-blue-700 dark:text-blue-400 shadow-sm dark:border-blue-400' 
                    : 'border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 text-gray-600 dark:text-gray-300 hover:border-blue-200 dark:hover:border-blue-700 hover:bg-gray-50 dark:hover:bg-gray-700'
                  }
                `}
              >
                <div className={`mb-3 ${isActive ? 'text-blue-600 dark:text-blue-400' : 'text-gray-400 dark:text-gray-500'}`}>
                  <p.icon />
                </div>
                <span className="font-semibold">{p.label}</span>
              </button>
            )
          })}
        </div>
      </div>

      {/* 2. Connection Mode Selection (Segmented Control) */}
      <div className="animate-slide-up">
        <label className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-3 block">
          連線模式
        </label>
        <div className="flex bg-gray-100 dark:bg-gray-700 p-1 rounded-lg inline-flex">
          {currentVariants.map((variant) => {
            const isSelected = selectedProtocol === variant.fullValue
            return (
              <button
                key={variant.fullValue}
                onClick={(e) => {
                  e.preventDefault()
                  // 只調用 onProtocolChange，因為協議字符串已經包含了模式信息
                  // handleProtocolChange 會自動提取並更新模式
                  onProtocolChange(variant.fullValue)
                  // 同時更新連線模式，確保狀態同步
                  onConnectionModeChange(variant.mode)
                }}
                className={`
                  px-6 py-2 rounded-md text-sm font-medium transition-all duration-200 cursor-pointer
                  active:scale-95
                  ${isSelected 
                    ? 'bg-white dark:bg-gray-600 text-blue-700 dark:text-blue-300 shadow-sm font-semibold ring-2 ring-blue-200 dark:ring-blue-500/30' 
                    : 'text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-600/50'
                  }
                `}
                aria-pressed={isSelected}
                type="button"
              >
                {variant.label}
              </button>
            )
          })}
        </div>
      </div>
    </div>
  )
}