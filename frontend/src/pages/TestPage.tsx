import { useState, useEffect, useLayoutEffect, useRef, useMemo } from 'react'
import ProtocolSelector from '../components/ProtocolSelector'
import ConfigForm from '../components/ConfigForm'
import TestOperations from '../components/TestOperations'
import DebugPanel from '../components/DebugPanel'
import MonitorControl from '../components/MonitorControl'
import ProfileSelector from '../components/ProfileSelector'
import MinimizedCardsBar from '../components/MinimizedCardsBar'
import CardMinimizeButton from '../components/CardMinimizeButton'
import DeviceScanner from '../components/DeviceScanner'
import RTUPollingCard from '../components/RTUPollingCard'
import { useProfiles } from '../hooks/useProfiles'
import { useCardMinimize, type MinimizedCardInfo } from '../hooks/useCardMinimize'
import type { Profile, ConnectionModeConfigs } from '../types/profile'

type PageConfig = Record<string, string | number | string[] | undefined>

export default function TestPage() {
  const [selectedProtocol, setSelectedProtocol] = useState<string>('modbus_tcp')
  const [connectionMode, setConnectionMode] = useState<string>('tcp')
  const [config, setConfig] = useState<PageConfig>({})
  const [connectionId, setConnectionId] = useState<string | null>(null)
  const [isConfigMinimized, setIsConfigMinimized] = useState<boolean>(false)
  
  const { currentProfile, updateCurrentProfileConfig, updateProfile, profiles, currentProfileId } = useProfiles()
  const { isMinimized, minimizeCardsBatch } = useCardMinimize()
  const isProfileLoadingRef = useRef(false)
  const isInitializedRef = useRef(false)
  
  /**
   * 從 Profile 的配置中獲取當前模式的配置
   * 確保不同模式的配置完全獨立，不會互相影響
   */
  const getConfigForMode = (profileConfig: ConnectionModeConfigs, mode: string): PageConfig => {
    if (!profileConfig || typeof profileConfig !== 'object') {
      return {}
    }
    // 如果是舊格式（直接是 Record<string, any>），返回它
    if (!('tcp' in profileConfig || 'udp' in profileConfig || 'serial' in profileConfig || 'rtu' in profileConfig)) {
      return profileConfig as PageConfig
    }
    // 新格式：從對應的模式中獲取配置
    // 確保只返回該模式的配置，如果沒有則返回完全空白的物件
    const modeConfig = profileConfig[mode as keyof ConnectionModeConfigs]
    if (!modeConfig || typeof modeConfig !== 'object') {
      return {}
    }
    // 返回深拷貝，避免引用問題
    return { ...(modeConfig as PageConfig) }
  }

  /**
   * 當 Profile 切換時，載入 Profile 的配置
   */
  useEffect(() => {
    if (currentProfile && !isProfileLoadingRef.current) {
      isProfileLoadingRef.current = true
      setSelectedProtocol(currentProfile.protocol)
      setConnectionMode(currentProfile.connectionMode)
      // 從 Profile 的配置中獲取當前模式的配置
      const modeConfig = getConfigForMode(currentProfile.config, currentProfile.connectionMode)
      setConfig(modeConfig)
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
  }, [currentProfile]) // 只在 profile 變更時觸發

  /**
   * 處理 Profile 切換
   */
  const handleProfileChange = (profile: Profile) => {
    isProfileLoadingRef.current = true
    setSelectedProtocol(profile.protocol)
    setConnectionMode(profile.connectionMode)
    // 從 Profile 的配置中獲取當前模式的配置
    const modeConfig = getConfigForMode(profile.config, profile.connectionMode)
    setConfig(modeConfig)
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
   * 確保只保存當前模式的配置，不會影響其他模式
   */
  useEffect(() => {
    // 如果正在載入 Profile 或切換模式，不觸發自動保存
    if (currentProfile && !isProfileLoadingRef.current && connectionMode) {
      // 清除之前的定時器
      if (saveTimerRef.current) {
        clearTimeout(saveTimerRef.current)
      }
      // 設置新的定時器（防抖：1秒後保存）
      saveTimerRef.current = setTimeout(() => {
        // 確保使用當前的 connectionMode 和 config
        // 只保存當前模式的配置，不會影響其他模式
        if (config && Object.keys(config).length > 0) {
          updateCurrentProfileConfig(selectedProtocol, connectionMode, config)
        }
      }, 1000)
    }
    
    return () => {
      if (saveTimerRef.current) {
        clearTimeout(saveTimerRef.current)
      }
    }
  }, [selectedProtocol, connectionMode, config, currentProfile, updateCurrentProfileConfig])

  /**
   * 處理配置變更
   */
  const handleConfigChange = (newConfig: PageConfig) => {
    setConfig(newConfig)
  }

  /**
   * 保存當前配置並載入新配置的統一處理函數
   * 確保不同模式的配置完全獨立，不會互相影響
   */
  const saveAndLoadConfig = (newProtocol: string, newMode: string) => {
    if (!currentProfile || !currentProfileId || isProfileLoadingRef.current) return
    
    // 清除自動保存的防抖定時器，避免在切換時觸發自動保存
    if (saveTimerRef.current) {
      clearTimeout(saveTimerRef.current)
      saveTimerRef.current = null
    }
    
    // 先保存當前模式的配置（如果當前模式存在且有配置）
    // 使用同步方式立即保存，確保配置不會丟失
    if (connectionMode && config && Object.keys(config).length > 0) {
      // 立即保存當前配置，不等待防抖
      // 使用 profiles 中的最新配置，確保獲取最新的狀態
      const latestProfile = profiles.find(p => p.id === currentProfileId) || currentProfile
      const existingConfig: ConnectionModeConfigs = latestProfile.config || {}
      const updatedConfig: ConnectionModeConfigs = {
        ...existingConfig,
        [connectionMode]: { ...config }, // 深拷貝，避免引用問題
      }
      
      // 直接更新 Profile，不通過防抖機制
      updateProfile(currentProfileId, {
        protocol: selectedProtocol,
        connectionMode: connectionMode,
        config: updatedConfig,
      })
    }
    
    // 更新協議和模式
    setSelectedProtocol(newProtocol)
    setConnectionMode(newMode)
    
    // 載入新模式對應的配置（如果有的話）
    // 確保只載入該模式的配置，如果沒有則使用空配置
    // 使用 profiles 中的最新配置（因為 updateProfile 可能還沒更新 profiles 狀態）
    const latestProfile = profiles.find(p => p.id === currentProfileId) || currentProfile
    const modeConfig = getConfigForMode(latestProfile.config, newMode)
    // 確保返回的是全新的物件，避免引用問題
    // 如果新模式沒有配置，確保返回完全空白的物件
    const cleanConfig = Object.keys(modeConfig).length > 0 ? { ...modeConfig } : {}
    setConfig(cleanConfig)
  }

  /**
   * 處理協議變更
   * 切換時先保存當前配置，再載入新協議對應的配置
   */
  const handleProtocolChange = (protocol: string) => {
    // 從協議字符串中提取連線模式
    // 例如: modbus_tcp -> tcp, modbus_rtu -> serial
    let mode = connectionMode
    if (protocol.includes('_tcp')) {
      mode = 'tcp'
    } else if (protocol.includes('_udp')) {
      mode = 'udp'
    } else if (protocol.includes('_rtu') || protocol.includes('_serial')) {
      mode = 'serial'
    }
    // 設置標誌，表示正在切換模式
    isModeChangingRef.current = true
    saveAndLoadConfig(protocol, mode)
    // 重置標誌
    setTimeout(() => {
      isModeChangingRef.current = false
    }, 100)
  }

  /**
   * 處理連線模式變更
   * 切換時先保存當前模式的配置，再載入新模式對應的配置
   * 注意：當從 ProtocolSelector 的連線模式按鈕點擊時，會同時調用此函數和 handleProtocolChange
   * 為了避免重複保存，使用一個標誌來防止重複調用 saveAndLoadConfig
   */
  const isModeChangingRef = useRef(false)
  
  const handleConnectionModeChange = (mode: string) => {
    // 如果正在切換模式（通過 handleProtocolChange），這裡只需要更新模式狀態
    if (isModeChangingRef.current) {
      setConnectionMode(mode)
      return
    }
    // 否則，這是一個獨立的模式變更（不應該發生，但為了安全起見保留）
    saveAndLoadConfig(selectedProtocol, mode)
  }

  /**
   * 頁面初始化時，將除連線配置外的所有卡片預設為最小化
   * 使用 useLayoutEffect 在瀏覽器繪製前設置狀態，避免 UI 跳動
   */
  useLayoutEffect(() => {
    if (!isInitializedRef.current) {
      isInitializedRef.current = true
      // 批量設置所有卡片為最小化狀態，一次性更新，減少重新渲染
      const cardsToMinimize: MinimizedCardInfo[] = [
        {
          type: 'operations',
          title: '測試操作',
          summary: selectedProtocol,
          status: connectionId ? 'connected' : 'disconnected',
        },
        {
          type: 'scanner',
          title: '設備掃描',
          summary: selectedProtocol,
          status: 'idle',
        },
        {
          type: 'monitor',
          title: '即時監控',
          summary: connectionId ? '已連線' : '未連線',
          status: connectionId ? 'monitoring' : 'idle',
        },
        {
          type: 'rtu-polling',
          title: 'RTU Polling',
          summary: `${config.baudRate || 9600} baud`,
          status: connectionId ? 'connected' : 'disconnected',
        },
        {
          type: 'debug',
          title: '調試面板',
          summary: connectionId && connectionId.length > 0 ? `連線: ${connectionId.slice(0, 8)}...` : '未連線',
          status: connectionId ? 'connected' : 'disconnected',
        },
      ]
      
      // 一次性批量設置，只觸發一次重新渲染
      minimizeCardsBatch(cardsToMinimize)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []) // 只在組件掛載時執行一次

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
      return `${String(config.host ?? '-')}:${String(config.port ?? '-')}`
    }
    return String(config.port ?? '-')
  }, [connectionMode, config])

  return (
    <div className="mx-auto max-w-[1600px] space-y-6 rounded-2xl bg-gradient-to-br from-[#0B1220] via-[#0F172A] to-[#111827] p-4 pb-40 text-slate-100 sm:p-6">
      {/* 最小化 Card Bar */}
      <MinimizedCardsBar />
      {/* 頂部：協議選擇器和 Profile 選擇器 */}
      <div className="rounded-2xl border border-white/10 bg-slate-900/70 p-6 shadow-xl shadow-black/20 transition-colors">
        <div className="space-y-6">
          <div className="flex flex-col gap-2 border-b border-white/10 pb-4">
            <p className="text-[11px] font-semibold uppercase tracking-[0.2em] text-slate-400">
              Engineer Workspace
            </p>
            <div className="flex flex-col gap-1 lg:flex-row lg:items-end lg:justify-between">
              <div>
                <h1 className="text-2xl font-bold tracking-tight text-slate-50">
                  協議測試工作台
                </h1>
                <p className="mt-1 text-sm text-slate-400">
                  保持工程測試用途不變，僅收斂視覺語言與主要資訊層級。
                </p>
              </div>
              <div className="inline-flex items-center rounded-full border border-emerald-500/30 bg-emerald-500/10 px-3 py-1 text-[11px] font-medium text-emerald-200">
                TestPage 專用工具
              </div>
            </div>
          </div>
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
                className="mt-7 min-w-[280px] flex-shrink-0 cursor-pointer rounded-xl border border-white/10 bg-slate-950/70 p-4 shadow-lg shadow-black/20 animate-slide-up transition-[border-color,background-color,box-shadow] duration-200 hover:border-blue-500/30 hover:bg-slate-900/80 hover:shadow-xl"
              >
                <div className="flex items-center justify-between mb-3">
                  <h3 className="flex items-center gap-2 text-sm font-semibold text-slate-100">
                    <span className="w-1 h-4 bg-blue-500 rounded-full"></span>
                    連線配置
                  </h3>
                  <div className="w-2 h-2 rounded-full bg-green-500 animate-pulse"></div>
                </div>
                <div className="space-y-2 text-xs text-slate-300">
                  {connectionMode === 'tcp' || connectionMode === 'udp' ? (
                    <>
                      <div className="flex justify-between">
                        <span className="text-slate-500">主機:</span>
                        <span className="font-mono">{String(config.host ?? '-')}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-slate-500">埠號:</span>
                        <span className="font-mono">{String(config.port ?? '-')}</span>
                      </div>
                    </>
                  ) : (
                    <div className="flex justify-between">
                      <span className="text-slate-500">串列埠:</span>
                      <span className="font-mono">{String(config.port ?? '-')}</span>
                    </div>
                  )}
                  {connectionId && (
                    <div className="flex justify-between border-t border-white/10 pt-2">
                      <span className="text-slate-500">連線 ID:</span>
                      <span className="font-mono text-emerald-300">{connectionId.slice(0, 8)}...</span>
                    </div>
                  )}
                </div>
                <div className="mt-3 text-center text-xs text-blue-300">
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
            <div className="relative h-fit rounded-2xl border border-white/10 bg-slate-900/70 p-6 shadow-xl shadow-black/20 animate-fade-in transition-colors">
              {/* 最小化按鈕 */}
              <CardMinimizeButton
                cardType="config"
                title="連線配置"
                summary={configSummary}
                status={connectionId ? 'connected' : 'disconnected'}
              />
              <h2 className="mb-6 flex items-center gap-2 text-lg font-bold text-slate-100">
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
            <div className={`relative h-fit rounded-2xl border border-white/10 bg-slate-900/70 p-6 shadow-xl shadow-black/20 ${isConfigMinimized ? 'lg:col-span-1' : ''} transition-colors`}>
              {/* 最小化按鈕 */}
              <CardMinimizeButton
                cardType="operations"
                title="測試操作"
                summary={selectedProtocol}
                status={connectionId ? 'connected' : 'disconnected'}
              />
              <h2 className="mb-6 flex items-center gap-2 text-lg font-bold text-slate-100">
                <span className="w-1.5 h-6 bg-amber-500 rounded-full"></span>
                測試操作
              </h2>
              <TestOperations
                protocol={selectedProtocol}
                connectionId={connectionId}
              />
            </div>
          )}

          {/* 掃描設備 */}
          {!isMinimized('scanner') && (
            <div className={`relative ${!isConfigMinimized ? 'lg:col-span-2' : ''} rounded-2xl border border-white/10 bg-slate-900/70 p-6 shadow-xl shadow-black/20 transition-colors`}>
              {/* 最小化按鈕 */}
              <CardMinimizeButton
                cardType="scanner"
                title="設備掃描"
                summary={selectedProtocol}
                status="idle"
              />
              <h2 className="mb-6 flex items-center gap-2 text-lg font-bold text-slate-100">
                <span className="w-1.5 h-6 bg-green-500 rounded-full"></span>
                設備掃描
              </h2>
              <DeviceScanner
                protocol={selectedProtocol}
                connectionId={connectionId}
                baseConfig={config}
              />
            </div>
          )}

          {/* 監控模式 (放在操作下方) */}
          {!isMinimized('monitor') && (
            <div className={`relative ${!isConfigMinimized ? 'lg:col-span-2' : ''} rounded-2xl border border-white/10 bg-slate-900/70 p-6 shadow-xl shadow-black/20 transition-colors`}>
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

          {/* RTU Polling Card */}
          {!isMinimized('rtu-polling') && (
            <div className={`relative ${!isConfigMinimized ? 'lg:col-span-2' : ''} transition-colors`}>
              <RTUPollingCard
                protocol={selectedProtocol}
                connectionId={connectionId}
                config={config}
                onConfigChange={handleConfigChange}
                onConnectionChange={setConnectionId}
              />
            </div>
          )}
        </div>

        {/* 右側：Debug 面板 (佔 5/12) - 獨立顯示以獲得更好的寬度 */}
        {!isMinimized('debug') && (
          <div className="xl:col-span-5 flex flex-col gap-6 sticky top-6">
            <div className="relative rounded-2xl border border-white/10 bg-slate-900/70 p-6 shadow-xl shadow-black/20 transition-colors">
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
