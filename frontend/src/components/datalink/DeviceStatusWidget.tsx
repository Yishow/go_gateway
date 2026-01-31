import { useMemo } from 'react';
import { useDevicesQuery } from '../../hooks/datalink';
import { Link } from 'react-router-dom';

export default function DeviceStatusWidget() {
  const { data: devices = [], isLoading } = useDevicesQuery();

  const stats = useMemo(() => {
    const total = devices.length;
    const active = devices.filter(d => d.status === 'active').length;
    const error = devices.filter(d => {
        // Check readiness status if available
        if (d.readiness_status) {
            try {
                const status = JSON.parse(d.readiness_status);
                return status.status === 'error';
            } catch { return false; }
        }
        return false;
    }).length;
    const warning = devices.filter(d => {
         if (d.readiness_status) {
            try {
                const status = JSON.parse(d.readiness_status);
                return status.status === 'warning';
            } catch { return false; }
        }
        return false;
    }).length;

    return { total, active, error, warning };
  }, [devices]);

  if (isLoading) return <div className="h-48 animate-pulse bg-slate-800 rounded-2xl"></div>;

  return (
    <div className="bg-slate-800 rounded-2xl border border-slate-700 p-6 shadow-lg">
      <div className="flex justify-between items-center mb-6">
        <h3 className="text-lg font-bold text-slate-100">Device Health</h3>
        <Link to="/datalink/devices" className="text-sm text-blue-400 hover:text-blue-300">View All</Link>
      </div>

      <div className="grid grid-cols-2 gap-4">
         <div className="p-4 bg-slate-900 rounded-xl border border-slate-700">
             <div className="text-slate-400 text-sm mb-1">Total Devices</div>
             <div className="text-2xl font-bold text-slate-100">{stats.total}</div>
         </div>
         <div className="p-4 bg-green-500/10 rounded-xl border border-green-500/20">
             <div className="text-green-400 text-sm mb-1">Active</div>
             <div className="text-2xl font-bold text-green-400">{stats.active}</div>
         </div>
         <div className="p-4 bg-yellow-500/10 rounded-xl border border-yellow-500/20">
             <div className="text-yellow-400 text-sm mb-1">Warnings</div>
             <div className="text-2xl font-bold text-yellow-400">{stats.warning}</div>
         </div>
         <div className="p-4 bg-red-500/10 rounded-xl border border-red-500/20">
             <div className="text-red-400 text-sm mb-1">Errors</div>
             <div className="text-2xl font-bold text-red-400">{stats.error}</div>
         </div>
      </div>
    </div>
  );
}
