
import { useMemo } from 'react';
import {
  DndContext,
  closestCenter,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
  DragEndEvent,
} from '@dnd-kit/core';
import {
  arrayMove,
  SortableContext,
  sortableKeyboardCoordinates,
  verticalListSortingStrategy,
  useSortable,
} from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import type { Device } from '../../types/datalink';

// --- Icons (Lucide) ---
// Using SVGs directly to avoid adding dependency if lucide-react is not installed,
// but ideally should use lucide-react. For now, inline SVGs to match project style.
const Icons = {
  Device: () => (
    <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01" />
    </svg>
  ),
  DragHandle: () => (
    <svg className="w-4 h-4 text-gray-400 cursor-grab active:cursor-grabbing" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 8h16M4 16h16" />
    </svg>
  ),
  ChevronRight: () => (
    <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
    </svg>
  ),
  ChevronDown: () => (
    <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
    </svg>
  )
};

interface DeviceTreeNavProps {
  devices: Device[];
  selectedDeviceId: string | null;
  onSelectDevice: (deviceId: string) => void;
  isCollapsed: boolean;
  onToggleCollapse: () => void;
  onReorder?: (deviceIds: string[]) => void;
}

interface SortableDeviceItemProps {
  device: Device;
  isSelected: boolean;
  onSelect: () => void;
  isCollapsed: boolean;
}

const SortableDeviceItem = ({ device, isSelected, onSelect, isCollapsed }: SortableDeviceItemProps) => {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: device.id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.4 : 1,
    zIndex: isDragging ? 999 : 'auto',
  };

  return (
    <div
      ref={setNodeRef}
      style={style}
      className={`
        group flex items-center gap-2 px-3 py-2 rounded-lg cursor-pointer transition-colors mb-1
        ${isSelected 
          ? 'bg-slate-700 text-white border-l-2 border-blue-500' 
          : 'text-slate-300 hover:bg-slate-700/50 hover:text-slate-100 border-l-2 border-transparent'}
        ${isDragging ? 'ring-2 ring-blue-500 bg-blue-500/20' : ''}
      `}
      onClick={onSelect}
    >
      {/* Drag Handle */}
      <div {...attributes} {...listeners} className="opacity-0 group-hover:opacity-100 transition-opacity p-1">
        <Icons.DragHandle />
      </div>

      {/* Basic Icon */}
      <Icons.Device />

      {/* Label */}
      {!isCollapsed && (
        <span className="text-sm font-medium truncate flex-1">
          {device.name}
        </span>
      )}
      
      {/* Status Dot */}
      {!isCollapsed && device.status === 'active' && (
        <span className="w-2 h-2 rounded-full bg-green-500 shadow-sm ml-auto"></span>
      )}
    </div>
  );
};

export function DeviceTreeNav({
  devices,
  selectedDeviceId,
  onSelectDevice,
  isCollapsed,
  onToggleCollapse,
  onReorder
}: DeviceTreeNavProps) {
  const sensors = useSensors(
    useSensor(PointerSensor),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    })
  );

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;

    if (active.id !== over?.id && onReorder) {
      const oldIndex = devices.findIndex((d) => d.id === active.id);
      const newIndex = devices.findIndex((d) => d.id === over?.id);
      
      const newOrder = arrayMove(devices, oldIndex, newIndex).map(d => d.id);
      onReorder(newOrder);
    }
  };
  
  // SortableContext needs IDs
  const deviceIds = useMemo(() => devices.map(d => d.id), [devices]);

  return (
    <div 
      className={`flex flex-col h-full bg-slate-900 border-r border-slate-700 transition-all duration-300 ${isCollapsed ? 'w-16' : 'w-60'}`}
      data-testid="device-tree-nav"
    >
      {/* Header */}
      <div className="h-14 flex items-center justify-between px-4 border-b border-slate-700">
        {!isCollapsed && <span className="text-slate-100 font-semibold tracking-wide">DEVICES</span>}
        <button 
          onClick={onToggleCollapse}
          className="p-1.5 text-slate-400 hover:text-white rounded-md hover:bg-slate-800 transition-colors"
        >
          {isCollapsed ? <Icons.ChevronRight /> : <Icons.ChevronDown />}
        </button>
      </div>

      {/* List */}
      <div className="flex-1 overflow-y-auto p-2">
        <DndContext
          sensors={sensors}
          collisionDetection={closestCenter}
          onDragEnd={handleDragEnd}
        >
          <SortableContext 
            items={deviceIds}
            strategy={verticalListSortingStrategy}
          >
            {devices.map((device) => (
              <SortableDeviceItem
                key={device.id}
                device={device}
                isSelected={selectedDeviceId === device.id}
                onSelect={() => onSelectDevice(device.id)}
                isCollapsed={isCollapsed}
              />
            ))}
          </SortableContext>
        </DndContext>
      </div>
    </div>
  );
};
