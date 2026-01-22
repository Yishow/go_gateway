import { useTranslation } from 'react-i18next';

export default function Dashboard() {
  const { t } = useTranslation();
  
  const stats = [
    { label: t('dashboard.totalDevices'), value: '12', unit: t('dashboard.active'), color: 'blue' },
    { label: t('dashboard.dataPoints'), value: '1,240', unit: t('dashboard.monitored'), color: 'emerald' },
    { label: t('dashboard.throughput'), value: '450', unit: t('dashboard.opsPerSec'), color: 'purple' },
    { label: t('dashboard.errors'), value: '0', unit: t('dashboard.last24h'), color: 'red' },
  ];

  return (
    <div className="space-y-8">
      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {stats.map((stat) => (
          <div key={stat.label} className="bg-slate-800/50 backdrop-blur-sm p-6 rounded-2xl border border-slate-700/50 hover:border-slate-600 transition-colors">
            <p className="text-sm text-slate-400 font-medium">{stat.label}</p>
            <div className="mt-2 flex items-baseline space-x-2">
              <span className="text-3xl font-bold text-slate-100">{stat.value}</span>
              <span className={`text-xs font-medium px-2 py-0.5 rounded-full bg-${stat.color}-500/10 text-${stat.color}-400 border border-${stat.color}-500/20`}>
                {stat.unit}
              </span>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
