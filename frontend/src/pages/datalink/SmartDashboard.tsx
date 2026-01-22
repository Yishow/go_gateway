
import { useState } from 'react';
import type { Device } from '../../../types/datalink';
import { DeviceTreeNav } from '../../components/datalink/DeviceTreeNav';
import { MemoryGrid } from '../../components/datalink/MemoryGrid';
import { QuickActions } from '../../components/datalink/QuickActions';
import { SlidePanel } from '../../components/datalink/SlidePanel';

export default function SmartDashboard() {
  const [selectedDeviceId, setSelectedDeviceId] = useState<string | null>(null);
  const [selectedAddresses, setSelectedAddresses] = useState<string[]>([]);
  const [isTreeCollapsed, setIsTreeCollapsed] = useState(false);
  const [isPanelOpen, setIsPanelOpen] = useState(false);
  
  // Mocks
  const devices: any[] = [ // Keep any for quick mock to avoid full typed object
    { id: '1', name: 'PLC-001', protocol: 'modbus_tcp', status: 'active' },
    { id: '2', name: 'FAT-002', protocol: 'fatek_fbs', status: 'draft' }
  ];
  
  const selectedDevice = devices.find(d => d.id === selectedDeviceId) || null;

  return (
    <div className="flex h-[calc(100vh-64px)] overflow-hidden bg-slate-900 border-t border-slate-700">
      {/* Left: Device Tree */}
      <DeviceTreeNav 
        devices={devices}
        selectedDeviceId={selectedDeviceId}
        onSelectDevice={setSelectedDeviceId}
        isCollapsed={isTreeCollapsed}
        onToggleCollapse={() => setIsTreeCollapsed(!isTreeCollapsed)}
        onReorder={(newOrder) => console.log('Reorder', newOrder)}
      />
      
      {/* Main: Memory Grid */}
      <div className="flex-1 flex flex-col bg-slate-50 dark:bg-[#0f172a] relative overflow-hidden">
        {selectedDevice ? (
          <div className="flex-1 overflow-auto p-4">
             <div className="mb-4">
               <h2 className="text-xl font-bold text-slate-800 dark:text-white mb-1">
                 {selectedDevice.name}
               </h2>
               <p className="text-sm text-slate-500 dark:text-slate-400">
                 Protocol: {selectedDevice.protocol}
               </p>
             </div>
             
             <MemoryGrid
               deviceId={selectedDevice.id}
               protocol={selectedDevice.protocol as any}
               centerAddress={selectedDevice.protocol.startsWith('modbus') ? '40001' : 'D0'}
               range={200}
               existingPoints={[]}
               selectedAddresses={selectedAddresses}
               onSelect={setSelectedAddresses}
               onCellClick={(addr) => {
                 setSelectedAddresses([addr]);
                 setIsPanelOpen(true);
               }}
             />
          </div>
        ) : (
          <div className="flex-1 flex items-center justify-center text-slate-400">
            請從左側選擇設備
          </div>
        )}
      </div>
      
      {/* Right: Quick Actions (Collapsible logic or always visible on detailed view?) */}
      {/* The requirement says 3-column. But quick actions might be hidden if screen small or just floating? */}
      {/* Let's put it as a sidebar on the right */}
      <div className="w-72 bg-white dark:bg-slate-900 border-l border-slate-200 dark:border-slate-800 p-4">
        <QuickActions
          device={selectedDevice}
          selectedCount={selectedAddresses.length}
          onBatchCreate={() => setIsPanelOpen(true)}
          onQuickMapping={() => {}}
          onTestConnection={() => {}}
        />
      </div>
      
      {/* Slide Panel Overlay */}
      <SlidePanel
        isOpen={isPanelOpen}
        title="配置點位"
        onClose={() => setIsPanelOpen(false)}
      >
        <div>Panel Content Form</div>
      </SlidePanel>
    </div>
  );
}
