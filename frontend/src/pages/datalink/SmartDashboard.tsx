import { useState, useMemo, useCallback } from 'react';
import type { Point } from '../../types/datalink';
import { DeviceTreeNav } from '../../components/datalink/DeviceTreeNav';
import { MemoryGrid } from '../../components/datalink/MemoryGrid';
import { QuickActions } from '../../components/datalink/QuickActions';
import { SlidePanel } from '../../components/datalink/SlidePanel';
import { BatchPointCreator } from '../../components/datalink/BatchPointCreator';
import { PointDetailPanel } from '../../components/datalink/PointDetailPanel';
import { useDevicesQuery } from '../../hooks/datalink/useDevices';
import { usePollingGroupsQuery } from '../../hooks/datalink/usePollingGroups';
import { usePointsQuery } from '../../hooks/datalink/usePoints';
import { useSmartDashboardShortcuts } from '../../hooks/useKeyboardShortcuts';
import { Search, Bell, Settings, Box, Cpu, Sparkles, Keyboard } from 'lucide-react';

export default function SmartDashboard() {
  const [selectedDeviceId, setSelectedDeviceId] = useState<string | null>(null);
  const [selectedAddresses, setSelectedAddresses] = useState<string[]>([]);
  const [selectedPoint, setSelectedPoint] = useState<Point | null>(null);
  const [isTreeCollapsed, setIsTreeCollapsed] = useState(false);
  const [panelType, setPanelType] = useState<'batch' | 'detail' | 'shortcuts' | null>(null);
  
  // Data Fetching
  const { data: devices = [] } = useDevicesQuery();
  const { data: pollingGroups = [] } = usePollingGroupsQuery();
  const { data: allPoints = [] } = usePointsQuery({ device_id: selectedDeviceId || undefined });
  
  const selectedDevice = useMemo(() => 
    devices.find(d => d.id === selectedDeviceId) || null
  , [devices, selectedDeviceId]);

  // Keyboard shortcut handlers
  const handleBatchCreate = useCallback(() => {
    if (selectedDeviceId) {
      setPanelType('batch');
    }
  }, [selectedDeviceId]);

  const handleClosePanel = useCallback(() => {
    setPanelType(null);
    setSelectedPoint(null);
  }, []);

  const handleToggleSidebar = useCallback(() => {
    setIsTreeCollapsed(prev => !prev);
  }, []);

  // Register keyboard shortcuts
  const shortcuts = useSmartDashboardShortcuts({
    onBatchCreate: handleBatchCreate,
    onClosePanel: handleClosePanel,
    onToggleSidebar: handleToggleSidebar,
  });

  const handleCellClick = (_addr: string, point?: Point) => {
    if (point) {
      setSelectedPoint(point);
      setPanelType('detail');
    } 
    // If no point, we just let MemoryGrid handle selection. 
    // Batch panels are opened via QuickActions now.
  };

  return (
    <div className="flex flex-col h-screen bg-gradient-to-br from-[#0B0F19] via-[#111827] to-[#0F172A] text-slate-100 font-sans overflow-hidden selection:bg-blue-500/30">
      
      {/* 1. Modern Header (Glassmorphism) */}
      <header className="h-16 px-6 flex items-center justify-between shrink-0 z-30">
        <div className="flex items-center gap-3">
           <div className="flex items-center justify-center w-10 h-10 rounded-xl bg-blue-600/20 text-blue-400 ring-1 ring-blue-500/30">
             <Box className="w-6 h-6" />
           </div>
           <div>
             <h1 className="text-xl font-bold tracking-tight bg-gradient-to-r from-white to-slate-400 bg-clip-text text-transparent">
               GoGateway
             </h1>
             <p className="text-[10px] text-slate-500 font-medium tracking-wider uppercase">Industrial Data Collector</p>
           </div>
        </div>

        {/* Floating Search Bar */}
        <div className="flex-1 max-w-lg mx-8 relative group">
          <div className="absolute inset-0 bg-blue-500/20 blur-xl opacity-0 group-focus-within:opacity-100 transition-opacity rounded-full" />
          <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400 group-focus-within:text-blue-400 transition-colors" />
          <input 
            type="text" 
            placeholder="Search anything..." 
            className="w-full bg-slate-900/50 backdrop-blur-md border border-slate-700/50 rounded-full pl-11 pr-4 py-2.5 text-sm text-slate-200 placeholder:text-slate-500 focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-blue-500/50 transition-all shadow-lg shadow-black/20"
          />
        </div>

        {/* Right Actions */}
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2 p-1 bg-slate-800/50 rounded-full border border-slate-700/50">
            <button className="p-2 text-slate-400 hover:text-white hover:bg-slate-700/50 rounded-full transition-all">
              <Bell className="w-5 h-5" />
            </button>
            <button className="p-2 text-slate-400 hover:text-white hover:bg-slate-700/50 rounded-full transition-all">
              <Settings className="w-5 h-5" />
            </button>
          </div>
          <div className="w-9 h-9 rounded-full bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center text-sm font-bold shadow-lg shadow-blue-500/25 ring-2 ring-slate-900 cursor-pointer hover:ring-offset-2 hover:ring-offset-slate-900 transition-all">
            Y
          </div>
        </div>
      </header>

      {/* 2. Main Floating Layout */}
      <div className="flex-1 flex overflow-hidden p-4 gap-4 pt-0">
        
        {/* Left: Floating Device Tree */}
        <div 
          className={`flex-shrink-0 transition-all duration-500 ease-[cubic-bezier(0.34,1.56,0.64,1)] ${
            isTreeCollapsed ? 'w-20' : 'w-[260px]'
          }`}
        >
          <div className="h-full bg-slate-900/60 backdrop-blur-xl border border-white/5 rounded-3xl overflow-hidden shadow-2xl relative">
            {/* Visual Flair */}
            <div className="absolute top-0 left-0 right-0 h-40 bg-gradient-to-b from-blue-500/5 to-transparent pointer-events-none" />
            
            <DeviceTreeNav 
              devices={devices}
              selectedDeviceId={selectedDeviceId}
              onSelectDevice={setSelectedDeviceId}
              isCollapsed={isTreeCollapsed}
              onToggleCollapse={() => setIsTreeCollapsed(!isTreeCollapsed)}
              onReorder={(newOrder) => console.log('Reorder', newOrder)}
            />
          </div>
        </div>
        
        {/* Middle: Content Island */}
        <div className="flex-1 min-w-0 flex flex-col">
          <div className="flex-1 bg-slate-900/40 backdrop-blur-md border border-white/5 rounded-3xl overflow-hidden shadow-2xl flex flex-col relative">
            
            {selectedDevice ? (
              <>
                {/* Context Header */}
                <div className="h-16 px-8 flex items-center justify-between border-b border-white/5 bg-white/[0.02]">
                   <div className="flex items-center gap-4">
                     <div className="w-10 h-10 rounded-2xl bg-gradient-to-br from-emerald-500/20 to-teal-500/20 flex items-center justify-center border border-emerald-500/30">
                       <Cpu className="w-5 h-5 text-emerald-400" />
                     </div>
                     <div>
                       <div className="flex items-center gap-2">
                         <h2 className="text-lg font-bold text-white tracking-wide">{selectedDevice.name}</h2>
                         <span className="px-2 py-0.5 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-[10px] font-bold text-emerald-400 uppercase tracking-wider">
                           Active
                         </span>
                       </div>
                       <div className="flex items-center gap-2 text-xs text-slate-400">
                         <span className="font-mono opacity-70">ID: {selectedDevice.id.slice(0, 8)}</span>
                         <span>•</span>
                         <span>{selectedDevice.protocol}</span>
                       </div>
                     </div>
                   </div>
                   
                   <div className="flex gap-3">
                     <div className="px-4 py-2 rounded-xl bg-slate-800/50 border border-white/5 flex flex-col items-end min-w-[100px]">
                       <span className="text-[10px] text-slate-500 uppercase font-bold tracking-wider">Cycle Time</span>
                       <span className="text-sm font-mono text-blue-400">100 ms</span>
                     </div>
                   </div>
                </div>

                {/* Grid Canvas */}
                <div className="flex-1 overflow-auto p-8 scrollbar-thin scrollbar-thumb-slate-700/50 scrollbar-track-transparent">
                  <MemoryGrid
                    deviceId={selectedDevice.id}
                    protocol={selectedDevice.protocol as any}
                    centerAddress={selectedDevice.protocol.startsWith('modbus') ? '40001' : 'D0'}
                    range={300}
                    existingPoints={allPoints}
                    selectedAddresses={selectedAddresses}
                    onSelect={setSelectedAddresses}
                    onCellClick={handleCellClick}
                  />
                </div>
              </>
            ) : (
              // Option B Empty State (Vibrant)
              <div className="flex-1 flex flex-col items-center justify-center text-slate-400 p-8 relative overflow-hidden">
                <div className="absolute inset-0 bg-gradient-to-t from-blue-500/5 to-transparent opacity-50" />
                <div className="relative z-10 flex flex-col items-center">
                  <div className="w-32 h-32 rounded-3xl bg-gradient-to-br from-blue-500/20 to-purple-500/20 flex items-center justify-center mb-8 border border-white/5 shadow-[0_0_50px_rgba(59,130,246,0.2)]">
                    <Sparkles className="w-12 h-12 text-blue-400" />
                  </div>
                  <h3 className="text-2xl font-bold text-white mb-3">Welcome to Smart Dashboard</h3>
                  <p className="text-slate-400 text-center max-w-md mb-8 leading-relaxed text-lg">
                    Experience the next generation of industrial data collection. Select a device to visualize memory & configure points.
                  </p>
                  <div className="animate-bounce text-blue-400 opacity-50">
                    Which device would you like to configure?
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
        
        {/* Right: Floating Quick Actions */}
        <div className="w-[300px] flex-shrink-0">
          <div className="h-full bg-slate-900/60 backdrop-blur-xl border border-white/5 rounded-3xl overflow-hidden shadow-2xl">
            <QuickActions
              device={selectedDevice}
              selectedCount={selectedAddresses.length}
              onBatchCreate={() => setPanelType('batch')}
              onQuickMapping={() => {}}
              onTestConnection={() => {}}
            />
            {/* Keyboard Shortcuts Button */}
            <div className="p-4 border-t border-white/5">
              <button
                onClick={() => setPanelType('shortcuts')}
                className="w-full flex items-center justify-center gap-2 px-4 py-2 text-xs text-slate-400 hover:text-white hover:bg-slate-800/50 rounded-lg transition-colors"
              >
                <Keyboard className="w-4 h-4" />
                <span>鍵盤快捷鍵</span>
                <span className="ml-auto text-[10px] font-mono opacity-60">?</span>
              </button>
            </div>
          </div>
        </div>
      </div>
      
      {/* Slide Panel Overlay */}
      <SlidePanel
        isOpen={panelType !== null}
        title={panelType === 'batch' ? '批量建立點位' : panelType === 'shortcuts' ? '鍵盤快捷鍵' : '點位詳情'}
        onClose={() => setPanelType(null)}
      >
        {panelType === 'batch' && selectedDevice && (
          <BatchPointCreator
            deviceId={selectedDevice.id}
            protocol={selectedDevice.protocol as any}
            preselectedAddresses={selectedAddresses}
            pollingGroups={pollingGroups}
            onCreated={() => {
               setPanelType(null);
               setSelectedAddresses([]);
            }}
            onCancel={() => setPanelType(null)}
          />
        )}
        
        {panelType === 'detail' && selectedPoint && (
          <PointDetailPanel
            point={selectedPoint}
            onUpdate={() => setPanelType(null)}
            onDelete={() => setPanelType(null)}
            onClose={() => setPanelType(null)}
          />
        )}

        {panelType === 'shortcuts' && (
          <div className="space-y-4 p-4">
            <p className="text-sm text-slate-400 mb-4">
              使用以下快捷鍵提升操作效率
            </p>
            <div className="space-y-3">
              {shortcuts.map((shortcut, index) => (
                <div key={index} className="flex items-center justify-between py-2 px-3 bg-slate-800/50 rounded-lg">
                  <span className="text-sm text-slate-300">{shortcut.description}</span>
                  <kbd className="px-2 py-1 text-xs font-mono bg-slate-700 rounded border border-slate-600 text-slate-300">
                    {shortcut.ctrl && 'Ctrl+'}
                    {shortcut.alt && 'Alt+'}
                    {shortcut.shift && 'Shift+'}
                    {shortcut.key}
                  </kbd>
                </div>
              ))}
            </div>
          </div>
        )}
      </SlidePanel>
    </div>
  );
}
