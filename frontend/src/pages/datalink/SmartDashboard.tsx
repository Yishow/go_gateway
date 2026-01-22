import { useState } from 'react';
import type { Device, Point } from '../../../types/datalink';
import { DeviceTreeNav } from '../../components/datalink/DeviceTreeNav';
import { MemoryGrid } from '../../components/datalink/MemoryGrid';
import { QuickActions } from '../../components/datalink/QuickActions';
import { SlidePanel } from '../../components/datalink/SlidePanel';
import { BatchPointCreator } from '../../components/datalink/BatchPointCreator';
import { PointDetailPanel } from '../../components/datalink/PointDetailPanel';

export default function SmartDashboard() {
  const [selectedDeviceId, setSelectedDeviceId] = useState<string | null>(null);
  const [selectedAddresses, setSelectedAddresses] = useState<string[]>([]);
  const [selectedPoint, setSelectedPoint] = useState<Point | null>(null);
  const [isTreeCollapsed, setIsTreeCollapsed] = useState(false);
  const [panelType, setPanelType] = useState<'batch' | 'detail' | null>(null);
  
  // Mocks
  const devices: any[] = [
    { id: '1', name: 'PLC-001', protocol: 'modbus_tcp', status: 'active' },
    { id: '2', name: 'FAT-002', protocol: 'fatek_fbs', status: 'draft' }
  ];
  
  const selectedDevice = devices.find(d => d.id === selectedDeviceId) || null;

  const handleCellClick = (addr: string, point?: Point) => {
    if (point) {
      setSelectedPoint(point);
      setPanelType('detail');
    } else {
      setSelectedAddresses([addr]);
      setPanelType('batch');
    }
  };

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
               existingPoints={[]} // TODO: Fetch points
               selectedAddresses={selectedAddresses}
               onSelect={setSelectedAddresses}
               onCellClick={handleCellClick}
             />
          </div>
        ) : (
          <div className="flex-1 flex items-center justify-center text-slate-400">
            請從左側選擇設備
          </div>
        )}
      </div>
      
      <div className="w-72 bg-white dark:bg-slate-900 border-l border-slate-200 dark:border-slate-800 p-4">
        <QuickActions
          device={selectedDevice}
          selectedCount={selectedAddresses.length}
          onBatchCreate={() => setPanelType('batch')}
          onQuickMapping={() => {}}
          onTestConnection={() => {}}
        />
      </div>
      
      {/* Slide Panel Overlay */}
      <SlidePanel
        isOpen={panelType !== null}
        title={panelType === 'batch' ? '批量建立點位' : '點位詳情'}
        onClose={() => setPanelType(null)}
      >
        {panelType === 'batch' && selectedDevice && (
          <BatchPointCreator
            deviceId={selectedDevice.id}
            protocol={selectedDevice.protocol as any}
            preselectedAddresses={selectedAddresses}
            pollingGroups={[]} // TODO: Fetch polling groups
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
      </SlidePanel>
    </div>
  );
}
