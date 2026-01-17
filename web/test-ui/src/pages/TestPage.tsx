import { useState, useEffect, useRef, useMemo } from 'react'
import ProtocolSelector from '../components/ProtocolSelector'
import ConfigForm from '../components/ConfigForm'
import TestOperations from '../components/TestOperations'
import DebugPanel from '../components/DebugPanel'
import MonitorControl from '../components/MonitorControl'
import ProfileSelector from '../components/ProfileSelector'
import MinimizedCardsBar from '../components/MinimizedCardsBar'
import CardMinimizeButton from '../components/CardMinimizeButton'
import { useProfiles } from '../hooks/useProfiles'
import { useCardMinimize } from '../hooks/useCardMinimize'
import type { Profile } from '../types/profile'

export default function TestPage() {
  const [selectedProtocol, setSelectedProtocol] = useState<string>('modbus_tcp')
  const [connectionMode, setConnectionMode] = useState<string>('tcp')
  const [config, setConfig] = useState<Record<string, any>>({})
  const [connectionId, setConnectionId] = useState<string | null>(null)
  const [isConfigMinimized, setIsConfigMinimized] = useState<boolean>(false)
  
  const { currentProfile, updateCurrentProfileConfig } = useProfiles()
  const { isMinimized } = useCardMinimize()
  const isProfileLoadingRef = useRef(false)
  
  /**
   * 當 Profile 切換時，載入 Profile 的配置
   */
  useEffect(() => {
    if (currentProfile && !isProfileLoadingRef.current) {
      isProfileLoadingRef.current = true
      setSelectedProtocol(currentProfile.protocol)
      setConnectionMode(currentProfile.connectionMode)
      setConfig(currentProfile.config || {})
      // 如果已連線，先斷線（因為配置變更了）
      setConnectionId((prevId) => {
        if (prevId) {
          return null
        }
        return prevId
      })
      // 重置標記
      setTimeout(() => {
        isProfileLoadingRef.current = false
      }, 100)
    }
  }, [currentProfile?.id]) // 只在 profile ID 變更時觸發

  /**
   * 處理 Profile 切換
   */
  const handleProfileChange = (profile: Profile) => {
    isProfileLoadingRef.current = true
    setSelectedProtocol(profile.protocol)
    setConnectionMode(profile.connectionMode)
    setConfig(profile.config || {})
    // 如果已連線，先斷線（因為配置變更了）
    if (connectionId) {
      setConnectionId(null)
    }
    setTimeout(() => {
      isProfileLoadingRef.current = false
    }, 100)
  }

  /**
   * 自動保存配置到當前 Profile（使用防抖）
   */
  const saveTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  
  /**
   * 當配置變更時自動保存
   */
  useEffect(() => {
    if (currentProfile && !isProfileLoadingRef.current) {
      // 清除之前的定時器
      if (saveTimerRef.current) {
        clearTimeout(saveTimerRef.current)
      }
      // 設置新的定時器（防抖：1秒後保存）
      saveTimerRef.current = setTimeout(() => {
        updateCurrentProfileConfig(selectedProtocol, connectionMode, config)
      }, 1000)
    }
    
    return () => {
      if (saveTimerRef.current) {
        clearTimeout(saveTimerRef.current)
      }
    }
  }, [selectedProtocol, connectionMode, config, currentProfile?.id, updateCurrentProfileConfig])

  /**
   * 處理配置變更
   */
  const handleConfigChange = (newConfig: Record<string, any>) => {
    setConfig(newConfig)
  }

  /**
   * 處理協議變更
   */
  const handleProtocolChange = (protocol: string) => {
    setSelectedProtocol(protocol)
  }

  /**
   * 處理連線模式變更
   */
  const handleConnectionModeChange = (mode: string) => {
    setConnectionMode(mode)
  }

  /**
   * 當連線狀態變更時，自動縮小或展開配置區塊
   */
  useEffect(() => {
    if (connectionId) {
      // 連線成功後，延遲一點時間再縮小，讓用戶看到連線成功的提示
      const timer = setTimeout(() => {
        setIsConfigMinimized(true)
      }, 500)
      return () => clearTimeout(timer)
    } else {
      // 斷線後自動展開
      setIsConfigMinimized(false)
    }
  }, [connectionId])

  /**
   * 處理配置區塊回復
   */
  const handleConfigExpand = () => {
    setIsConfigMinimized(false)
  }

  /**
   * 處理配置區塊最小化
   */
  const handleConfigMinimize = () => {
    setIsConfigMinimized(true)
  }

  /**
   * 生成配置 card 的摘要資訊
   */
  const configSummary = useMemo(() => {
    if (connectionMode === 'tcp' || connectionMode === 'udp') {
      return `${config.host || '-'}:${config.port || '-'}`
    }
    return config.port || '-'
  }, [connectionMode, config])

  return (
    <div className="space-y-6 max-w-[1600px] mx-auto pb-40">
      {/* 最小化 Card Bar */}
      <MinimizedCardsBar />
      {/* 頂部：協議選擇器和 Profile 選擇器 */}
      <div className="bg-white shadow-sm border border-gray-100 rounded-2xl p-6">
        <div className="space-y-6">
          {/* Profile 選擇器 */}
          <ProfileSelector
            currentProtocol={selectedProtocol}
            currentConnectionMode={connectionMode}
            currentConfig={config}
            onProfileChange={handleProfileChange}
          />
          
          {/* 協議選擇器和縮小後的配置區塊 */}
          <div className="flex items-start gap-4">
            {/* 協議選擇器 */}
            <div className="flex-1">
              <ProtocolSelector
                selectedProtocol={selectedProtocol}
                onProtocolChange={handleProtocolChange}
                connectionMode={connectionMode}
                onConnectionModeChange={handleConnectionModeChange}
              />
            </div>
            
            {/* 縮小後的配置區塊 - 對齊協議按鈕頂部 */}
            {isConfigMinimized && (
              <div 
                onClick={handleConfigExpand}
                className="bg-white shadow-sm border border-gray-100 rounded-xl p-4 cursor-pointer hover:shadow-md transition-all duration-200 animate-slide-up flex-shrink-0 min-w-[280px] mt-7"
              >
                <div className="flex items-center justify-between mb-3">
                  <h3 className="text-sm font-semibold text-gray-800 flex items-center gap-2">
                    <span className="w-1 h-4 bg-blue-500 rounded-full"></span>
                    連線配置
                  </h3>
                  <div className="w-2 h-2 rounded-full bg-green-500 animate-pulse"></div>
                </div>
                <div className="space-y-2 text-xs text-gray-600">
                  {connectionMode === 'tcp' || connectionMode === 'udp' ? (
                    <>
                      <div className="flex justify-between">
                        <span className="text-gray-500">主機:</span>
                        <span className="font-mono">{config.host || '-'}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-500">埠號:</span>
                        <span className="font-mono">{config.port || '-'}</span>
                      </div>
                    </>
                  ) : (
                    <div className="flex justify-between">
                      <span className="text-gray-500">串列埠:</span>
                      <span className="font-mono">{config.port || '-'}</span>
                    </div>
                  )}
                  {connectionId && (
                    <div className="flex justify-between pt-2 border-t border-gray-100">
                      <span className="text-gray-500">連線 ID:</span>
                      <span className="font-mono text-green-600">{connectionId.slice(0, 8)}...</span>
                    </div>
                  )}
                </div>
                <div className="mt-3 text-xs text-blue-600 text-center">
                  點擊展開
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 xl:grid-cols-12 gap-6 items-start">
        {/* 左側：配置與操作 (佔 7/12) */}
        <div className={`xl:col-span-7 grid grid-cols-1 ${!isConfigMinimized ? 'lg:grid-cols-2' : ''} gap-6`}>
          {/* 配置表單 - 當縮小時隱藏 */}
          {!isConfigMinimized && !isMinimized('config') && (
            <div className="relative bg-white shadow-sm border border-gray-100 rounded-2xl p-6 h-fit animate-fade-in">
              {/* 最小化按鈕 */}
              <CardMinimizeButton
                cardType="config"
                title="連線配置"
                summary={configSummary}
                status={connectionId ? 'connected' : 'disconnected'}
              />
              <h2 className="text-lg font-bold text-gray-800 mb-6 flex items-center gap-2">
                <span className="w-1.5 h-6 bg-blue-500 rounded-full"></span>
                連線配置
              </h2>
              <ConfigForm
                protocol={selectedProtocol}
                mode={connectionMode}
                config={config}
                onConfigChange={handleConfigChange}
                connectionId={connectionId}
                onConnectionChange={setConnectionId}
                onMinimize={handleConfigMinimize}
              />
            </div>
          )}

          {/* 測試操作 */}
          {!isMinimized('operations') && (
            <div className={`relative bg-white shadow-sm border border-gray-100 rounded-2xl p-6 h-fit ${isConfigMinimized ? 'lg:col-span-1' : ''}`}>
              {/* 最小化按鈕 */}
              <CardMinimizeButton
                cardType="operations"
                title="測試操作"
                summary={selectedProtocol}
                status={connectionId ? 'connected' : 'disconnected'}
              />
              <h2 className="text-lg font-bold text-gray-800 mb-6 flex items-center gap-2">
                <span className="w-1.5 h-6 bg-amber-500 rounded-full"></span>
                測試操作
              </h2>
              <TestOperations
                protocol={selectedProtocol}
                connectionId={connectionId}
              />
            </div>
          )}

          {/* 監控模式 (放在操作下方) */}
          {!isMinimized('monitor') && (
            <div className={`relative ${!isConfigMinimized ? 'lg:col-span-2' : ''} bg-white shadow-sm border border-gray-100 rounded-2xl p-6`}>
              {/* 最小化按鈕 */}
              <CardMinimizeButton
                cardType="monitor"
                title="即時監控"
                summary={connectionId ? '已連線' : '未連線'}
                status={connectionId ? 'monitoring' : 'idle'}
              />
              <MonitorControl connectionId={connectionId} protocol={selectedProtocol} />
            </div>
          )}
        </div>

        {/* 右側：Debug 面板 (佔 5/12) - 獨立顯示以獲得更好的寬度 */}
        {!isMinimized('debug') && (
          <div className="xl:col-span-5 flex flex-col gap-6 sticky top-6">
            <div className="relative">
              {/* 最小化按鈕 */}
              <CardMinimizeButton
                cardType="debug"
                title="調試面板"
                summary={connectionId ? `連線: ${connectionId.slice(0, 8)}...` : '未連線'}
                status={connectionId ? 'connected' : 'disconnected'}
              />
              <DebugPanel connectionId={connectionId} />
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
